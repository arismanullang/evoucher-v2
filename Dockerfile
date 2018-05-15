FROM golang:1-stretch

# Copy files
COPY . "$GOPATH/src/github.com/gilkor/evoucher"
COPY ./files/etc/* "/etc/"
RUN go install github.com/gilkor/evoucher

# Start
EXPOSE 8080
