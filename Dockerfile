FROM golang:1.21-alpine@sha256:4db4aac30880b978cae5445dd4a706215249ad4f43d28bd7cdf7906e9be8dd6b AS build

ENV CGO_ENABLED=0
COPY . /src
WORKDIR /src
RUN go build -o service-plugin cmd/service-plugin/main.go

FROM scratch
COPY --from=build /src/service-plugin /service-plugin
WORKDIR /
ENTRYPOINT ["/service-plugin"]
CMD []