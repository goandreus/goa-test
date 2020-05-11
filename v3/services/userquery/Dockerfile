# Dockerfile for userquery micro-service
FROM alpine:3.9

RUN addgroup -g 10001 -S wsgroup && adduser -u 10001 -S wsuser -G wsgroup

FROM scratch

ARG BUILD_DATE

ENV BUILD_DATE=$BUILD_DATE

LABEL org.label-schema.name="userquery"
LABEL org.label-schema.description="API service to manage authentication and security. The service integrates its own cache and is optimized for query. Data synchronization with other services is done via a shared NATS Streaming event bus."
LABEL org.label-schema.build-date=$BUILD_DATEma
LABEL org.label-schema.vendor="WiserSKILLS"
LABEL org.label-schema.version="1.0.0"

LABEL maintainer="christophe.peillet@wiserskills.com"


WORKDIR /root/

COPY --from=0 /etc/passwd /etc/passwd
COPY build/userquery .

USER wsuser

EXPOSE 8080 8090
ENTRYPOINT ["./userquery"]
