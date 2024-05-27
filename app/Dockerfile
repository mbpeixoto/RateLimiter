# Dockerfile
FROM golang:1.22.2 as build
WORKDIR /app
COPY . .
WORKDIR /app/cmd
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ratelimiter

FROM scratch
WORKDIR /app
COPY --from=build /app/cmd/ratelimiter .

ENV LIMITE_REQUISICOES_TOKEN=10
ENV LIMITE_REQUISICOES_IP=10
ENV TEMPO=30

ENTRYPOINT [ "./ratelimiter" ]
