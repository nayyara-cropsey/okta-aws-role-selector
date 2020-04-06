build:
	@echo "Build"
	go build -o main .

run:
	@echo "Build"
	./main -c config.yaml

docker_build:
	@echo "Docker Build"
	docker build -t nayyarasamuel7/okta-aws-role-selector:local .

docker_run:
	@echo "Docker Run"
	docker run -p 8080:8080 -v $(PWD):/root/config d4244b29d434  -c config/config.yaml

.PHONY: build
