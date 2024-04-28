# build stage
FROM golang:1.22-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./


# test stage it's going to skip if DOCKER_BUILDKIT is enabled
FROM build-stage AS test-stage
RUN go test -v ./...


# release stage
FROM gcr.io/distroless/base-debian12 AS release-stage

WORKDIR /

COPY --from=build-stage /server /server

USER nonroot:nonroot

ENTRYPOINT ["/server"]