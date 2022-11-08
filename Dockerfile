FROM golang:alpine as app-builder
WORKDIR /go/src/app
COPY . .
RUN echo "Cache break counter: 7"
# Static build required so that we can safely copy the binary over.
# `-tags timetzdata` embeds zone info from the "time/tzdata" package.
RUN CGO_ENABLED=0 go install -ldflags '-extldflags "-static"' -tags timetzdata ./...

FROM scratch

COPY --from=app-builder /go/bin/splatbackend /go-server-splatbackend

CMD ["/go-server-splatbackend"]