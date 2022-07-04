.DEFAULT_GOAL := prepare
.PHONY: prepare

prepare:
	chmod +x scripts/fake-docker.py && chmod +x scripts/mscrunner.py
	rm -rf src && git clone https://github.com/ComPlat/chem-spectra-app.git src && cd src && git checkout $${SPECTRA_VERSION-development} && rm -rf .git
	[[ -f miniconda.sh ]] || (curl -L -o miniconda.sh https://repo.anaconda.com/miniconda/Miniconda3-latest-Linux-x86_64.sh && chmod +x miniconda.sh)

clean:
	rm -rf src
	rm -rf miniconda.sh
