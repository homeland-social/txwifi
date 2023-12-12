FROM golang:1.21-alpine3.18 AS builder

ENV GOPATH /go

RUN mkdir -p /go/src/github.com/homeland-social/txwifi
COPY . /go/src/github.com/homeland-social/txwifi

WORKDIR /go/src/github.com/homeland-social/txwifi
RUN CGO_ENABLED=0 GOOS=linux go build -o /go/bin/wifi-server main.go

FROM alpine:3.18

RUN apk update
RUN apk add bridge hostapd wireless-tools wpa_supplicant dnsmasq iw ethtool

RUN mkdir -p /etc/wpa_supplicant/
COPY ./dev/configs/wpa_supplicant.conf /etc/wpa_supplicant/wpa_supplicant.conf

WORKDIR /

COPY --from=builder /go/bin/wifi-server /wifi-server
ENTRYPOINT ["/wifi-server"]


