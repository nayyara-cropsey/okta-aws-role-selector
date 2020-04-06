build:
	@echo "Build"
	go build -o okta-aws-role-selector .

run:
	@echo "Build"
	./okta-aws-role-selector -c config.yaml

.PHONY: build
