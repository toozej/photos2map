# setup project and deps
FROM golang:1.23-bullseye AS init

WORKDIR /go/golang-starter/

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

RUN CGO_ENABLED=0 go build -ldflags="${LDFLAGS}"

# runtime image
FROM scratch
# Copy our static executable.
COPY --from=build /go/golang-starter/golang-starter /go/bin/golang-starter
# Run the binary.
ENTRYPOINT ["/go/bin/golang-starter"]
