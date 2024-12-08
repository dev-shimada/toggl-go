FROM --platform=$BUILDPLATFORM golang:1.23.1-bookworm as vscode

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH
ENV LESSCHARSET utf-8

WORKDIR /service

RUN  <<EOF
apt-get update && apt-get install -y vim git unzip less bash-completion
go install -v golang.org/x/tools/gopls@latest
go install -v github.com/go-delve/delve/cmd/dlv@latest
EOF

RUN  <<EOF
go install -v github.com/goreleaser/goreleaser/v2@latest
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin latest
EOF

COPY . /service
