FROM golang:1.21.3

WORKDIR /opt/kitsuneC2

COPY . /opt/kitsuneC2

RUN go mod download

WORKDIR /opt/kitsuneC2/server

RUN go build .

ENTRYPOINT [ "./server" ]