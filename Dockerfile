FROM golang:1.10 AS builder
ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep
COPY . "/go/src/github.com/dannolan/webparrot"
WORKDIR "/go/src/github.com/dannolan/webparrot"
RUN dep ensure --vendor-only
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o /webparrot

FROM scratch
COPY --from=builder /webparrot .
EXPOSE 5000
EXPOSE 443
EXPOSE 80
CMD ["./webparrot"]
VOLUME ["/certs"]  