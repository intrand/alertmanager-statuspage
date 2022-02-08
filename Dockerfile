# docker build --pull -t main:dev
FROM golang:alpine as builder
ARG Version
ARG Commit
ARG CommitDate
ARG Builder="buildx"
COPY . "${GOPATH}/src/package/app/"
WORKDIR "${GOPATH}/src/package/app"
RUN apk add --no-cache --upgrade \
		git \
		ca-certificates && \
	adduser -D -g '' app && \
	go get -d -v && \
	CGO_ENABLED=0 go build -a -o "/go/bin/main" -ldflags "-s -w -X 'main.version=${Version}' -X 'main.commit=${Commit}' -X 'main.date=${CommitDate}' -X 'main.builtBy=${Builder}'";

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/bin/main /go/bin/main
EXPOSE 8080
USER app
WORKDIR "/go/bin"
ENTRYPOINT ["/go/bin/main"]
