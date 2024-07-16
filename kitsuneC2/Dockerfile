FROM golang:1.21.3

RUN mkdir -p /opt/kitsuneC2

WORKDIR /opt/kitsuneC2

ADD . /opt/kitsuneC2

RUN go mod download

WORKDIR /opt/kitsuneC2/server

RUN go build .

ENTRYPOINT [ "./server" ]