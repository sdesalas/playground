FROM golang:1.21.8-bullseye AS builder

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY *.go /app/
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/service
RUN useradd -ms /bin/bash golang && chown -R golang /app 

# Second stage --> slimmer image
FROM golang:1.21.8-bullseye

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/service /app/service

WORKDIR /app
RUN chown -R golang /app

# Run as user, allowing write to db
USER golang
RUN touch /app/books.db

CMD ["/app/service"]
