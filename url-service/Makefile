up :
	@echo "==========Starting up the application=========="
	docker-compose -f docker-compose.yml up -d --build --remove-orphans
down :
	@echo "==========Stopping the application=========="
	docker-compose -f docker-compose.yml down --remove-orphans

build :
	@echo "==========Building the application=========="
	docker build -t url-service .
	@echo "==========Application built=========="
run:
	@echo "==========Running the application=========="
	docker run -p 8080:8080 --name url-container url-service
	@echo "==========Application is running=========="	
stop :
	@echo "==========Stopping the application=========="
	docker stop url-container
	@echo "==========Application stopped=========="
test : 
	@echo "==========Running tests=========="
	go test -coverprofile=coverage.out ./...
test-html :
	@echo "==========Generating test coverage report=========="
	go tool cover -html=coverage.out -o coverage.html

