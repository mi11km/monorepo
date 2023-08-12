FROM golang:1.21.0-alpine as builder
ENV HOME /app
WORKDIR $HOME
COPY go.mod go.sum $HOME
RUN go mod download
COPY . $HOME
RUN go build -o /bin/app ./cmd/main.go
EXPOSE 8080
CMD ["/bin/app"]

FROM alpine:latest as prd
RUN adduser -D app
COPY --from=builder --chown=app:app /bin/app /app
EXPOSE 8080
USER app
ENTRYPOINT ["/app"]
