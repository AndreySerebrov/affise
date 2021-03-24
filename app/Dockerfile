FROM golang:1.15 as builder

WORKDIR /app
COPY . .

RUN go build -o checkout ./cmd/main.go 


FROM ubuntu 

WORKDIR /app
COPY --from=builder /app/checkout .
COPY --from=builder /app/db ./db/

CMD /app/checkout
