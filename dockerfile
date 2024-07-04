FROM node:20-alpine as ui

RUN npm config set registry https://registry.npmmirror.com/  && yarn config set registry https://registry.npmmirror.com/
COPY frontend/package.json frontend/yarn.lock frontend/

RUN yarn --cwd frontend install --network-timeout 1000000


COPY frontend frontend

RUN yarn --cwd frontend build

####################################################################################################

FROM golang:1.21-alpine3.18 as builder

ARG version_path="/out/backend/pilotgo/"
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk update && apk add --no-cache \
    git \
    make \
    ca-certificates \
    wget \
    curl \
    gcc \
    bash \
    build-base

WORKDIR /PilotGo
COPY . .
COPY --from=ui /frontend/dist/assets ./src/app/server/resource/assets
COPY --from=ui /frontend/dist/index.html  ./src/app/server/resource/index.html 
COPY --from=ui /frontend/dist/pilotgo.ico ./src/app/server/resource/pilotgo.ico

RUN cd /PilotGo/src/app/server && GOWORK=off GO111MODULE=on go build -mod=vendor -o ${version_path}/server/pilotgo-server -tags="production" main.go 
RUN  chmod a+x ${version_path}/server/pilotgo-server

####################################################################################################

# # FROM scratch as pilotgo-server 
FROM alpine:3.16.2 as pilotgo-server 


EXPOSE 8888 8889
WORKDIR /home/pilotgo

COPY --from=builder /out/backend/pilotgo/server/pilotgo-server  pilotgo-server 
COPY --from=builder /PilotGo/src/config_server.yaml.templete  config_server.yaml
COPY --from=builder /PilotGo/src/user.xlsx.templete   user.xlsx
CMD [ "./pilotgo-server" ]