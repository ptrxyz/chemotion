process.env.NODE_PATH = '/usr/pnpm-global/5/node_modules/'
import module from 'module'
module._initPaths()

import util from 'util'
import child_process from 'child_process'
const exec = util.promisify(child_process.exec)

import fs from 'fs'
import path from 'path'
import { globbySync } from 'globby'

import fetch from 'node-fetch'
import semver from 'semver'

import { dirSort } from './dirsort.js'
import { yamlLoad, yamlDump } from './rbyml.js'

let raw = fs.readFileSync(
    path.resolve(process.cwd(), 'prepare', 'volinfo.yaml'),
    'utf-8'
)
let volumeInfo = yamlLoad(raw)

const inSrcDir = (x) => path.join(`${config.srcDir}`, x)

// desiredNodeVersion:
//      - SemVer range
//      - "package" to use the version set in package.json
//      - "package-upgrade" to upgrade the version set in package.json to the highest possible minor release.
const config = {
    tag: process.env.CHEMOTION_BUILD_TAG || 'v1.1.2',
    release: process.env.CHEMOTION_BUILD_RELEASE || 'local-build',
    srcDir: `${process.cwd()}/src/app`,
    dataDir: `${process.cwd()}/src/data`,
    // seedDir: `${process.cwd()}/seed`,
    rawDir: `${process.cwd()}/raw`,
    fixDir: `${process.cwd()}/fixes`,
    ephemeralDir: `/dev/shm/chemotion`,
    exposeDir: `/shared`,
    cwd: process.cwd(),
    desiredNodeVersion: 'package-upgrade',
    version: { chemotion: {}, crt: {} },
    exposedPaths: [],
    needsSeeding: [],
}

const bash = async (cmd, options = {}) => {
    const bashOptions = {
        shell: '/bin/bash',
        encoding: 'utf-8',
        env: { GIT_TERMINAL_PROMPT: '0' },
    }
    return await exec(cmd, { ...bashOptions, ...options })
}

async function taskCheckoutRepo(config) {
    const tag2param = (tag) => {
        tag = tag.replace(/\s/g, '')
        return tag.length > 0 && !tag.match(/(default|latest)/i) ? `${tag}` : ''
    }

    const dirCopy = async (src, dst) => {
        return await bash(
            `tar cH posix --directory ${src} . | tar xv --one-top-level=${dst}`
        )
    }

    fs.rmSync(config.srcDir, { recursive: true, force: true })
    let repoURL =
        process.env.CHEMOTION_REPO_URL ||
        'https://github.com/ComPlat/chemotion_ELN'

    let recreateRaw = false

    if (!fs.existsSync(config.rawDir)) {
        // raw doesnt exist
        recreateRaw = true
    } else if (fs.existsSync(config.rawDir)) {
        // raw exists but it's from a different origin
        let { stdout, stderr } = await bash(
            `git config --get remote.origin.url`,
            {
                cwd: config.rawDir,
            }
        )
        if (stdout.trim() != repoURL) {
            recreateRaw = true
        }
    }

    if (recreateRaw) {
        console.log(`Cloning into [${config.srcDir}]...`)
        // https://github.com/megorei/chemotion_ELN/commit/4e4f541de9ba7204b741b8cd7864e590b6129f70
        console.log(` > Cloning from [${repoURL}].`)
        await bash(
            `git clone -c advice.detachedHead=false ${repoURL} ${config.rawDir} 2>/dev/null`
        )
    }

    console.log(`Using raw cache for [${config.srcDir}]...`)
    //fs.cpSync(config.rawDir, config.srcDir, { preserveTimestamps: true, force: true, recursive: true, dereference: true })
    fs.rmSync(config.srcDir, { recursive: true, force: true })
    fs.rmSync(config.dataDir, { recursive: true, force: true })
    await dirCopy(config.rawDir, config.srcDir)
    await bash(`git fetch --depth 1 origin ${tag2param(config.tag)}`, {
        cwd: config.srcDir,
    })
    await bash(`git checkout FETCH_HEAD`, { cwd: config.srcDir })

    // echo -e "CHEMOTION_REF=${ELNREF}\nCHEMOTION_TAG=${ELNTAG}\nBUILDSYSTEM_REF=${BLDREF}\nBUILDSYSTEM_TAG=${BLDTAG}"
    config.version.chemotion = {
        ref: (
            await bash(`git rev-parse --short HEAD || echo "unknown"`, {
                cwd: config.srcDir,
            })
        ).stdout.trim(),
        tag: (
            await bash(`git describe --abbrev=0 --tags || echo "untagged"`, {
                cwd: config.srcDir,
            })
        ).stdout.trim(),
    }
    config.version.crt = {
        ref: (
            await bash(`git rev-parse --short HEAD || echo "unknown"`, {
                cwd: config.cwd,
            })
        ).stdout.trim(),
        tag: (
            await bash(`git describe --abbrev=0 --tags || echo "untagged"`, {
                cwd: config.cwd,
            })
        ).stdout.trim(),
    }

    const env = [
        ['CHEMOTION_REF', config.version.chemotion.ref],
        ['CHEMOTION_TAG', config.version.chemotion.tag],
        ['BUILDSYSTEM_REF', config.version.crt.ref],
        ['BUILDSYSTEM_TAG', config.version.crt.ref],
        ['RELEASE', config.release],
    ]

    const versionString = env.map((x) => `${x[0]}=${x[1]}`).join('\n') + '\n'
    fs.writeFileSync(inSrcDir('.version'), versionString)
    fs.writeFileSync(inSrcDir('.metadata'), versionString)
}

