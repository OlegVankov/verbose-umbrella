SERVER_PORT=8080
AGENT_PATH=./agent
SERVER_PATH=./server
ADDRESS="127.0.0.1:${SERVER_PORT}"

.DEFAULT_GOAL := default

.PHONY: clean
clean:
	@rm -rf ${AGENT_PATH} ${SERVER_PATH}

.PHONY: build
build: clean
	@go build -o ${AGENT_PATH} ./cmd/agent/...
	@go build -o ${SERVER_PATH} ./cmd/server/...

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

.PHONY: default
default: clean build mtest1 mtest2 mtest3 mtest4 mtest5 mtest6 mtest7 mtest8
