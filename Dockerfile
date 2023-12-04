FROM golang:1.21-alpine@sha256:70afe55365a265f0762257550bc38440e0d6d6b97020d3f8c85328f00200dd8e AS build

ENV CGO_ENABLED=0
COPY . /src
WORKDIR /src
RUN go build -o service-plugin cmd/service-plugin/main.go

FROM scratch
COPY --from=build /src/service-plugin /service-plugin
WORKDIR /
ENTRYPOINT ["/service-plugin"]
CMD []