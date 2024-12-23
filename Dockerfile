FROM golang:1.23 AS build

RUN apt-get update && apt-get install -y make

COPY . /build
WORKDIR /build

ENV CGO_ENABLED=0
RUN make all

FROM gcr.io/distroless/static-debian12

COPY --from=build /build/dodfile-syntax /bin/dodfile-syntax
ENTRYPOINT ["/bin/dodfile-syntax"]
