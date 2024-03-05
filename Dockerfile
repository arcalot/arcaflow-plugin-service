FROM golang:1.22-alpine@sha256:fc5e5848529786cf1136563452b33d713d5c60b2c787f6b2a077fa6eeefd9114 AS build

ENV CGO_ENABLED=0
COPY . /src
WORKDIR /src
RUN go build -o service-plugin cmd/service-plugin/main.go

FROM scratch
COPY --from=build /src/service-plugin /service-plugin
WORKDIR /
ENTRYPOINT ["/service-plugin"]
CMD []
