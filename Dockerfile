FROM golang:1.23.4-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o main ./cmd/backend/main.go
RUN chmod +x main
EXPOSE 4040
CMD [ "./main" ]