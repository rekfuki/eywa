FROM scratch

COPY go-watchdog /fwatchdog
COPY handler /handler

ADD https://raw.githubusercontent.com/ConnectCorp/docker-scratch-ssl/master/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENV fprocess="/handler"
ENV mode="http"
ENV upstream_url="http://127.0.0.1:8082"

CMD ["/fwatchdog"]