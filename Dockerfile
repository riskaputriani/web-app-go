FROM golang:1.22-alpine AS build

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/web-app .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /app/web-app ./web-app
COPY --from=build /app/templates ./templates
RUN adduser -D -g '' appuser
USER appuser
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/app/web-app"]
