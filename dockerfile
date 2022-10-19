FROM alpine:latest

ARG ARCH

WORKDIR /opt/PilotGo/server

COPY ./out/${ARCH}/pilotgo-v0.0.1/server/ /opt/PilotGo/server

EXPOSE 8888 8889


CMD [ "/opt/PilotGo/server/pilotgo-server" ]

