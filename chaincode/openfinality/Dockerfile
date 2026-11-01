FROM golang:1.26 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -o openfinality .

FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/openfinality .
EXPOSE 9999
CMD ["./openfinality"]
