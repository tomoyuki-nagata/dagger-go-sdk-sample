FROM golang:1.20.1

WORKDIR /src

COPY . /src

RUN pwd
RUN go mod download

ENTRYPOINT ["go", "run", "./main.go"]