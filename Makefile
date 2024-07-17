include makefiles/const.mk
include makefiles/dependent.mk

APP_VERSION = v2.1.1
GOARCH=amd64
# Build pilotgo-front binary
pilotgo-front: ; $(info $(M)...Begin to build pilotgo-front binary.)  @ ## Build pilotgo-front.
	scripts/build.sh front;
# Build pilotgo-server binary
pilotgo-server: ; $(info $(M)...Begin to build pilotgo-server binary.)  @ ## Build pilotgo-server.
	scripts/build.sh backend ${GOARCH};
container: ;$(info $(M)...Begin to build the docker image.)  @ ## Build the docker image.
	DOCKER_BUILDKIT=0  docker build  --target pilotgo-server  -t pilotgo_server:latest . --no-cache
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
lint: golangci
	@{ \
	$(INFO) lint ;\
	GOFLAGS="-buildvcs=false" ;\
	go work edit -json | jq -r '.Use[].DiskPath'  | xargs -I{} $(GOLANGCILINT) run {}/... ;\
	}
# check-license-header: Check license header
check-license-header:
	./scripts/licence/header-check.sh
fix-license-header:
	./scripts/licence/header-check.sh fix
