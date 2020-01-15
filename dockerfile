FROM golang:latest

RUN mkdir -p /src
WORKDIR /src

CMD ["go", "run", "coffee.go"]
