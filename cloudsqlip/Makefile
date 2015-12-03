
.PHONY: compile build_image push
VERSION := 0.3

compile:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

build_image:
	docker build -t jordic/cloudsqlip:$(VERSION) .

push:
	docker push jordic/cloudsqlip:$(VERSION) 
