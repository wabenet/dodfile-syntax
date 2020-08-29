FROM golang:1.14-alpine AS build

RUN apk add -U make

COPY . /build
WORKDIR /build

ENV CGO_ENABLED 0
RUN make all

FROM scratch

COPY --from=build /build/dodfile-syntax /bin/dodfile-syntax
ENTRYPOINT ["/bin/dodfile-syntax"]
