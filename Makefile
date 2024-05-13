all:uibuild
	go build ./
uibuild:
	cd ui && npm run build