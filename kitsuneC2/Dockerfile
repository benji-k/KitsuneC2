FROM golang:1.21.3

WORKDIR /opt/kitsuneC2

COPY . /opt/kitsuneC2

RUN go mod download

WORKDIR /opt/kitsuneC2/server

RUN go build .

RUN apt update && apt install sqlite3

RUN addgroup --system appgroup && adduser --system --ingroup appgroup appuser

USER appuser

ENTRYPOINT [ "./server" ]