import yaml from 'yaml'

class RawTag {
    constructor(v) {
        this.raw = v
    }

    toString() {
        return `${this.raw}`
    }
}

const regexp = {
    identify: (value) => {
        return value instanceof RawTag === true
    },
    tag: '!ruby/regexp',
    resolve(doc, cst) {
        return new RawTag(cst.strValue)
    },
}

yaml.defaultOptions.customTags = [regexp]

function yamlLoad(str) {
    return yaml.parse(str)
}

function yamlDump(yml) {
    return yaml.stringify(yml, {
        simpleKeys: true,
        defaultKeyType: 'BLOCK_LITERAL',
        keepSourceTokens: true,
        schema: 'yaml-1.1',
    })
}

export { yamlLoad, yamlDump }
