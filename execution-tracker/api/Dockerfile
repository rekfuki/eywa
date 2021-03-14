FROM scratch

COPY migrations /migrations

COPY execution-tracker-api /execution-tracker-api

ADD https://raw.githubusercontent.com/ConnectCorp/docker-scratch-ssl/master/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

CMD ["/execution-tracker-api"]