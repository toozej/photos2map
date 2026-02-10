# setup project and deps
FROM golang:1.25-bookworm AS init

WORKDIR /go/photos2map/

COPY go.mod* go.sum* ./
RUN go mod download

COPY . ./

FROM init AS vet
RUN go vet ./...

# run tests
FROM init AS test
RUN go test -coverprofile c.out -v ./... && \
	echo "Statements missing coverage" && \
	grep -v -e " 1$" c.out

# build binary
FROM init AS build
ARG LDFLAGS

# Install coreutils for sleep and other utilities utilized in devcontainer
RUN apt-get update && apt-get install --no-install-recommends -y coreutils

RUN CGO_ENABLED=0 go build -ldflags="${LDFLAGS}"

# runtime image
FROM scratch
# Copy our static executable.
COPY --from=build /go/photos2map/photos2map /go/bin/photos2map
# Run the binary.
USER non-root
ENTRYPOINT ["/go/bin/photos2map"]
