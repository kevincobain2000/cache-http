FROM golang:1.17.0-alpine

WORKDIR /app

COPY .env.example ./.env
COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN mkdir assets

COPY *.go ./

RUN go build -o /cache-http

EXPOSE 3000

CMD [ "/cache-http", "3000" ]
