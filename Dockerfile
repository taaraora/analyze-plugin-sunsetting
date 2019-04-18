FROM golang:1.11.8 as back_builder

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates

FROM scratch
COPY --from=back_builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY ./dist/analyze-sunsetting /bin/analyze-sunsetting

ENTRYPOINT ["/bin/analyze-sunsetting"]