async function taskCleanRepo(config) {
    // Removes "test" and "development" sections from YML files
    const removeUnusedSections = (fileName) => {
        const fileContents = fs.readFileSync(fileName, { encoding: 'utf-8' })
        const yml = yamlLoad(fileContents)
        if (!yml) return fileContents

        // we do not want development or testing configurations
        delete yml['development']
        delete yml['test']

        return yamlDump(yml)
    }

    const globPaths = [
        'config/deploy.rb',
        'config/environments/development.rb',
        'config/environments/test.rb',
        '**/*.example',
        '**/*github*',
        '**/*gitlab*',
        '**/*travis*',
        '**/*.bak',
        '**/*git*',
    ]

    const toDelete = globbySync(globPaths.map(inSrcDir), {
        dot: true,
        onlyFiles: false,
        expandDirectories: true,
    })
    toDelete.sort(dirSort).forEach((element) => {
        fs.rmSync(element, { recursive: true, force: true })
        console.log(`[${element}] was purged from repo.`)
    })

    // files to 'clean':
    // - if it's a yml file, remove all but the production sections
    const toClean = ['**/*.yml']
    for (const fileName of globbySync(toClean.map(inSrcDir), {
        dot: true,
        onlyFiles: true,
    }).sort(dirSort)) {
        const fc = removeUnusedSections(fileName)
        fs.writeFileSync(fileName, fc, { encoding: 'utf-8' })
        console.log(`Test- and dev-sections cleaned from ${fileName}`)
    }
}

async function taskDetermineNodeVersion(config) {
    // determine what node version we should use.
    const getDesiredVersion = (cfgKey, pkgJson) => {
        if (cfgKey.match(/^package|package-upgrade$/)) {
            return semver.validRange(pkgJson?.engines?.node)
        }
        return semver.validRange(config.desiredNodeVersion)
    }

    const pkgRaw = fs.readFileSync(inSrcDir('package.json'))
    const pkgJson = JSON.parse(pkgRaw)

    let desiredVersion = getDesiredVersion(config.desiredNodeVersion, pkgJson)
    if (!desiredVersion)
        throw new Error('Desired Node version can not be determined.')

    const nodeIndexText = await (
        await fetch('https://nodejs.org/dist/index.json')
    ).text()
    const nodeIndexJson = JSON.parse(nodeIndexText)
    const nodeAvailableVersions = nodeIndexJson.map((x) => x.version)

    let chosenVersion = semver.maxSatisfying(
        nodeAvailableVersions,
        desiredVersion
    )
    if (config.desiredNodeVersion == 'package-upgrade') {
        chosenVersion = semver.maxSatisfying(
            nodeAvailableVersions,
            `^${semver.major(chosenVersion)}`
        )
    }
    console.log(`Node version: ${chosenVersion}`)

    // replace engine property from package.json to enable
    // the usage of newer Node versions
    pkgJson.engines.node = chosenVersion
    fs.writeFileSync(
        inSrcDir('package.json'),
        JSON.stringify(pkgJson, undefined, 2)
    )

    // write Node info to metadata file for later usage during the build process.
    const metadata = {
        NODE_VERSION: chosenVersion,
        NODE_URL: `https://nodejs.org/dist/${chosenVersion}/node-${chosenVersion}-linux-x64.tar.gz`,
    }
    const metadataString = Object.entries(metadata).reduce(
        (acc, [k, v]) => (acc += `${k}=${v}\n`),
        ''
    )
    fs.appendFileSync(inSrcDir('.metadata'), metadataString)
}

