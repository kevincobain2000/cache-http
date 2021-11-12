FROM golang:1.17.0-alpine

WORKDIR /app

COPY .env.example ./.env
COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN mkdir assets

COPY *.go ./

RUN go build -o /cache-http

# Because this is designed to run on a private Docker network
# this will not conflict with anything running on port 80 on the host
EXPOSE 80

# Listen to requests for all IP addresses so that the container
# can respond both the private name that Github Actions may assign
# as a service container, while also responding to "localhost" requests
# if you map the port to be accessible directly from the host.
#
# If you want to expose this service to the host and not the whole
# internet, make sure you Docker port numbers restrict to localhost:
#  -p 127.0.0.1:3000:80

CMD [ "/cache-http", "--host=", "--port=80" ]
