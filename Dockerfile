# Estágio 1: Compilação do binário
FROM golang:1.27-alpine AS builder

# Definir o diretório de trabalho dentro do container
WORKDIR /app

# Copiar os ficheiros de dependências primeiro para aproveitar o cache do Docker
COPY go.mod go.sum ./
RUN go mod download

# Copiar todo o código-fonte do projeto para o container
COPY . .

# Compilar o binário estático apontando para a nova localização do main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o devices-api ./cmd/api/main.go

# Estágio 2: Ambiente de execução reduzido
FROM alpine:latest

# Instalar certificados de segurança caso a API precise de fazer pedidos HTTPS externos
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copiar o binário final compilado no estágio anterior
COPY --from=builder /app/devices-api .

# Expor a porta padrão da aplicação
EXPOSE 8080

# Comando para iniciar o servidor da API
CMD ["./devices-api"]
