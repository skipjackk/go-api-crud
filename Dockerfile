FROM golang:latest

LABEL maintainer="Skipjackk crew"

WORKDIR /app

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

ENV PORT 3000

RUN go build

CMD ["go-api-crud.exe"]