async function taskBuild(config) {
    console.log(`docker-compose build eln`)
}

async function taskApplyPatches(config) {
    for (let fileName of globbySync(`${config.fixDir}/*.patch`).sort()) {
        const fileContents = fs
            .readFileSync(fileName, { encoding: 'utf-8' })
            .split('\n')

        let fixedBy = []
        if (fileContents[0].startsWith('FIXEDBY ')) {
            fixedBy = fileContents[0].split(' ').slice(1)
        }

        let applyFix = true
        for (const fix of fixedBy) {
            const { stdout } = await bash(
                `git merge-base ${fix} --is-ancestor HEAD && echo "discard" || echo "apply"`, // needs git 1.8.0
                {
                    cwd: config.srcDir,
                }
            )
            if (stdout.includes('discard')) {
                applyFix = false
                break
            }
        }

        if (applyFix) {
            await bash(`git apply ${fileName} || true`, { cwd: config.srcDir })
            console.log(`Patch [${path.basename(fileName)}] applied.`)
        } else {
            // patch is unnecessary in this release
            console.log(`Patch [${path.basename(fileName)}] is obsolete.`)
        }
    }
}

async function taskConfigureChemotion(config) {
    // configure the ELN:
    // step 1: basic cleaning
    //      - copy example files, remove .example extension
    // step 2: configure context specifics (nothing here yet...)

    // this is the data structure for manualFiles later in the code:
    // store what's in 'fc' as file named 'fn'.
    const databaseConfig = {
        fn: 'config/database.yml', // file newname
        fc: (_) =>
            yamlDump({
                // file content
                production: {
                    adapter: 'postgresql',
                    encoding: 'unicode',
                    database: 'chemotion',
                    pool: 5,
                    username: 'postgres',
                    password: 'postgres',
                    host: 'db',
                    port: '5432',
                },
            }),
    }

    const secretsConfig = {
        fn: 'config/secrets.yml',
        fc: (_) =>
            yamlDump({
                production: {
                    secret_key_base: "<%= ENV['SECRET_KEY_BASE'] %>",
                },
            }),
    }

    // As of now, we use static files in fixes directory to configure database.yml and secrets.yml.
    // Feel free to enable this, if you need context specific configurations
    const manualConfigs = [
        // example: this array could be
        // databaseConfig,
        // secretsConfig
    ]

    // rename all example files
    const exampleFiles = ['**/*.example']
    for (const fileName of globbySync(exampleFiles.map(inSrcDir), {
        dot: true,
        onlyFiles: true,
    })) {
        const newFileName = fileName.replace(/\.example$/, '')
        fs.renameSync(fileName, newFileName)
    }

    // write files manually
    for (const manualConfig of manualConfigs) {
        const fileName = inSrcDir(manualConfig.fn)
        const srcFileContent = fs.readFileSync(fileName, { encoding: 'utf-8' })
        const dstFileContent = manualConfig.fc(srcFileContent)
        fs.writeFileSync(fileName, dstFileContent, { encoding: 'utf-8' })
        console.log(`Wrote [${fileName}].`)
    }
}

async function taskEmbed(config) {
    const includePath = path.join(config.fixDir, 'include')
    if (fs.existsSync(includePath)) {
        console.log(`Including [${includePath}] as is.`)
        fs.cpSync(includePath, config.srcDir, { recursive: true, force: true })
    }
}

async function taskCreateSeedArchive(config) {
    const tarfile = path.join(config.srcDir, 'seed.tar.gz')
    // TODO: fixme -- if there are no files in seedDir, * is nothing
    // which results in tar to fail. so we put the '|| true' part.
    // That seems to be bad practice. Should be fixed.
    await bash(`cd ${config.seedDir} && tar cfz ${tarfile} * || true`)
}

