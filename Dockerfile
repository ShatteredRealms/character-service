# Build application
FROM golang:1.24 AS build
WORKDIR /src
ENV CGO_ENABLED=0
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
ARG APP_VERSION=v0.0.1
RUN go build \
	-ldflags="-X 'github.com/ShatteredRealms/character-service/pkg/config/default.Version=${APP_VERSION}'" \
	-o /out/character ./cmd/character-service

# Run server
FROM alpine:3.21.3
WORKDIR /app
COPY --from=build /out/character ./
EXPOSE 8081
ENTRYPOINT [ "./character" ]
