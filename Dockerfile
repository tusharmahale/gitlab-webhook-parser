FROM golang:1.20 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
COPY src /app/src
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-w" -o /gitlab-webhook-parser

# Run the tests in the container
# FROM build-stage AS run-test-stage
# RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /gitlab-webhook-parser /gitlab-webhook-parser

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/gitlab-webhook-parser"]