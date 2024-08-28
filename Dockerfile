FROM golang:1.22-alpine as buildbase

RUN apk add git build-base ca-certificates

WORKDIR /go/src/github.com/rarimo/verificator-svc
COPY . .

RUN go mod tidy && go mod vendor
RUN CGO_ENABLED=1 GO111MODULE=on GOOS=linux GOOS=linux go build -o /usr/local/bin/verificator-svc /go/src/github.com/rarimo/verificator-svc

FROM scratch
COPY --from=alpine:3.9 /bin/sh /bin/sh
COPY --from=alpine:3.9 /usr /usr
COPY --from=alpine:3.9 /lib /lib

COPY --from=buildbase /usr/local/bin/geo-forms-svc /usr/local/bin/verificator-svc
COPY --from=buildbase /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["verificator-svc"]