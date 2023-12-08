DOCKER_VERSION=0.1.1

docker.job:job.build
	docker build --no-cache --build-arg PROJECT_NAME=go-kbs-job-example -f Dockerfile.job.container -t "liburdi/go-k8s-job-example:$(DOCKER_VERSION)" ./

docker.operator:operator.build
	docker build --no-cache --build-arg PROJECT_NAME=go-kbs-operator -f Dockerfile.operator.container -t "liburdi/go-k8s-operator:$(DOCKER_VERSION)" ./

docker.build:docker.job docker.operator

job.build:
	GOOS=linux GOARCH=amd64 go build -o release/go-k8s-job-example ./container/main.go

operator.build:
	GOOS=linux GOARCH=amd64 go build -o release/go-k8s-operator ./cmd/main.go

docker.push:docker.build
	docker push liburdi/go-k8s-job-example:$(DOCKER_VERSION)
	docker push liburdi/go-k8s-operator:$(DOCKER_VERSION)


docker.clean:
	rm -rf release


run.operator:
	go run  ./cmd/main.go --image=liburdi/go-k8s-job-example:0.0.5 --name=job-example