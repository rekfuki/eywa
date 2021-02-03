FROM scratch

COPY root.pem /usr/local/ntt/ssl/
ADD https://raw.githubusercontent.com/ConnectCorp/docker-scratch-ssl/master/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ADD sawmill-api current_version /
ADD migrations /migrations

CMD ["/sawmill-api"]
