
FROM golang:1.26-alpine AS builder


WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/org-api ./cmd/main.go


FROM alpine:latest

WORKDIR /app


COPY --from=builder /app/bin/org-api .


COPY --from=builder /app/migrations ./migrations


EXPOSE 6969



CMD ["./org-api", "-dbserver=postgres_db", "-dbuser=postgres", "-dbpass=password", "-dbname=test_department", "-dbport=5432"]