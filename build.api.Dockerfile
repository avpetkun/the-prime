FROM golang:alpine AS app

RUN apk add ca-certificates git gcc musl-dev make

WORKDIR /srv

ARG GIT_CREDENTIALS
RUN git config --global --add url."https://gitlab-ci-token:${GIT_CREDENTIALS}@gitlab.com".insteadOf "https://gitlab.com"

COPY go.mod go.sum ./
RUN  go mod download

COPY cmd      cmd
COPY pkg      pkg
COPY internal internal

ARG GIT_HASH

RUN go build -o app ./cmd/api



FROM alpine:latest

WORKDIR /srv

COPY --from=app /srv/app .

CMD /srv/app
