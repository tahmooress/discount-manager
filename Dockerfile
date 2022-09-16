# Set the base image for subsequent instructions
# FROM golang as builder

# # Update packages
# RUN apt-get update

# RUN apt-get install -qq \
#     make \
#     curl 
#     # ca-certificates \
#     # git \
#     # libssl-dev
 

# # Clear out the local repository of retrieved package files
# RUN apt-get clean

# COPY . /src
# WORKDIR /src

# RUN make build

# ##
# # app
# ##
# FROM golang as app

# COPY --from=builder /bin /go/bin/app

# CMD ["ping","4.2.2.4"]



# Start from golang base image
FROM golang:alpine

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git && apk add --no-cach bash && apk add build-base

# Setup folders
RUN mkdir /src
WORKDIR /src

# Copy the source from the current directory to the working Directory inside the container
COPY . .
COPY .env .

# RUN export $(grep -v '^#' .env | xargs -d '\n')

# Download all the dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# Build the Go app
RUN make build

# Expose port 8080 to the outside world
EXPOSE 8080

CMD ["./bin/src-app"]
# CMD ["ping","192.168.1.1"]
