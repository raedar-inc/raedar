############################
# STEP 1 build executable binary
############################
# Accept the Go version for the image to be set as a build argument.
# Default to Go 1.13
ARG GO_VERSION=1.13

FROM golang:${GO_VERSION}-alpine AS builder

# We create an /app directory within our
# image that will hold our application source
# files
RUN mkdir /app

# Create the user and group files that will be used in the running container to
# run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

# Install git.
# Git is required for fetching the dependencies.
# Allow Go to retrieve the dependencies for the buld
RUN apk update && apk add --no-cache ca-certificates git

# Force the go compiler to use modules 
ENV GO111MODULE=on

# We copy everything in the root directory
# into our /app directory
ADD . /app

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Compile the binary, we don't want to run the cgo
# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/main cmd/app/main.go

# RUN CGO_ENABLED=0 GOOS=darwin -e GOARCH=amd64 -v go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main .

############################
# STEP 2 build a small image
############################

# Final stage: the running container.
FROM scratch AS final

WORKDIR /root/

# Import the user and group files from the first stage.
COPY --from=builder /user/group /user/passwd /etc/

# Import the Certificate-Authority certificates for enabling HTTPS.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/


# Import the compiled executable from the first stage.
COPY --from=builder /app .

# Declare the port on which the webserver will be exposed. ==> 8080
# As we're going to run the executable as an unprivileged user, we can't bind
# to ports below 1024.
EXPOSE 8080

# Perform any further action as an unprivileged user.
USER nobody:nobody

# Run the compiled binary.
ENTRYPOINT ["./bin/main"]
