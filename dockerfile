# -------- build stage --------
FROM golang:1.23-alpine AS build
WORKDIR /src

# cache deps
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY . .

# build binary (sesuaikan path main.go kamu)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /out ./app/main.go


FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=build /out /app/server

# Railway kasih PORT via env
ENV PORT=8080
EXPOSE 8080
CMD ["/app/server"]
