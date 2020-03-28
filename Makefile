.PHONY: images

build:
	$(eval TAG=`git tag | tail -1`)
	docker build . -t gobottas:$(TAG) -t gobottas:latest

deploy:
	docker run -d gobottas:latest