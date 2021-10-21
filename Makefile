build:
	# build backend
	docker build -t say-backend .
	# build client
	cd client && go build -o client
push:
	# add the image we built to micro8s cluster
	# so that kubectl can use it for creating deployment
	# defined in backend/kubernetes.yml
	docker save say-backend > say-backend.tar
	microk8s ctr image import say-backend.tar
	rm say-backend.tar