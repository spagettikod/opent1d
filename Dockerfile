FROM golang:1.20.5 AS test
WORKDIR /source
COPY . /source
RUN CGO_ENABLED=1 GOOS=linux go test -a -ldflags '-linkmode external -extldflags "-static"' -v ./...