FROM golang:1.24

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
  && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN cp ./config/docker.toml ./config/default.toml

EXPOSE 8080

CMD ["go", "run", "./cmd/api"]
