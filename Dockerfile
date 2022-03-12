FROM golang:1.17.8-alpine3.15 AS builder
WORKDIR /go/src/app

COPY . .
RUN go get -d -v ./...
RUN go install -v ./...

#==============
FROM scratch

COPY --from=builder /go/bin/users /users

CMD ["/users"]