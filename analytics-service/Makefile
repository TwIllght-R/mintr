APP_NAME = url-service
IMAGE_NAME = $(APP_NAME)-dev
DOCKER_DEV_FILE = Dockerfile.dev
CONTAINER_NAME = $(APP_NAME)-dev-container
PORT = 8080

dev-build:
	docker build -f $(DOCKER_DEV_FILE) -t $(IMAGE_NAME) .
dev-up:
	docker run --rm -it \
		--name $(CONTAINER_NAME) \
		-p $(PORT):$(PORT) \
		-v $$PWD:/app \
		-w /app \
		$(IMAGE_NAME)

dev-shell:
	docker run -it \
		--name $(CONTAINER_NAME) \
		-p $(PORT):$(PORT) \
		-v $$PWD:/app \
		-w /app \
		$(IMAGE_NAME) sh

dev-clean:
	docker rmi $(IMAGE_NAME) || true
test : 
	@echo "==========Running tests=========="
	go test -coverprofile=coverage.out ./...
test-html :
	@echo "==========Generating test coverage report=========="
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run