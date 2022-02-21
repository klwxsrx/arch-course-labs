.PHONY: clean docker-image/%

all: clean bin/auth bin/userprofile docker-image/auth docker-image/userprofile

clean:
	rm -rf bin/*

bin/%:
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o ./bin/$(notdir $@) ./cmd/$(notdir $@)

docker-image/%:
	docker build -f docker/$(notdir $@)/Dockerfile -t klwxsrx/arch-course-labs-05-$(notdir $@)-service .