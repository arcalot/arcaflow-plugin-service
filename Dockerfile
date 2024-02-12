FROM golang:1.21-alpine@sha256:0ff68fa7b2177e8d68b4555621c2321c804bcff839fd512c2681de49026573b7 AS build

ENV CGO_ENABLED=0
COPY . /src
WORKDIR /src
RUN go build -o service-plugin cmd/service-plugin/main.go

FROM scratch
COPY --from=build /src/service-plugin /service-plugin
WORKDIR /
ENTRYPOINT ["/service-plugin"]
CMD []