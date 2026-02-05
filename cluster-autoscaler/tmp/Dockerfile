FROM docker.io/alpine:3.22.3 as certificates
RUN apk add --no-cache ca-certificates bash

FROM scratch
COPY --from=certificates /etc/ssl /etc/ssl
COPY cluster-autoscaler-amd64 /cluster-autoscaler
CMD ["/cluster-autoscaler"]