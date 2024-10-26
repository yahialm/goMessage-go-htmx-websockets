FROM golang:latest

WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . .

RUN go build -o gchat

EXPOSE 8080

CMD [ "./gchat" ]