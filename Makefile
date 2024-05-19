run:build
	./go_books

build:uibuild
	go build ./
uibuild:
	cd ui && npm run build

