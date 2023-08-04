SERVER_PORT=8080
AGENT_PATH=./bin/agent
SERVER_PATH=./bin/server
TEMP_FILE=/tmp/metric-server-storage.json
DSN='postgres://postgres:postgres@localhost:5432/praktikum?sslmode=disable'
ADDRESS="127.0.0.1:${SERVER_PORT}"

.DEFAULT_GOAL := default

.PHONY: clean
clean:
	@rm -rf ${AGENT_PATH} ${SERVER_PATH}

.PHONY: build
build: clean
	@go build -o ${AGENT_PATH} ./cmd/agent/...
	@go build -o ${SERVER_PATH} ./cmd/server/...

.PHONY: deps
deps:
	@go get github.com/jmoiron/sqlx
	@go get github.com/jackc/pgx
	@go get github.com/go-chi/chi/v5

.PHONY: test
test:
	@go test -v -cover ./...

.PHONY: mtest1
mtest1: build
	@metricstest -test.v -test.run="^TestIteration1$$" \
            -binary-path=${SERVER_PATH}

.PHONY: mtest2
mtest2: build
	@metricstest -test.v -test.run="^TestIteration2[AB]*$$" \
            -agent-binary-path=${AGENT_PATH} \
            -source-path=.

.PHONY: mtest3
mtest3: build
	@metricstest -test.v -test.run="^TestIteration3[AB]*$$" \
            -agent-binary-path=${AGENT_PATH} \
            -binary-path=${SERVER_PATH} \
            -source-path=.

.PHONY: mtest4
mtest4: build
	@metricstest -test.v -test.run="^TestIteration4$$" \
		-agent-binary-path=${AGENT_PATH} \
		-binary-path=${SERVER_PATH} \
		-server-port=${SERVER_PORT} \
		-source-path=.

.PHONY: mtest5
mtest5: build
	@metricstest -test.v -test.run="^TestIteration5$$" \
		-agent-binary-path=${AGENT_PATH} \
		-binary-path=${SERVER_PATH} \
		-server-port=${SERVER_PORT} \
		-source-path=.

.PHONY: mtest6
mtest6: build
	@metricstest -test.v -test.run="^TestIteration6$$" \
		-agent-binary-path=${AGENT_PATH} \
		-binary-path=${SERVER_PATH} \
		-server-port=${SERVER_PORT} \
		-source-path=.

.PHONY:	mtest7
mtest7: build
	@metricstest -test.v -test.run="^TestIteration7$$" \
		-agent-binary-path=${AGENT_PATH} \
		-binary-path=${SERVER_PATH} \
		-server-port=${SERVER_PORT} \
		-source-path=.

.PHONY: mtest8
mtest8: build
	@metricstest -test.v -test.run="^TestIteration8$$" \
		-agent-binary-path=${AGENT_PATH} \
		-binary-path=${SERVER_PATH} \
		-server-port=${SERVER_PORT} \
		-source-path=.

.PHONY: mtest9
mtest9: build
	@metricstest -test.v -test.run="^TestIteration9$$" \
		-agent-binary-path=${AGENT_PATH} \
		-binary-path=${SERVER_PATH} \
		-file-storage-path=${TEMP_FILE} \
		-server-port=${SERVER_PORT} \
		-source-path=.

.PHONY: mtest10
mtest10: build
	@metricstest -test.v -test.run="^TestIteration10[AB]$$" \
            -agent-binary-path=${AGENT_PATH} \
            -binary-path=${SERVER_PATH} \
            -database-dsn=${DSN} \
            -server-port=${SERVER_PORT} \
            -source-path=.

.PHONY: mtest11
mtest11: build
	@metricstest -test.v -test.run="^TestIteration11$$" \
            -agent-binary-path=${AGENT_PATH} \
            -binary-path=${SERVER_PATH} \
            -database-dsn=${DSN} \
            -server-port=${SERVER_PORT} \
            -source-path=.

.PHONY: mtest12
mtest12: build
	@metricstest -test.v -test.run="^TestIteration12$$" \
            -agent-binary-path=${AGENT_PATH} \
            -binary-path=${SERVER_PATH} \
            -database-dsn=${DSN} \
            -server-port=${SERVER_PORT} \
            -source-path=.

.PHONY: mtest13
mtest13: build
	@metricstest -test.v -test.run="^TestIteration13$$" \
            -agent-binary-path=${AGENT_PATH} \
            -binary-path=${SERVER_PATH} \
            -database-dsn=${DSN} \
            -server-port=${SERVER_PORT} \
            -source-path=.

.PHONY: default
default: clean build \
		 mtest1 mtest2 mtest3 mtest4 mtest5 mtest6 mtest7 mtest8 mtest9 mtest10 \
		 mtest11 mtest12 mtest13
