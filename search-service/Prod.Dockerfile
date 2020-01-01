# Compile stage
FROM golang:1.20-alpine3.17 AS build-env
ENV CGO_ENABLED 0
ADD . /go/src/search-service
WORKDIR /go/src/search-service
RUN go build -o search_services main.go

# Final stage
FROM alpine:3.17
EXPOSE 8080

WORKDIR /

COPY .env /
COPY --from=build-env /go/src/search-service/search_services /
CMD ["/search_services"]

# ENTRYPOINT ["tail", "-f", "/dev/null"]