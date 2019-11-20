############################
# STEP 1 build executable binary
############################

FROM golang:alpine AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

RUN mkdir -p /go/raedarAPI/backend

WORKDIR /go/raedarAPI/backend

COPY . .

# Fetch dependencies.
# Using go get.
RUN go get -d -v

RUN go install raedarAPI

# Using go mod.
# RUN go mod download
# RUN go mod verify

# Build the binary.
RUN go build -o /go/bin/api

############################
# STEP 2 build a small image
############################

FROM scratch

# Copy our static executable.
COPY --from=builder /go/bin/api /go/bin/api

# Run the hello binary.
ENTRYPOINT ["/go/bin/api"]

EXPOSE 8000
