FROM golang:1.20.6 AS golang_common
ARG TARGETOS TARGETARCH
WORKDIR /source
COPY ./backend /source
ENV CGO_ENABLED=1
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH

FROM --platform=linux/amd64 golang_common AS backend_test_amd64
RUN dpkg --add-architecture amd64 \
    && apt-get update \
    && apt-get install -y --no-install-recommends gcc-x86-64-linux-gnu libc6-dev-amd64-cross
# RUN CC=x86_64-linux-gnu-gcc go test -a -ldflags '-linkmode external -extldflags "-static"' ./...
RUN CC=x86_64-linux-gnu-gcc go test ./...

FROM --platform=linux/arm64 golang_common AS backend_test_arm64
# RUN go test -a -ldflags '-linkmode external -extldflags "-static"' ./...
RUN go test ./...

FROM --platform=linux/amd64 backend_test_amd64 AS backend_build_amd64
# RUN CC=x86_64-linux-gnu-gcc go build -o opent1d -a -ldflags '-linkmode external -extldflags "-static"' .
RUN CC=x86_64-linux-gnu-gcc go build -o opent1d .

FROM --platform=linux/arm64 backend_test_arm64 AS backend_build_arm64
# RUN go build -o opent1d -a -ldflags '-linkmode external -extldflags "-static"' .
RUN go build -o opent1d .

FROM node:18.16.1-slim AS frontend_build
WORKDIR /source
COPY ./www /source
RUN corepack enable
RUN corepack prepare yarn@stable --activate
RUN yarn
RUN yarn build

FROM alpine:3.18.2 AS common_final
COPY --from=frontend_build /source/dist /www
VOLUME [ "/data" ]
EXPOSE 8080
ENV OPENT1D_DBPATH=file:/data/opent1d.sqlite
ENV OPENT1D_LOGLEVEL=error
ENTRYPOINT [ "/opent1d" ]

FROM --platform=linux/amd64 common_final
COPY --from=backend_build_amd64 /source/opent1d /opent1d

FROM --platform=linux/arm64 common_final
COPY --from=backend_build_arm64 /source/opent1d /opent1d