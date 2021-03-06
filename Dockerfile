FROM golang:1.17.8-alpine3.15 AS builder
WORKDIR /go/src/app

RUN apk add --no-cache make

COPY . .
RUN go get -u github.com/go-swagger/go-swagger/cmd/swagger@latest
RUN make swagger
RUN CGO_ENABLED=0 go install -ldflags '-extldflags "-static"' -tags timetzdata

#==============
FROM scratch

COPY --from=builder /go/bin/users /users
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/users"]