# ======= 開発用ベース（ソース入り） =======
FROM golang:1.24 AS local

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

CMD [ "go", "run", "main.go" ]

# ======= 本番用ビルドステージ =======
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

# ======= 実行用イメージ（最小限） =======
FROM debian:bullseye-slim AS prod

WORKDIR /app

# 必要なツール（gitとか）を明示的に入れたければここで
# RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# builder からビルド済みバイナリのみをコピー
COPY --from=builder /app/main .

CMD ["./main"]