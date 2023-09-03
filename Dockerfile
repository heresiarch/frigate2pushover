# Use an official Go development image as the base
FROM golang:1.21.0 AS build
# Set the working directory
WORKDIR /app
# Copy the Go source code into the container
COPY main.go    .
COPY go.mod     .
COPY go.sum     .
# Build the Go program
RUN go build -o frigate2pushover main.go

# Build a very small container based on
# https://klotzandrew.com/blog/smallest-golang-docker-image
FROM scratch
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group
COPY --from=build /app/frigate2pushover .
USER small-user:small-user
CMD ["./frigate2pushover"]