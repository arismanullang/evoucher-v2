# Build sources
FROM golang:1-alpine
ADD . "$GOPATH/src/github.com/gilkor/evoucher-v2"
RUN go install github.com/gilkor/evoucher-v2

# Copy built executables
FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
RUN [ ! -e /etc/nsswitch.conf ] && echo 'hosts: files dns' > /etc/nsswitch.conf
COPY --from=0 /go/bin/* /usr/local/bin/
CMD echo "Use the app commands."; exit 1
