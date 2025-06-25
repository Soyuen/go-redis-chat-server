IMAGE_NAME=go-redis-chat-server
REDIS_CONTAINER_NAME=redis-server
APP_CONTAINER_NAME=go-redis-chat
DOCKER_NETWORK=go-redis-chat-server_default

# 取消外部環境中敏感 build args 影響，改用空值覆蓋
BUILD_ARGS := --build-arg API_KEY= --build-arg OTHER_SECRET=

.PHONY: build
build:
	docker build --no-cache $(BUILD_ARGS) -t $(IMAGE_NAME) .

.PHONY: start-redis
start-redis:
	@if [ -z "$$(docker ps -q -f name=$(REDIS_CONTAINER_NAME))" ]; then \
		echo "Starting Redis container..."; \
		docker run -d --name $(REDIS_CONTAINER_NAME) -p 6379:6379 --network $(DOCKER_NETWORK) redis:7; \
	else \
		echo "Redis container already running."; \
	fi

.PHONY: run
run: start-redis build
	@if [ -z "$$(docker ps -q -f name=$(APP_CONTAINER_NAME))" ]; then \
		echo "Starting app container..."; \
		docker run -d --name $(APP_CONTAINER_NAME) -p 8080:8080 --env-file .env --network $(DOCKER_NETWORK) $(IMAGE_NAME); \
	else \
		echo "App container already running."; \
	fi

.PHONY: stop
stop:
	docker stop $(APP_CONTAINER_NAME) $(REDIS_CONTAINER_NAME) || true
	docker rm $(APP_CONTAINER_NAME) $(REDIS_CONTAINER_NAME) || true

.PHONY: logs
logs:
	docker logs -f $(APP_CONTAINER_NAME)
