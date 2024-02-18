FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -o apisvc

EXPOSE 8082

CMD ["./apisvc"]
