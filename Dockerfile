FROM golang:1.18-alpine AS build

ENV CGO_ENABLED=0
COPY . /src
WORKDIR /src
RUN go build -o service-plugin cmd/service-plugin/main.go

FROM scratch
COPY --from=build /src/service-plugin /service-plugin
WORKDIR /
ENTRYPOINT ["/service-plugin"]
CMD []