FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -o udpSvc

EXPOSE 3000

CMD ["./udpSvc"]
