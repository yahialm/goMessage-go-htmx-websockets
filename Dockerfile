FROM golang:latest

WORKDIR /app

COPY go.mod .
RUN go mod download

COPY . .

RUN go build -o gchat

EXPOSE 3000

CMD [ "./gchat" ]