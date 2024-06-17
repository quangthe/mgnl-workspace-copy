ARG GO_VERSION=1.21

FROM postgres:15-alpine AS tools

FROM postgres:11-alpine AS postgres-11
FROM postgres:12-alpine AS postgres-12
FROM postgres:13-alpine AS postgres-13
FROM postgres:14-alpine AS postgres-13

FROM golang:${GO_VERSION}-alpine AS builder

ARG APP_VERSION="unversioned"

RUN apk add --no-cache ca-certificates git 
ENV CGO_ENABLED=0

WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./ ./

RUN --mount=type=cache,target=/root/.cache/go-build \
    go build \
    -installsuffix 'static' \
    -ldflags "-X main.version=$APP_VERSION" \
    -o /app

FROM tools AS final

# Copy postgres binaries.
COPY --from=postgres-11 /usr/local/bin /usr/lib/postgresql/11/bin
COPY --from=postgres-12 /usr/local/bin /usr/lib/postgresql/12/bin
COPY --from=postgres-13 /usr/local/bin /usr/lib/postgresql/13/bin
COPY --from=postgres-14 /usr/local/bin /usr/lib/postgresql/14/bin

COPY --from=postgres-11 /usr/local/lib /usr/lib/postgresql/11/lib
COPY --from=postgres-12 /usr/local/lib /usr/lib/postgresql/12/lib
COPY --from=postgres-13 /usr/local/lib /usr/lib/postgresql/13/lib
COPY --from=postgres-14 /usr/local/lib /usr/lib/postgresql/14/lib

COPY --from=postgres-11 /usr/local/share /usr/lib/postgresql/11/share
COPY --from=postgres-12 /usr/local/share /usr/lib/postgresql/12/share
COPY --from=postgres-13 /usr/local/share /usr/lib/postgresql/13/share
COPY --from=postgres-14 /usr/local/share /usr/lib/postgresql/14/share

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app /app

ENTRYPOINT [ "/app" ]
