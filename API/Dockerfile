FROM golang:1.20.13-alpine3.19

WORKDIR /build

COPY . . 

RUN go build main.go

EXPOSE 3000

CMD ["./main"]
