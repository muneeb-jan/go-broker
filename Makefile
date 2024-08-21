.PHONY: build run image deploy destroy help

all: image deploy

build:
	go build cmd/go-broker/main.go

run:
	go run cmd/go-broker/main.go

image:
	docker build -t ${IMAGE_NAME} .

deploy:
	docker-compose -f compose.yaml up

destroy:
	docker-compose -f compose.yaml up

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build  	Package and Install Dependencies"
	@echo "  run		Run the Application"
	@echo "  image		Create Docker Image"
	@echo "  deploy    	Deploy in Docker Container"
	@echo "  destroy    	Destroy Docker Containers"
	@echo "  all	   	Run Install and Deploy"
	@echo "  help      	Display this help message"
