FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder

ARG TARGETOS
ARG TARGETARCH
ARG VERSION

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags "-s -w -X main.version=${VERSION}" -o /checks-qe ./cmd/checks-qe/

FROM registry.access.redhat.com/ubi9/ubi-minimal:9.6

COPY --from=builder /checks-qe /usr/local/bin/checks-qe

USER 1001

ENTRYPOINT ["checks-qe"]
CMD ["--help"]
