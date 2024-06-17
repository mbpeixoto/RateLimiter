# Estágio de construção
FROM golang:1.22.2 AS build
WORKDIR /app
COPY . .
WORKDIR /app/cmd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main

# Estágio final
FROM scratch
WORKDIR /app
COPY --from=build /app/cmd/main .
COPY --from=build /app/cmd/.env .

ENTRYPOINT ["./main"]

