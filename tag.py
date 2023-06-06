import os
import sys
from dataclasses import dataclass
from functools import cache
import subprocess

import argparse


class Globals:
    services = [
        "base",
        "converter",
        "ketchersvc",
        "ketchersvc-sc",
        "spectra",
        "msconvert",
        "eln",
    ]


class RepoNameTag:
    def __init__(self, repo, name, tag):
        self.repo = repo
        self.name = name
        self.tag = tag

    def toStr(self):
        return self.__str__()

    def __str__(self):
        if "/" in self.repo:
            return f"{self.repo}:{self.name}-{self.tag}"
        else:
            return f"{self.repo}/{self.name}:{self.tag}"


class Docker:
    ALLOWED_NAME_CHARACTERS = "abcdefghijklmnopqrstuvwxyz0123456789"
    ALLOWED_TAG_CHARACTERS = (
        "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789._-"
    )
    TAG_SEPARATORS = ".-"
    NAME_SEPARATORS = [".", "_", "__", "-", "--", "---"]

    @staticmethod
    def validTag(s: str):
        return (
            all([x in Docker.ALLOWED_TAG_CHARACTERS for x in s])
            and all([not s.startswith(y) for y in Docker.TAG_SEPARATORS])
            and len(s) > 0
            and len(s) < 128
        )

    @staticmethod
    def validName(s: str):
        return (
            all([x in Docker.ALLOWED_NAME_CHARACTERS for x in s])
            and all(
                [
                    not s.startswith(y) and not s.endswith(y)
                    for y in Docker.NAME_SEPARATORS
                ]
            )
            and len(s) > 0
            and len(s) < 128
        )

    @staticmethod
    def validImageName(s: str):
        if ":" in s:
            name, tag = s.split(":")
            return Docker.validName(name) and Docker.validTag(tag)
        else:
            return Docker.validName(s)

    @staticmethod
    @cache
    def allImages():
        def splitLSLine(line):
            repoNSvc, id, tag, *_ = line.split("|")
            repo, svc = repoNSvc.split("/")
            imageInfo = ImageInfo(repo, svc, id, tag)
            return imageInfo

        proc = subprocess.run(
            ["docker", "image", "ls", "--format", "{{.Repository}}|{{.ID}}|{{.Tag}}"],
            stdout=subprocess.PIPE,
        )
        if proc.returncode != 0:
            print("Error: docker image ls failed")
            sys.exit(1)

        imageList = [
            splitLSLine(x.strip())
            for x in proc.stdout.decode("utf-8").splitlines()
            if x
        ]

        return imageList

    @staticmethod
    def addTag(src: RepoNameTag, dst: RepoNameTag):
        proc = subprocess.run(
            ["docker", "image", "tag", src.toStr(), dst.toStr()],
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
        )
        Docker.allImages.cache_clear()
        if proc.returncode != 0:
            print(f"Error: failed to tag image.")
            print(
                "  "
                + proc.stderr.decode("utf-8")
                .replace("Error response from daemon:", "")
                .strip()
                .replace("\n", "\n  ")
            )
            print()
            return False
        else:
            print(f"Tagged {src.toStr()}\t\t\t\t as {dst.toStr()}")
            return True

    @staticmethod
    def pushImage(what: RepoNameTag):
        proc = subprocess.run(
            ["docker", "image", "push", what.toStr()],
        )
        if proc.returncode != 0:
            print(f"Error: failed to push image {what.toStr()}")
            return False
        else:
            print(f"Pushed {what.toStr()}")
            return True


@dataclass
class Args:
    sourceTag: str
    newTag: str
    sourceRepo: str
    targetRepo: str
    push: bool
    latest: bool
    sourceLatest: bool

    def printConfig(self):
        print("Config:")
        print("Source Tag  :", self.sourceTag)
        print("New Tag     :", self.newTag)
        print("Source Repo :", self.sourceRepo)
        print("Target Repo :", self.targetRepo)
        print("Push        :", self.push)
        print("Apply Latest:", self.latest)
        print("--\n")


@dataclass
class ImageInfo:
    repo: str
    svc: str
    id: str
    tag: str


parser = argparse.ArgumentParser()
parser.add_argument("sourceTag", help="Tag shared for all services.", type=str)
parser.add_argument(
    "newTag",
    help="New tag to attach to all services. Use '-' to use sourceTag as targetTag.",
    type=str,
    default="-",
)
parser.add_argument(
    "--sourceRepo",
    help="source repo (default: chemotion-build)",
    default="chemotion-build",
)
parser.add_argument(
    "--targetRepo",
    help="target repo (default: ptrxyz/chemotion)",
    default="ptrxyz/chemotion",
)
parser.add_argument("--push", help="Also push to Docker Hub.", action="store_true")
parser.add_argument("--latest", help="Also tag as latest", action="store_true")
parser.add_argument(
    "--sourceLatest", help="Also tag as latest in source repo", action="store_true"
)

# if targetRepo does contain "/", the new image name will be
# `<targetRepo>:<service>-<tag>`, otherwise `<targetRepo>/<service>:<tag>`

args = Args(**dict(parser.parse_args()._get_kwargs()))
if args.newTag == "-":
    args.newTag = args.sourceTag

if (not Docker.validTag(args.sourceTag)) or (not Docker.validTag(args.newTag)):
    print("Invalid tag")
    sys.exit(1)


class Helpers:
    @staticmethod
    def grep(cmp, lines):
        return [x for x in lines if cmp(x)]

    @staticmethod
    def checkDockerImageExists(repo, svc, tag):
        allDockerImages = Docker.allImages()
        return (
            len(
                Helpers.grep(
                    lambda x: x.repo == repo and x.svc == svc and x.tag == tag,
                    allDockerImages,
                )
            )
            > 0
        )

    @staticmethod
    def checkAllServiceImagesExist(repo, tag):
        ret = dict(
            [
                (
                    RepoNameTag(repo, svc, tag),
                    Helpers.checkDockerImageExists(repo, svc, tag),
                )
                for svc in Globals.services
            ]
        )

        if not all(ret.values()):
            print("Error: Not all images exist")
            for k, v in ret.items():
                if not v:
                    print(f"  Missing: {k}")
            return False
        else:
            return True


# if not Helpers.checkAllServiceImagesExist(args.sourceRepo, args.sourceTag):
#     sys.exit(1)

for svc in Globals.services:
    Docker.addTag(
        src=RepoNameTag(args.sourceRepo, svc, args.sourceTag),
        dst=RepoNameTag(args.targetRepo, svc, args.newTag),
    )

    if args.latest:
        Docker.addTag(
            src=RepoNameTag(args.sourceRepo, svc, args.sourceTag),
            dst=RepoNameTag(args.targetRepo, svc, "latest"),
        )

    if args.sourceLatest:
        Docker.addTag(
            src=RepoNameTag(args.sourceRepo, svc, args.sourceTag),
            dst=RepoNameTag(args.sourceRepo, svc, "latest"),
        )

print("Tagged all images.")

if args.push:
    print("Pushing images...")
    for svc in Globals.services:
        Docker.pushImage(RepoNameTag(args.targetRepo, svc, args.newTag))
        if args.latest:
            Docker.pushImage(RepoNameTag(args.targetRepo, svc, "latest"))
