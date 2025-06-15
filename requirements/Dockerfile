# Build Stage
FROM golang:1.24.3 AS builder
WORKDIR /app/project
COPY ./project/. ./
RUN go build -o forum main.go

# Final Stage
FROM debian:stable-slim
RUN apt-get update && apt-get upgrade -y
WORKDIR /project

# bring in the compiled binary & assets
COPY --from=builder /app/project/forum     .
COPY ./project/static   ./static
COPY ./project/templates ./templates


# seed DB in the image (optional if you always bind-mount)
# COPY ./project/database ./database

EXPOSE 8080

CMD ["./forum"]
