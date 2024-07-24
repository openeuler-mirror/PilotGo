include makefiles/const.mk
include makefiles/dependent.mk

# Build pilotgo-front binary
pilotgo-front: ; $(info $(M)...Begin to build pilotgo-front binary.)  @ ## Build pilotgo-front.
	scripts/build.sh front;

# Build pilotgo-server binary
pilotgo-server: ; $(info $(M)...Begin to build pilotgo-server binary.)  @ ## Build pilotgo-server.
	scripts/build.sh backend ${GOARCH};

pilotgo-server-debug: ; $(info $(M)...Begin to build pilotgo-server binary.)  @ ## Build pilotgo-server.
	CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build  -tags=production -ldflags '${LDFLAGS}' -o ./pilotgo-server ./src/app/server/main.go

container: ;$(info $(M)...Begin to build the docker image.)  @ ## Build the docker image.
	DOCKER_BUILDKIT=0  docker build  --target pilotgo-server  -t pilotgo_server:${IMAGE_TAG} . --no-cache

clean: ; $(info $(M)...Begin to clean out dir.)  @ ## clean out dir.
	scripts/build.sh clean;

all: clean pilotgo-front pilotgo-server

pack: ; $(info $(M)...Begin to pack tar package  dir.)  @ ##  pack tar package .
	scripts/build.sh pack ${GOARCH}

docker-compose-up: ; $(info $(M)...Begin to deploy by docker-compose.)  @ ## deploy by docker-compose.
	DOCKER_BUILDKIT=0 docker-compose -f scripts/dockercompose/docker-compose.yml up

docker-compose-down: ; $(info $(M)...Begin to stop by docker-compose.)  @ ## stop by docker-compose.
	DOCKER_BUILDKIT=0 docker-compose -f scripts/dockercompose/docker-compose.yml down -v

docker-compose-build: ; $(info $(M)...Begin to build image by docker-compose.)  @ ## build image by docker-compose.
	DOCKER_BUILDKIT=0 docker-compose -f scripts/dockercompose/docker-compose.yml build --no-cache

## lint: Run the golangci-lint
lint: golangci ; $(info $(M)...Begin to check  code.)  @ ## check  code.
	@{ \
	$(INFO) lint ;\
	$(INFO) $(ROOT_DIR) ;\
	GOFLAGS="-buildvcs=false" ;\
	# go list -f '{{.Dir}}' -m | xargs -I {} bash -c ' cd "{}" && $(GOLANGCILINT) run ./... --config $(ROOT_DIR)/.golangci.yml'  ;\
	go work edit -json | jq -r '.Use[].DiskPath' | xargs -I {} bash -c ' cd "{}" && $(GOLANGCILINT) run ./... --config $(ROOT_DIR)/.golangci.yml ' ;\
	}
# check-license-header: Check license header
check-license-header:
	./scripts/licence/header-check.sh

fix-license-header:
	./scripts/licence/header-check.sh fix

build-server-templete:  ; $(info $(M)...Begin to build config_server.yaml.templete.)  @ ## build config_server.yaml.templete
	go run src/app/server/main.go templete
	cp -R config_server.yaml.templete ./src/config_server.yaml.templete
	rm -rf ./config_server.yaml.templete
vendor: ; $(info $(M)...Begin to update vendor.)  @ ## update vendor
	cd src && GOWORK=off && go mod vendor