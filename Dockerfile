ARG GO_VERSION=latest
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} AS build
WORKDIR /app

ARG VERSION="undefined"
ENV VERSION=${VERSION}

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT
ENV GOOS=${TARGETOS}
ENV GOARCH=${TARGETARCH}
ENV GOARM=${TARGETVARIANT}

COPY . .

RUN GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 go build -trimpath \
     -ldflags="-X 'main.version=${VERSION}'" \
     -o "/app/envgen" ./cmd/envgen

FROM --platform=$BUILDPLATFORM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app

COPY --from=build /app/envgen /app/envgen

ENTRYPOINT ["/app/envgen"]