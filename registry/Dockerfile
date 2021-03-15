FROM alpine

RUN apk add --no-cache shadow-uidmap libseccomp-dev img

COPY templates /templates
COPY migrations /migrations
COPY registry /registry
COPY img /img

ADD https://raw.githubusercontent.com/ConnectCorp/docker-scratch-ssl/master/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt


CMD ["/registry"]