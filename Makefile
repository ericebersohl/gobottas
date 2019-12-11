.PHONY: images

images:
	$(eval TAG=`git tag | tail -1`)
	docker build . -t gobottas:$(TAG)