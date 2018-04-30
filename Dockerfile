FROM golang:1-stretch

# Copy files
COPY . "$GOPATH/src/github.com/gilkor/evoucher"
COPY ./files/etc/* "/etc/"
RUN go install github.com/gilkor/evoucher

# Set environment variables
ENV EVOUCHER_CONFIG "/etc/evoucher/config.yml"

# Start
EXPOSE 8080
