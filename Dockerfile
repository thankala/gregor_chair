# First stage: build the application
FROM golang:1.22.4-alpine AS builder

ARG SERVICE_NAME

ARG SSH_PRIVATE_KEY

ENV GO111MODULE=on

ENV GOOS=linux

ENV GOARCH=amd64

# Set the Current Working Directory inside the container
WORKDIR /build

# Install git and openssh.
RUN apk add --no-cache git openssh-client gcc musl-dev

# Set .gitconfig to use ssh instead of https
RUN git config --global url.ssh://git@github.com/.insteadOf https://github.com/

# Make the root foler for our ssh
#RUN mkdir -p /root/.ssh && \
#  chmod 0700 /root/.ssh && \
#  ssh-keyscan github.com > /root/.ssh/known_hosts

# Copy go mod and sum files
COPY go.mod go.sum ./

# Add your SSH key and download all dependencies. Dependencies will be
# cached if the go.mod and go.sum files are not changed
#RUN echo "$SSH_PRIVATE_KEY" > /root/.ssh/id_ed25519 && \
#  chmod 600 /root/.ssh/id_ed25519 && \
#  go mod download && \
#  rm -rf /root/.ssh/

RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN GOOS=$GOOS GOARCH=$GOARCH go build -tags musl -ldflags '-extldflags "-static"' -o $SERVICE_NAME .

# Second stage: prepare the final image
FROM alpine:latest

ARG SERVICE_NAME

# Set the argument value as an environment variable
ENV SERVICE_NAME=$SERVICE_NAME

RUN apk --no-cache add ca-certificates curl

WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /build/$SERVICE_NAME .

# Ensure executable permissions
RUN chmod +x ./$SERVICE_NAME

# Command to run the executable
CMD ./$SERVICE_NAME