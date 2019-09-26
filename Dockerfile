# Build sources
FROM golang:1-alpine
ADD . "$GOPATH/src/github.com/gilkor/evoucher-v2"
RUN go install github.com/gilkor/evoucher-v2

# Copy built executables
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
RUN [ ! -e /etc/nsswitch.conf ] && echo 'hosts: files dns' > /etc/nsswitch.conf
COPY --from=0 /go/bin/evoucher-v2 /usr/local/bin/evoucher-v2
COPY ./files/etc/evoucher-v2 "/etc/evoucher-v2"
CMD echo "Use the app commands."; exit 1