# Accept the Go version for the image to be set as a build argument.
# Default to Go 1.13
ARG GO_VERSION=1.13
# First stage: build the executable.
FROM golang:${GO_VERSION}-alpine AS builder

ARG APP_NAME
RUN echo "Building for ${APP_NAME}"

# Create the user and group files that will be used in the running container to
# run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

# Install the Certificate-Authority certificates for the app to be able to make
# calls to HTTPS endpoints.
# Git is required for fetching the dependencies.
RUN apk add --no-cache \ 
    ca-certificates \
    git


# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /src

# Fetch dependencies first; they are less susceptible to change on every build
# and will therefore be cached for speeding up the next build
COPY ./go.mod ./go.sum ./
RUN go mod download

# Import the code from the context.
COPY ./ ./

#COPY ./certs/ /etc/ssl/certs/

#RUN go get ./...
# Build the executable to `/app`. Mark the build as statically linked.
RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /app \
    ./cmd/${APP_NAME}

# Final stage: the running container.
FROM scratch AS final

# Import the user and group files from the first stage.
COPY --from=builder /user/group /user/passwd /etc/

# Import the Certificate-Authority certificates for enabling HTTPS.
#COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Import the GCLoud credentials
#COPY --from=builder --chown=nobody:nobody /etc/ssl/certs/google-key.json /etc/ssl/certs/

# Import the compiled executable from the first stage.
COPY --from=builder /app /app

# ENV GRPC_GO_REQUIRE_HANDSHAKE=off
# ENV GOOGLE_APPLICATION_CREDENTIALS=/etc/ssl/certs/google-key.json

# Perform any further action as an unprivileged user.
USER nobody:nobody


# Run the compiled binary.
ENTRYPOINT ["/app"]