FROM --platform=$BUILDPLATFORM golang:1.20.6 AS backend_test
ARG TARGETOS TARGETARCH
WORKDIR /source
RUN dpkg --add-architecture amd64 \
    && apt-get update \
    && apt-get install -y --no-install-recommends gcc-x86-64-linux-gnu libc6-dev-amd64-cross
COPY ./backend /source
ENV CGO_ENABLED=1
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH
ENV CC=x86_64-linux-gnu-gcc
RUN go test -a -ldflags '-linkmode external -extldflags "-static"' ./...

FROM --platform=$BUILDPLATFORM backend_test AS backend_build
RUN go build -o opent1d -a -ldflags '-linkmode external -extldflags "-static"' .

FROM node:18.16.1-slim AS frontend_build
WORKDIR /source
COPY ./www /source
RUN corepack enable
RUN corepack prepare yarn@stable --activate
RUN yarn
RUN yarn build

FROM scratch
COPY --from=backend_build /source/opent1d /opent1d
COPY --from=frontend_build /source/dist /www
VOLUME [ "/data" ]
EXPOSE 8080
ENV OPENT1D_DBPATH=file:/data/opent1d.sqlite
ENV OPENT1D_DEBUG=false
ENTRYPOINT [ "/opent1d" ]