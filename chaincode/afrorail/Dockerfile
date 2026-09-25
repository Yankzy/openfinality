FROM golang:1.22-bullseye AS builder
WORKDIR /app
COPY . .
RUN go build -o afrorail .

FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/afrorail .
EXPOSE 9999
CMD ["./afrorail"]
