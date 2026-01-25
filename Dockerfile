FROM golang:1.25.1-alpine

LABEL authors="nmori"

RUN apk --update add tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
#    apk del tzdata && \
    rm -rf /var/cache/apk/*

RUN apk upgrade --update && \
    apk --no-cache add \
    git \
    bash \
    curl \
    procps \
    ripgrep \
    build-base

# 開発用ツールをインストール（バージョン固定）
ARG GCI_VERSION=v0.13.7
ARG GOLANGCI_LINT_VERSION=v1.64.8
ENV PATH="$PATH:/go/bin"
ENV GOTOOLCHAIN=go1.25.1
ARG AIR_VER=v1.63.0
ARG DLV_VER=v1.25.2
ARG WIRE_VER=v0.7.0
ARG GOIMPORTS_VER=v0.24.0
ARG GOTESUM_VER=v1.13.0
ARG COBERTURA_VER=aaee18c8195c
ARG MOCKGEN_VER=v0.6.0

# go install はキャッシュしやすいようにまとめて実行
# ※ go get は使わない（モジュール外では禁止）
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go install github.com/air-verse/air@${AIR_VER} && \
    go install github.com/go-delve/delve/cmd/dlv@${DLV_VER} && \
    go install github.com/google/wire/cmd/wire@${WIRE_VER} && \
    go install golang.org/x/tools/cmd/goimports@${GOIMPORTS_VER} && \
    go install gotest.tools/gotestsum@${GOTESUM_VER} && \
    go install github.com/t-yuki/gocover-cobertura@${COBERTURA_VER} && \
    go install go.uber.org/mock/mockgen@${MOCKGEN_VER} && \
    go install github.com/daixiang0/gci@${GCI_VERSION} && \
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLANGCI_LINT_VERSION} && \
    go install mvdan.cc/gofumpt@latest

# note WSL環境だからか？ディレクトリを信頼できてない見たい
RUN git config --global --add safe.directory /go/src

WORKDIR /go/src

# 依存を先に解決（キャッシュ効かせる）
COPY go.mod go.sum ./
RUN go mod download

# ソースコピー
COPY . .

EXPOSE 8080 40000

CMD ["air", "-c", ".air.toml"]
#CMD ["dlv", "exec", "$(which air)", "--headless", "--listen=:40000", "--api-version=2", "--accept-multiclient", "--log"]
