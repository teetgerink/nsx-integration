ARG VERSION=1.17

FROM golang:${VERSION}-alpine AS builder

# See https://github.com/hadolint/hadolint/wiki/DL4006
SHELL ["/bin/ash", "-eo", "pipefail", "-c"]

# Hadolint expects you to specify exact versions. Since this is a build container ignore this check
# hadolint ignore=DL3018
RUN apk add --update --no-cache \
    # Git is required for fetching the dependencies.
    git \
    # For CGO_ENABLED=1 we need these gcc and alpine SDK packages
    gcc alpine-sdk

ARG UID=10001
ARG USER=appuser

# See https://stackoverflow.com/a/55757473/12429735
# Hadolint expects you to specify exact versions. Since this is a build container ignore this check
# hadolint ignore=DL3018
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR $GOPATH/src/package

COPY . .





RUN go get -v \
    # Get swag command
    && go get -u github.com/swaggo/swag/cmd/swag


# Generate swagger docs
# Hadolint check SC2046 is a false positive
# hadolint ignore=SC2046
RUN swag init --output swagger --parseDependency --parseDepth=4 \
    # Run tests, we need CGO_ENABLED=1 for sqlite \
    # We also use '-short' to prevent the postgresql tests from being run. These can not run on Alpine.
    && CGO_ENABLED=1 go test ./... -short \
    # 🦢 Build the binary.
    && GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/nsx-api

############################
# STEP 2 build a small image
############################
FROM scratch

# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group


# Default environment variables
ENV HTTP_SERVER_ADDRESS 0.0.0.0:8080
ENV GIN_MODE=release

# Copy our static executable.
COPY --from=builder /go/bin/nsx-api /go/bin/nsx-api

# Use an unprivileged user.
USER appuser:appuser

# Run the security-api binary.
ENTRYPOINT ["/go/bin/nsx-api"]
