FROM golang:1.24.4-alpine

LABEL authors="nmori"

RUN apk --update add tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    apk del tzdata && \
    rm -rf /var/cache/apk/*

WORKDIR /go/src
RUN mkdir logs

COPY . .

RUN apk upgrade --update && \
    apk --no-cache add \
    git \
    bash \
    curl \
    procps

RUN go get -u github.com/air-verse/air && \
    go build -o /go/bin/air github.com/air-verse/air && \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install github.com/google/wire/cmd/wire@latest && \
    go install go.uber.org/mock/mockgen@latest && \
    go get go.uber.org/mock

ENV PATH="$PATH:/go/bin"

# note WSL環境だからか？ディレクトリを信頼できてない見たい
RUN git config --global --add safe.directory /go/src

# expose both Gin port and Delve
EXPOSE 8080 40000

CMD ["dlv", "exec", "$(which air)", "--headless", "--listen=:40000", "--api-version=2", "--accept-multiclient", "--log"]
#CMD ["air", "-c", ".air.toml"]