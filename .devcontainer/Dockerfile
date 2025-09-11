FROM migrate/migrate AS mig

FROM golang:1.25

COPY --from=mig /usr/local/bin/migrate /usr/local/bin/migrate

RUN apt update && apt install -y libasound2-dev

RUN go install golang.org/x/tools/gopls@latest && \
    go install honnef.co/go/tools/cmd/staticcheck@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install go.uber.org/mock/mockgen@v0.6.0 && \
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(go env GOPATH)/bin v2.4.0
