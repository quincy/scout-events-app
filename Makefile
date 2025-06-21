APP_NAME = scout-events-app
SRC_DIR = ./src
MAIN_FILE = main.go
BIN_DIR = ./bin
BIN_FILE = $(BIN_DIR)/$(APP_NAME)

.PHONY: all build fmt test test-unit test-integration clean docker-up docker-down app-run app-stop

UNIT_TEST_PKGS := $(shell find ${SRC_DIR} -name '*_test.go' -not -path '${SRC_DIR}/integration-tests/*' -exec dirname {} \; | sort -u)

all: build test

build: fmt
	go build -o ${BIN_FILE} ${MAIN_FILE}

fmt:
	go fmt ${MAIN_FILE}
	go fmt ${SRC_DIR}/...

test: test-unit test-integration

test-unit:
	go test -v $(UNIT_TEST_PKGS)

docker-up:
	docker compose up -d

docker-down:
	docker compose down

app-run:
	@echo "Starting application..."
	@nohup go run ${MAIN_FILE} > app.log 2>&1 & echo $$! > app.pid
	@sleep 2

app-stop:
	@if [ -f app.pid ]; then \
		kill `cat app.pid` && rm app.pid; \
		 echo "Application stopped."; \
	else \
		echo "No app.pid file found. Application may not be running."; \
	fi

# Integration tests: start docker, app, run tests, then stop everything

test-integration: docker-up app-run
	go test -v ${SRC_DIR}/integration-tests/...
	$(MAKE) app-stop
	rm -f app.pid app.log
	$(MAKE) docker-down

clean:
	rm -rf ${BIN_FILE} app.log app.pid
