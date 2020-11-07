FROM --platform=${TARGETPLATFORM:-linux/amd64} openfaas/of-watchdog:0.8.1 as watchdog
FROM scratch

COPY --from=watchdog /fwatchdog /
COPY handler /

ENV fprocess="./handler"
ENV mode="http"
ENV upstream_url="http://127.0.0.1:8082"

CMD ["./fwatchdog"]