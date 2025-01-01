FROM golang:1.23.4

WORKDIR /app

COPY . /app

RUN go mod download

EXPOSE 8080

CMD ["go", "run", "."]
