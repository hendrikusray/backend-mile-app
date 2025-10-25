# ---- build stage ----
FROM golang:1.23-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# build dari app/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o server app/main.go

# ---- run stage ----
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /app/server /app/server
COPY config.cold.json /app/config.cold.json
EXPOSE 8080
USER 65532:65532
ENTRYPOINT ["/app/server"]
