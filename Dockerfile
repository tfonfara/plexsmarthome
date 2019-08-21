#https://docs.docker.com/develop/develop-images/multistage-build
#https://medium.com/@chemidy/create-the-smallest-and-secured-golang-docker-image-based-on-scratch-4752223b7324

FROM golang:alpine as builder

LABEL stage=intermediate

RUN apk --no-cache add \
    git \
    curl

RUN curl -s https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/github.com/tfonfara/plexsmarthome

COPY Gopkg.toml .
COPY Gopkg.lock .

RUN touch vendor.go
RUN dep ensure -vendor-only
RUN rm vendor.go

COPY . .

ARG VERSION

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -a \
    -installsuffix cgo \
    -o /go/bin/plexsmarthome .



FROM alpine:latest

RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    tini \
    su-exec

RUN update-ca-certificates 2>/dev/null || true

COPY --from=builder /go/bin/plexsmarthome /bin/plexsmarthome

COPY entrypoint.sh /
RUN chmod +x /entrypoint.sh

ENV PLEX_PORT 3000
ENV PLEX_CONFIG /config

VOLUME $PLEX_CONFIG
EXPOSE $PLEX_PORT/tcp

HEALTHCHECK --interval=30s --timeout=3s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:$PLEX_PORT/health || exit 1

ENTRYPOINT ["/sbin/tini", "--", "/entrypoint.sh"]
CMD ["/bin/plexsmarthome"]
