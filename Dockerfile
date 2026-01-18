FROM golang:1.25-alpine AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./bin/app cmd/main.go

FROM alpine

COPY --from=builder /usr/src/app/bin/app /

CMD ["/app"]