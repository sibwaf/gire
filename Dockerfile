FROM golang:1.19-alpine AS builder

WORKDIR /build
COPY src src
COPY gire.go .
COPY go.mod .
COPY go.sum .

RUN go build

FROM alpine:3.17

RUN apk add openssh && \
    apk add git && \
    add-shell /usr/bin/git-shell && \
    echo "StrictHostKeyChecking yes" >> /etc/ssh/ssh_config && \
    echo "PasswordAuthentication no" >> /etc/ssh/sshd_config && \
    echo "ForceCommand /bin/git-readonly" >> /etc/ssh/sshd_config

COPY git-readonly.sh /bin/git-readonly

WORKDIR /app
COPY --from=builder /build/gire /app/gire
COPY entrypoint.sh .

RUN chmod +x /bin/git-readonly gire entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]
