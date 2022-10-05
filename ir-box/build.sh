#!/bin/bash

[[ ! -e src ]] && git clone https://github.com/ComPlat/chem-dl-ir.git src

docker build -t chem-ir:latest .
