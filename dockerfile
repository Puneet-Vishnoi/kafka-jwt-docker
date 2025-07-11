FROM golang:1.23.2-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

RUN apk update && apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/main .

#EXPOSE 8080

ENTRYPOINT ["./main"]



# FROM golang:1.23.2-alpine AS builder

# ENV CGO_ENABLED=0 \
#     GOOS=linux \
#     GOARCH=amd64

# WORKDIR /app

# # Copy go mod files
# COPY go.mod go.sum ./

# # Run go mod tidy to ensure modules are up-to-date
# RUN go mod tidy

# # Download dependencies
# RUN go mod download

# # Copy the rest of the application
# COPY . .

# # Build the Go binary
# RUN go build -o main ./cmd

# # Final image
# FROM alpine:latest

# WORKDIR /root/

# # Copy the built binary from the builder stage
# COPY --from=builder /app/main .

# # Expose port (optional)
# # EXPOSE 8080

# # Set the entrypoint to run the Go binary
# ENTRYPOINT ["./main"]




# # Use a Windows-based image
# FROM mcr.microsoft.com/windows/servercore:ltsc2022 AS builder

# ENV CGO_ENABLED=0 \
#     GOOS=windows \
#     GOARCH=amd64

# WORKDIR /app

# COPY go.mod go.sum ./

# RUN go mod download

# COPY . .

# RUN go build -o main.exe ./cmd

# # Final image for Windows
# FROM mcr.microsoft.com/windows/servercore:ltsc2022

# WORKDIR /root/

# COPY --from=builder /app/main.exe .

# ENTRYPOINT ["C:\\root\\main.exe"]  # Use the full path for Windows executables
