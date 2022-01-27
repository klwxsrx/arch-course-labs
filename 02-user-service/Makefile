APP_VERSION=1.0.0

.PHONY: clean build docker-image

all: clean build docker-image

clean:
	rm -rf bin/*

build:
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o ./bin/user ./cmd/user

docker-image:
	docker build -f docker/Dockerfile -t klwxsrx/arch-course-labs-userservice:$(APP_VERSION) .