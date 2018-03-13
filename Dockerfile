FROM golang:latest
WORKDIR /go/src/github.com/cisco-sso/mh
RUN go get -u github.com/golang/dep/cmd/dep
COPY . .
RUN dep ensure \
&& CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mh .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /go/src/github.com/cisco-sso/mh/mh .
ENTRYPOINT ["./mh"]
