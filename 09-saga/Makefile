.PHONY: clean docker-image/%

all: clean \
	bin/order bin/payment bin/delivery bin/warehouse \
	docker-image/order docker-image/payment	docker-image/delivery docker-image/warehouse

clean:
	rm -rf bin/*

bin/%:
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o ./bin/$(notdir $@) ./cmd/$(notdir $@)

docker-image/%:
	docker build -f docker/$(notdir $@)/Dockerfile -t klwxsrx/arch-course-labs-09-$(notdir $@)-service .