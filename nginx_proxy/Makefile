
.PHONY: compile build_image push

VERSION := v1

compile:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

build_image:
	docker build -t jordic/nginxsp:$(VERSION) .

push:
	docker push jordic/nginxsp:$(VERSION)

