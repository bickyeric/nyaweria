FROM golang:1.23

RUN apt update && apt install -y libasound2-dev

RUN go install golang.org/x/tools/gopls@latest && \
    go install honnef.co/go/tools/cmd/staticcheck@latest
