apps = 'gocron-agent' 'gocron-server'

ORG ?= caas4

.PHONY: build
build: fmt checkfmt
	for app in $(apps) ;\
	do \
		CGO_ENABLED=0 go build -o dist/$$app -a -ldflags "-w -s" ./cmd/$$app;\
	done

.PHONY: fmt
fmt:
	gofmt -s -w ./

.PHONY: checkfmt
checkfmt:
	@echo checking gofmt...
	@res=$$(gofmt -d -e -s $$(find . -type d \( -path ./src/vendor -o -path ./tests \) -prune -o -name '*.go' -print)); \
	if [ -n "$${res}" ]; then \
		echo checking gofmt fail... ; \
		echo "$${res}"; \
		exit 1; \
	fi

.PHONY: swagger
swagger:
	go run tools/doc-gen/main.go --output=swagger-ui/swagger.json

.PHONY: swagger-server
swagger-server:
	go run swagger-ui/swagger.go

.PHONY: proto
proto:
	cd proto && make generate

.PHONY: build-image
build-image: build-gocron-server build-gocron-agent

.PHONY: build-gocron-server
build-gocron-server:
	docker build -f build/gocron-server/Dockerfile -t $(ORG)/gocron-server:latest .

.PHONY: build-gocron-agent
build-gocron-agent:
	docker build -f build/gocron-agent/Dockerfile -t $(ORG)/gocron-agent:latest .