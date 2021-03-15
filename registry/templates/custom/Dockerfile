FROM registry.eywa.rekfuki.dev/of-watchdog:latest as watchdog
FROM alpine:3.12
# Add non root user and certs
RUN apk --no-cache add ca-certificates tree \
    && addgroup -S app && adduser -S -g app app \
    && mkdir -p /home/app \
    && chown app /home/app

WORKDIR /home/app

ARG FPROCESS
COPY --chown=app ./function/  ./
COPY --from=watchdog --chown=app /of-watchdog ./fwatchdog
RUN chmod +x ./fwatchdog

USER app

RUN if [[ ! -e $FPROCESS ]]; then >&2 echo "Provided executable ${FPROCESS} does not exist..."; exit 1; fi
RUN RESULT=$(ldd $FPROCESS 2>/dev/null) \
    && {>&2 echo "Found DLLs:"; >&2 echo "${RESULT}"; >&2 echo "Custom runners must be compiled statically"; exit 1;} \
    || echo "Static Binary"

RUN chmod +x $FPROCESS

RUN ls -la
ENV fprocess="./${FPROCESS}"
ENV mode="http"
ENV upstream_url="http://127.0.0.1:8082"
RUN env

CMD ["./fwatchdog"]