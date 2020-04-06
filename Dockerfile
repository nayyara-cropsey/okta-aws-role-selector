# go build
FROM golang:latest
LABEL maintainer="Nayyara Samuel <nayyara.samuel@gmail.com>"
WORKDIR /app

COPY go.mod go.sum ./
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# run time build
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# to mount additional user-defined configs
RUN mkdir -p /root/config/
COPY --from=0 /app/main .
COPY --from=0 /app/views ./views
COPY --from=0 /app/assets ./assets
COPY --from=0 /app/config.yaml ./config/

EXPOSE 8080
ENV GIN_MODE=release
ENTRYPOINT ["./main", "-c", "config/config.yaml"]