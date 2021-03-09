ARG BUILDER_IMAGE=golang:buster
ARG BASE_IMAGE=scratch

FROM $BUILDER_IMAGE AS build

WORKDIR /app

ADD . /app

ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0

RUN go mod download
RUN go mod verify
RUN go build -o goapp -ldflags '-w -s' main.go

# ---

FROM $BASE_IMAGE as final

WORKDIR /app

COPY --from=build /app/goapp ./
# COPY --from=build /app/config/ ./config/

ENTRYPOINT ["/app/goapp"]