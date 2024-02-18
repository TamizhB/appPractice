FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -o kafkasvc

EXPOSE 9099

CMD ["./kafkasvc"]
