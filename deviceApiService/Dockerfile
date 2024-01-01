FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -o apisvc

EXPOSE 8080

CMD ["./apisvc"]
