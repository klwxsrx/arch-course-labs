.PHONY: clean docker/auth

all: clean bin/auth docker/auth

clean:
	rm -rf bin/*

bin/%:
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o ./bin/$(notdir $@) ./cmd/$(notdir $@)

docker/auth:
	docker build -f docker/auth/Dockerfile -t klwxsrx/arch-course-labs-05-auth-service .