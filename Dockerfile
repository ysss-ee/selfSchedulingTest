# Scheduler Framework plugins are compiled into kube-scheduler; they cannot be
# installed into a stock scheduler image using only a YAML configuration file.
FROM golang:1.26-bookworm AS build

ARG KUBERNETES_VERSION=v1.36.1
RUN apt-get update \
    && apt-get install -y --no-install-recommends git ca-certificates \
    && rm -rf /var/lib/apt/lists/*
RUN git clone --depth 1 --branch ${KUBERNETES_VERSION} https://github.com/kubernetes/kubernetes.git /kubernetes

WORKDIR /kubernetes
COPY cmd/self-scheduler /kubernetes/cmd/self-scheduler
RUN go test -mod=vendor ./cmd/self-scheduler/... \
    && go build -mod=vendor -trimpath -ldflags="-s -w" -o /out/self-scheduler ./cmd/self-scheduler
FROM registry.k8s.io/kube-scheduler:v1.36.1
COPY --from=build /out/self-scheduler /usr/local/bin/self-scheduler
USER 65532:65532
ENTRYPOINT ["/usr/local/bin/self-scheduler"]
