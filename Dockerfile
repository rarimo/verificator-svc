FROM golang:1.19.7-alpine as buildbase

RUN apk add git build-base ca-certificates

WORKDIR /go/src/verificator-svc
COPY . .
RUN GOOS=linux go build -o /usr/local/bin/verificator-svc /go/src/verificator-svc

FROM scratch
COPY --from=alpine:3.9 /bin/sh /bin/sh
COPY --from=alpine:3.9 /usr /usr
COPY --from=alpine:3.9 /lib /lib

COPY --from=buildbase /usr/local/bin/verificator-svc /usr/local/bin/verificator-svc
COPY --from=buildbase /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["verificator-svc"]
