#!/bin/bash
# RUNAS: root
# unblocks pubchem

# do not use -i here. Docker prevents it from working...
sed '/pubchem.ncbi.nlm.nih.gov/d' /etc/hosts > /etc/hosts
