FROM golang:1.17.8-alpine3.15 AS builder
WORKDIR /go/src/app

RUN apk add --no-cache make

COPY . .
RUN go get -u github.com/go-swagger/go-swagger/cmd/swagger@latest
RUN make swagger
RUN go get -d -v ./...
RUN go install -v ./...

#==============
FROM scratch

COPY --from=builder /go/bin/users /users

CMD ["/users"]