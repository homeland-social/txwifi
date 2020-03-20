FROM arm32v7/golang:1.13-alpine3.11 AS builder

ENV GOPATH /go
WORKDIR /go/src

RUN mkdir -p /go/src/github.com/kinokochat/txwifi
COPY . /go/src/github.com/kinokochat/txwifi

RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/wifi-server /go/src/github.com/kinokochat/txwifi/main.go

FROM arm32v7/alpine:3.11

RUN apk update
RUN apk add bridge hostapd wireless-tools wpa_supplicant dnsmasq iw

RUN mkdir -p /etc/wpa_supplicant/
COPY ./dev/configs/wpa_supplicant.conf /etc/wpa_supplicant/wpa_supplicant.conf

WORKDIR /

COPY --from=builder /go/bin/wifi-server /wifi-server
ENTRYPOINT ["/wifi-server"]


