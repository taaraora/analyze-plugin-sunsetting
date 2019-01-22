FROM golang:1.11.4 as back_builder

ARG ARCH=amd64
ARG GO111MODULE=on

WORKDIR $GOPATH/src/github.com/supergiant/analyze-plugin-sunsetting/

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates

COPY go.mod go.sum $GOPATH/src/github.com/supergiant/analyze-plugin-sunsetting/
RUN go mod download

COPY . $GOPATH/src/github.com/supergiant/analyze-plugin-sunsetting/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} \
    go build -o $GOPATH/bin/analyze-sunsetting -a -installsuffix cgo -ldflags='-extldflags "-static" -w -s'  ./cmd/analyze-sunsetting

FROM scratch
COPY --from=back_builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=back_builder /go/bin/analyze-sunsetting /bin/analyze-sunsetting

ENTRYPOINT ["/bin/analyze-sunsetting"]
