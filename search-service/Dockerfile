# Compile stage
FROM golang:1.20-alpine3.17 AS build-env
ENV CGO_ENABLED 0

ADD . /go/src/search-service
WORKDIR /go/src/search-service
RUN go build -gcflags "all=-N -l" -o search_services main.go

# download Delve
RUN go install github.com/go-delve/delve/cmd/dlv@latest

## Final stage
FROM alpine:3.17
EXPOSE 8080 40000

# Allow delve to run on Alpine based containers.
RUN apk add --no-cache libc6-compat

WORKDIR /
COPY --from=build-env /go/src/search-service/search_services /
COPY --from=build-env /go/bin/dlv /
COPY .env /
CMD ["/dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/search_services"]

#ENTRYPOINT ["tail", "-f", "/dev/null"]