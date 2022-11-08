ARG MONGODB_URI
ARG MONGO_MAX_IDLE_TIME_MS
ARG MONGO_MAX_POOL_SIZE
ARG MONGO_MIN_POOL_SIZE
ARG PORT
ARG REDIS_URI
ARG SECRET_KEY
FROM golang:alpine as app-builder
WORKDIR /go/src/app

COPY . .

FROM scratch

COPY --from=app-builder /go/bin/splatbackend /go-server-splatbackend

CMD ["/go-server-splatbackend"]