async function taskLinkStorages(config) {
    const uniq = (a) => {
        return Array.from(new Set(a))
    }

    fs.rmSync(config.seedDir, { recursive: true, force: true })
    fs.mkdirSync(config.seedDir, { recursive: true })

    volumeInfo.pullin = uniq(volumeInfo.pullin)
    volumeInfo.expose = uniq(volumeInfo.expose)
    volumeInfo.ephemeral = uniq(volumeInfo.ephemeral)

    async function exposeTo(list, exposePath, srcPath, skipList) {
        const isSymLink = (fn) => {
            try {
                const stats = fs.lstatSync(fn)
                if (stats.isSymbolicLink()) {
                    return fs.readlinkSync(fn)
                } else {
                    return false
                }
            } catch {
                return false
            }
        }

        for (const p of list.sort(dirSort)) {
            const resolvedLinkName = path.resolve(
                path.isAbsolute(p) ? p : path.join(srcPath, p)
            )
            const resolvedLinkTarget = path.resolve(path.join(exposePath, p))
            const resolvedLinkTargetPath = path.dirname(resolvedLinkName)

            let skip = false
            for (let saw of skipList) {
                if (resolvedLinkName.startsWith(saw + path.sep)) {
                    console.log(
                        `WARNING: [${resolvedLinkName}] is already exposed by me via [${saw}]`
                    )
                    skip = true
                    break
                }
            }

            if (skip) continue

            const elementsOfLinkName = resolvedLinkName.split(path.sep)

            // boundaries here are a bit tricky:
            // if the path has 8 elements, since we start with an absolute path, we have an array length
            // of 9 (first element in array being empty). In addition, we use the slice function, which excludes
            // the end from the list (i.e. slice(1, 3) returns 2 elements, element 1 and element 2). Thus i must
            // not always be bigger than 1.
            // The upper boundary is length - 1, since we do not want to check for the element itself.
            // i.e. if we have /a/b/c/d, we want to check /a/b/c at most.
            for (let i = elementsOfLinkName.length - 1; i > 1; i--) {
                const toCheck = elementsOfLinkName.slice(0, i).join(path.sep)
                if (toCheck === config.srcDir) break

                const trgt = isSymLink(toCheck)
                if (trgt) {
                    console.log(
                        `WARNING: [${resolvedLinkName}] is already exposed via [${toCheck} -> ${trgt}]`
                    )
                    skip = true
                    break
                } else {
                    // it's not a symlink
                }
            }
            if (skip) continue

            // console.log(`${resolvedLinkName} => ${resolvedLinkTarget} [${resolvedLinkTargetPath}]`)

            const seedTarget = path.join(config.seedDir, exposePath, p)
            const seedPath = path.dirname(seedTarget)
            fs.mkdirSync(seedPath, { recursive: true })
            if (fs.existsSync(resolvedLinkName)) {
                // if folder exists, move it to seed folder
                fs.renameSync(resolvedLinkName, seedTarget)
                console.log(
                    `Added [${resolvedLinkName}] to seed as [${seedTarget}]`
                )
            } else {
                // if it does not, create an empty one in the seed
                fs.mkdirSync(seedTarget, { recursive: true })
                console.log(
                    `Created an empty seed folder for [${resolvedLinkName}] as [${seedTarget}]`
                )
            }

            fs.mkdirSync(resolvedLinkTargetPath, { recursive: true })
            await bash(`ln -s ${resolvedLinkTarget} ${resolvedLinkName}`)
            skipList.push(resolvedLinkName)
        }
        return skipList
    }

    // await exposeTo(volumeInfo.ephemeral, config.ephemeralDir, config.srcDir, config.exposedPaths, false)
    // await exposeTo(volumeInfo.expose, config.exposeDir, config.srcDir, config.exposedPaths)
}

async function taskHandleUploadsAndImages(config) {
    fs.mkdirSync(path.join(config.dataDir, 'public', 'images'), {
        recursive: true,
    })
    fs.mkdirSync(path.join(config.dataDir, 'uploads'), { recursive: true })
    await bash(
        `mv "${path.join(config.srcDir, 'public', 'images')}" "${path.join(
            config.dataDir,
            'public'
        )}"`
    )
    await bash(
        `ln -s /chemotion/data/public/images "${path.join(
            config.srcDir,
            'public',
            'images'
        )}"`
    )
    await bash(
        `ln -s /chemotion/data/uploads "${path.join(config.srcDir, 'uploads')}"`
    )
    console.log('Moved [public/images] and [uploads] to data.')
}

// main function.
;(async () => {
    await taskCheckoutRepo(config)
    await taskDetermineNodeVersion(config)
    await taskApplyPatches(config)
    await taskConfigureChemotion(config)
    await taskCleanRepo(config)
    await taskEmbed(config)
    // await taskLinkStorages(config)
    await taskHandleUploadsAndImages(config)
    // await taskCreateSeedArchive(config)
    // await taskBuild(config)
})()
