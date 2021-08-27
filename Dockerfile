FROM alpine:latest as certs
RUN apk --update add ca-certificates
# https://medium.com/on-docker/use-multi-stage-builds-to-inject-ca-certs-ad1e8f01de1b

FROM scratch
EXPOSE 80
ENV PATH=/bin
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY ./tweb1 /
CMD ["/tweb1"]
