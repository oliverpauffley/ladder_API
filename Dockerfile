FROM golang as builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o /chess_ladder


FROM ubuntu
COPY --from=builder /chess_ladder /
COPY wait-for-postgres.sh /
RUN apt-get update \
    && apt-get install -y postgresql-client \
    && rm -rf /var/lib/apt/lists/*
EXPOSE 8000
CMD ["./wait-for-postgres.sh", "db", "/chess_ladder"]
