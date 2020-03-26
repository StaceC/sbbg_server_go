# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from golang v1.11 base image
FROM golang:1.11

# Add Maintainer Info
LABEL maintainer="Stace C <stace@hiddenfield.com>"

# Build Args
ARG HOME_DIR
ARG MODULE_NAME
ARG APP_NAME
ARG LOG_DIR_NAME
ARG VOLUME_DIR
ARG VERSION

RUN echo "Build number: $VERSION"

# Set the Current Working Directory inside the container
WORKDIR $MODULE_NAME/cmd/$APP_NAME

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Download all the dependencies
# https://stackoverflow.com/questions/28031603/what-do-three-dots-mean-in-go-command-line-invocations
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# Get gin. We all love gin. So, drink some gin
RUN go get github.com/codegangsta/gin

# Create Log Directory
# RUN mkdir -p ./$LOG_DIR_NAME

# This container exposes these port to the outside world
EXPOSE 8089

ENV LOG_FILE_LOCATION=$APP_NAME.log

# Run the executable
CMD ["sbbg"]
