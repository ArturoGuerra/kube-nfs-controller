.PHONY: all build clean

GOBUILD = go build

all: clean build


clean:
	rm -rf bin


build: clean
	COOS=linux $(GOBUILD) -o bin/controller ./cmd


docker-build: build
	test $(DOCKERREPO)
	docker build . -t $(DOCKERREPO)

docker-push: docker-build
	test $(DOCKERREPO)
	docker push $(DOCKERREPO)

docker: build docker-build docker-push clean
