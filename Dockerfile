
FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git make bash
ARG VERSION=dev
ENV IMAGE_TAG=${VERSION}

WORKDIR $GOPATH/src/scanner
COPY . .
RUN go mod download && go mod verify 
RUN make o=/go/bin/scanner

FROM alpine:3.11
ARG PUID=2000
ARG PGID=2000

RUN apk add --no-cache curl

COPY --from=builder /go/bin/scanner /go/bin/scanner

#Healthcheck to make sure container is ready
HEALTHCHECK --interval=5m --timeout=3s \
    CMD curl -f http://localhost:6080/healthCheck || exit 1

ENTRYPOINT ["/go/bin/scanner"]
EXPOSE 6080