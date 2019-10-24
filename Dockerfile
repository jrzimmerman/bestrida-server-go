FROM golang:1.13 AS build

ENV GO111MODULE=on

# Copy the code from the host and compile it
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bestrida-server-go .

FROM scratch
COPY --from=build /bestrida-server-go .

# Service listens on port 4001.
EXPOSE 4001
ENTRYPOINT ["./bestrida-server-go"]