# Build Stage
FROM golang:1.22.2 AS builder
WORKDIR /project
COPY project /project
RUN go build -o app main.go

# Final Stage
FROM debian:bullseye-slim
WORKDIR /project
COPY --from=builder /project/app /project/app
EXPOSE 8080
LABEL owner="soubaaiss" \
      owner="mlarbi" \
      owner="ebelfkih" \
      project="forum" \
      version="1.0" \
      description="This is a simple web application that works like a forum." \
      environment="zone 01"
CMD ["./app"]
