FROM golang:latest

LABEL maintainer="Peter Carr <carrpet@gmail.com>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

EXPOSE 8080

CMD ["./main"]
