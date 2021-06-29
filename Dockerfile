FROM scratch
ENTRYPOINT ["/highlander"]
COPY ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY highlander /
