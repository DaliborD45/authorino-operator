# Build the authorino binary
ARG DATA_RACE=false

# https://catalog.redhat.com/software/containers/ubi10/go-toolset
FROM --platform=$TARGETPLATFORM registry.access.redhat.com/ubi10/go-toolset:1.26 AS builder
WORKDIR /workspace

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download
# Copy the go source
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/
COPY pkg/ pkg/

ARG OPERATOR_VERSION=latest
ARG DEFAULT_AUTHORINO_IMAGE=quay.io/kuadrant/authorino:latest
ARG GIT_SHA=unknown
ARG DIRTY=unknown
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT
ARG DATA_RACE=false

RUN if [ "${DATA_RACE}" = "true" ]; then  \
    CGO_ENABLED=1 GOOS=${TARGETOS} GOARCH=${TARGETARCH} GOARM=${TARGETVARIANT} \
        go build -race -a -ldflags "-X main.version=${OPERATOR_VERSION} -X main.gitSHA=${GIT_SHA} -X main.dirty=${DIRTY} -X github.com/kuadrant/authorino-operator/pkg/reconcilers.DefaultAuthorinoImage=${DEFAULT_AUTHORINO_IMAGE}" \
        -o /tmp/manager main.go; \
    else \
    CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} GOARM=${TARGETVARIANT} \
        go build -a -ldflags "-X main.version=${OPERATOR_VERSION} -X main.gitSHA=${GIT_SHA} -X main.dirty=${DIRTY} -X github.com/kuadrant/authorino-operator/pkg/reconcilers.DefaultAuthorinoImage=${DEFAULT_AUTHORINO_IMAGE}" \
        -o /tmp/manager main.go; \
    fi


# CGO_ENABLED=1 (race) produces a dynamically linked binary requiring glibc (ubi9 full)
# https://catalog.redhat.com/software/containers/ubi9-minimal
FROM registry.access.redhat.com/ubi9:latest AS runtime-true
FROM registry.access.redhat.com/ubi9-minimal:latest AS runtime-false

FROM runtime-${DATA_RACE}
WORKDIR /
COPY --from=builder /tmp/manager .
USER 1001

ARG QUAY_IMAGE_EXPIRY=never
LABEL quay.expires-after=$QUAY_IMAGE_EXPIRY

ENTRYPOINT ["/manager"]
