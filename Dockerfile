#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
ENV GO111MODULE=on
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o birthday-service -v

#final stage
FROM golang:alpine
# RUN apk --no-cache add curl
COPY --from=builder /app/birthday-service /birthday-service
EXPOSE 50054
ENTRYPOINT ["/birthday-service"]