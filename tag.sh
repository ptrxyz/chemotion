#!/bin/bash

action=${1:-"tag"}
tag=${2:-}
services=("base" "converter" "ketchersvc" "ketchersvc-sc" "spectra" "msconvert" "eln")

case "$action" in
    "push")
        echo "Pushing containers:"
        docker image ls | grep -e "ptrxyz/internal:" -e "-dev"
        for i in "${services[@]}"; do
            echo -e "\n+ Pushing [$i]..."
            docker image ls "ptrxyz/internal:${i}-dev" | grep -q "${i}-dev" || {
                echo "Image not found: ptrxyz/internal:${i}-dev"
            }
            systemd-inhibit docker push "ptrxyz/internal:${i}-dev"
        done
        ;;
    "tag")
        echo "Tagging containers:"
        for i in "${services[@]}"; do
            echo "+ Tagging [$i]..."
            # echo "  [chemotion-build/${i}:latest]  ->  [ptrxyz/internal:${i}-dev]"
            docker tag "chemotion-build/${i}:latest" "ptrxyz/internal:${i}-dev"
        done
        echo -e "\nTagged containers:"
        docker image ls | grep -e "ptrxyz/internal:" -e "-dev"
        ;;
    "push-prod")
        if [[ -z "$tag" ]]; then
            echo "Missing tag."
            exit 1
        fi
        echo "Pushing containers:"
        docker image ls | grep -e "ptrxyz/chemotion:" -e "-${tag}"
        for i in "${services[@]}"; do
            echo -e "\n+ Pushing [$i]..."
            docker image ls "ptrxyz/chemotion:${i}-${tag}" | grep -q "${i}-${tag}" || {
                echo "Image not found: ptrxyz/chemotion:${i}-${tag}"
            }
            systemd-inhibit docker push "ptrxyz/chemotion:${i}-${tag}"
            systemd-inhibit docker push "ptrxyz/chemotion:${i}-latest"
        done
        ;;
    "tag-prod")
        if [[ -z "$tag" ]]; then
            echo "Missing tag."
            exit 1
        fi
        echo "Tagging containers with tag [${tag}]:"
        for i in "${services[@]}"; do
            echo "+ Tagging [$i]..."
            # echo "  [chemotion-build/${i}:latest]  ->  [ptrxyz/internal:${i}-dev]"
            docker tag "chemotion-build/${i}:latest" "ptrxyz/chemotion:${i}-latest"
            docker tag "chemotion-build/${i}:latest" "ptrxyz/chemotion:${i}-${tag}"
        done
        echo -e "\nTagged containers:"
        docker image ls | grep -e "ptrxyz/chemotion"
        ;;
    *)
        echo "Unknown action: $action"
        exit 1
        ;;
esac

