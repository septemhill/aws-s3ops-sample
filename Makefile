BINNAME=s3ops

all: build start-localstack-docker run-sample

build:
	@echo "running go build"
	@go build -o $(BINNAME)

start-localstack:
	@echo "running localstack docker"
	@docker-compose up -d

stop-localstack:
	@echo "stop localstack docker"
	@docker-compose down

run-sample:
	./$(BINNAME)

clean:
	rm -rf $(BINNAME)