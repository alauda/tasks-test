FROM golang:1.7

COPY . /go/src/app

WORKDIR /go/src/app

RUN go install

CMD ["/go/bin/app"]