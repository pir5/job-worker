FROM golang:latest AS build-env

ENV CGO_ENABLED=0
WORKDIR /go/src/health-worker
COPY . .

RUN mkdir -p /build
RUN go build -a  -ldflags="-s -w -extldflags \"-static\"" -o=/build/health-worker main.go

FROM alpine:3
# Timezone = Tokyo
RUN apk --no-cache add tzdata zlib && \
    apk add --upgrade --no-cache && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime

COPY --from=build-env /build/health-worker /build/health-worker
RUN chmod u+x /build/health-worker

ENTRYPOINT ["/build/health-worker", "server"]
