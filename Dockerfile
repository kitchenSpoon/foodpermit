FROM golang:1.14

COPY . /project
WORKDIR /project

RUN go get -d -v ./...
RUN go build -v ./cmd/project/main.go

CMD ["/project/main"]
