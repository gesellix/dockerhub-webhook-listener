FROM golang:1.6
ADD . /go/src/github.com/gesellix/dockerhub-webhook-listener
WORKDIR /go/src/github.com/gesellix/dockerhub-webhook-listener/hub-listener
RUN go get && go build
ENTRYPOINT ["/go/src/github.com/gesellix/dockerhub-webhook-listener/hub-listener/hub-listener"]
CMD ["-listen", "0.0.0.0:80"]
