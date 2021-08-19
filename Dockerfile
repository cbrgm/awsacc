FROM alpine:latest

ENV AWSACC_CONFIG=/data/accounts.json

RUN mkdir /data
COPY bin/awsacc_linux_docker_amd64 /app/awsacc

ENTRYPOINT ["/app/awsacc"]
