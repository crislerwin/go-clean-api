# Etapa 1: Builder
FROM golang:1.27-alpine AS builder

WORKDIR /app

# Instala dependências do sistema necessárias para compilar (se houver CGO)
# RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compila o binário. 
# CGO_ENABLED=0 garante um binário estático puro (ideal para rodar em alpine/scratch)
RUN CGO_ENABLED=0 GOOS=linux go build -o go-clean-api ./cmd/api/main.go

# Etapa 2: Runner (Imagem final leve)
FROM alpine:latest

WORKDIR /root/

# Copia o binário da etapa anterior
COPY --from=builder /app/go-clean-api .

EXPOSE 8080

CMD ["./go-clean-api"]