
# Build frontend.
FROM denoland/deno:2.8.2 as builder-deno

WORKDIR /app

COPY ./frontend/deno.lock ./frontend/package.json ./
RUN deno ci --prod

COPY ./frontend .
RUN deno task build

# Build server (and embed frontend).
FROM golang:1.26.4-alpine3.23 as builder-go

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download -x

COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./frontend/frontend.go ./frontend/frontend.go
COPY --from=builder-deno /app/build /app/frontend/build

ENV CGO_ENABLED=0
RUN go build -v -ldflags="-s -w -buildid=" -trimpath -o /app/Linea ./cmd/Linea
RUN go build -v -ldflags="-s -w -buildid=" -trimpath -o /app/linea-manager ./cmd/linea-manager

# Run server.
FROM gcr.io/distroless/static-debian13

COPY --from=builder-go /app/Linea /Linea
COPY --from=builder-go /app/linea-manager /linea-manager

ENV LINEA_ADDR=":9090"
ENV LINEA_MONGODB_DATABASE="linea-platform-v1"

EXPOSE 9090

USER nonroot:nonroot

ENTRYPOINT ["/Linea"]
