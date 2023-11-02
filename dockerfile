FROM golang:1.21-alpine

COPY . /t4m

WORKDIR /t4m

RUN go mod tidy

RUN go build -o t4m .

CMD ["/t4m/t4m"]