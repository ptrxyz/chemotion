build:
	cd ./eln2/ && docker build -f build.dockerfile -t chmtn:eln
	cd ./converter/ && docker build -f build.dockerfile -t chmtn:converter
	cd ./base/ && docker build -f build.dockerfile -t chmtn:base
	cd ./spectra/ && docker build -f build.dockerfile -t chmtn:spectra
	cd ./ketchersvc/ && docker build -f build.dockerfile -t chmtn:ketcher
