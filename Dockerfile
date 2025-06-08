FROM golang:1.24.4-alpine

LABEL authors="nmori"

WORKDIR /go/src
COPY . .

RUN apk upgrade --update && \
    apk --no-cache add \
    git \
    bash \
    curl \
    procps

RUN go get -u github.com/air-verse/air && \
    go build -o /go/bin/air github.com/air-verse/air

# note WSL環境だからか？ディレクトリを信頼できてない見たい
RUN git config --global --add safe.directory /go/src

CMD ["air", "-c", ".air.toml"]