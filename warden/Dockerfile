FROM scratch

COPY templates /templates
COPY migrations /migrations
COPY warden /warden

ADD https://raw.githubusercontent.com/ConnectCorp/docker-scratch-ssl/master/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt


CMD ["/warden"]