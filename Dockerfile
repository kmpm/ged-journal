FROM golang:1.21
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN CGO_ENABLED=0 go build -o ./main cmd/ged-journal


FROM scratch
WORKDIR /app
COPY --from=builder /app/main ./
ENTRYPOINT ["./main"]
