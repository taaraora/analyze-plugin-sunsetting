FROM golang:1.11.4 as job_builder

ARG ARCH=amd64
ARG GO111MODULE=on

WORKDIR $GOPATH/src/github.com/supergiant/analyze-plugin-sunsetting/

COPY go.mod go.sum $GOPATH/src/github.com/supergiant/analyze-plugin-sunsetting/
RUN go mod download

COPY . $GOPATH/src/github.com/supergiant/analyze-plugin-sunsetting/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} \
    go build -o $GOPATH/bin/analyze-sunsetting-job -a -installsuffix cgo -ldflags='-extldflags "-static" -w -s'  ./cmd/analyze-sunsetting-job

FROM busybox:latest
COPY --from=job_builder /go/bin/analyze-sunsetting-job /bin/analyze-sunsetting-job

ENTRYPOINT ["/bin/analyze-sunsetting-job"]
