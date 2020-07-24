FROM golang:1.14.1-buster AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download
COPY ./batch .
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-w -s' -o syncer main.go

FROM debian:buster-slim

RUN apt-get update \
  && apt-get install -y curl ca-certificates --no-install-recommends \
  && curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl \
  && chmod +x ./kubectl \
  && mv ./kubectl /usr/local/bin/kubectl \
  && rm -rf /var/lib/apt/lists/*

COPY --from=build /app/syncer /

CMD ["/syncer"]