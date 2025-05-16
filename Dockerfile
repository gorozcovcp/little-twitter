FROM golang:1.20 as builder

WORKDIR /app

# Copiar los archivos del proyecto
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compilamos la aplicación
RUN go build -o app

# Imagen final para ejecución
FROM gcr.io/distroless/base-debian11
WORKDIR /root/

COPY --from=builder /app/app .
CMD ["/root/app"]
