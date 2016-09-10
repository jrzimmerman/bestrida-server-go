# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.6

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/jrzimmerman/bestrida-server-go

# Build the bestrida-server-go command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go install github.com/jrzimmerman/bestrida-server-go

# Run the bestrida-server-go command by default when the container starts.
ENTRYPOINT /go/bin/bestrida-server-go

# Add environment variables
ENV PORT=4001

# Document that the service listens on port 4001.
EXPOSE 4001