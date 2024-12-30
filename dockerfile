FROM golang:1.23.1 

# Set working directory
WORKDIR /app

COPY go.mod go.sum* ./
RUN go mod download 
COPY . .

RUN go build -o gotex ./cmd/gotex/...

# Keep the container running
CMD ["tail", "-f", "/dev/null"]
