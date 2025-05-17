FROM golang:1.23-bookworm AS builder
WORKDIR /app
# Copiar los archivos del proyecto
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Compilamos la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd

# Imagen final para ejecución
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/app .
EXPOSE 8080
CMD ["/root/app"]