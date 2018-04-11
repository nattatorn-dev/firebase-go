FROM golang:1.9 as builder
ADD ./ /go/src/github.com/nattatorn-dev/log-manager
WORKDIR /go/src/github.com/nattatorn-dev/log-manager
RUN go get ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .


#Container to run
FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/nattatorn-dev/log-manager/app .
CMD ["./app"]  


EXPOSE 8080
