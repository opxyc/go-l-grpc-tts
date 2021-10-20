build:
	GOOS=linux go build -o app
	docker build -t say-backend .
	rm -f app
push:
	# add it to micro8s cluster
	docker save say-backend > say-backend.tar
	microk8s ctr image import say-backend.tar
	microk8s ctr images ls | grep say-latest