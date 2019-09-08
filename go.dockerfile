WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN go build .
COPY . .
EXPOSE 8000
