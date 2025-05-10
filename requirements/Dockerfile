# Build Stage
FROM golang:1.23.0 AS builder
WORKDIR /app/project
COPY ./project/. ./
RUN go build -o forum main.go

# Final Stage
FROM debian:stable-slim
RUN apt-get update && apt-get upgrade -y
WORKDIR /project
COPY --from=builder ./app/project/forum .
COPY ./project/static ./static
# COPY ./project/database ./database
COPY ./project/templates ./templates
EXPOSE 8080
LABEL owner="soubaaiss" \
      owner="mlarbi" \
      owner="ebelfkih" \
      project="forum" \
      version="1.0" \
      description="This is a simple web application that works like a forum." \
      environment="zone 01"
CMD ["./forum"]
