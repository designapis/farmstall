FROM golang:alpine as builder
RUN apk update && apk add --no-cache \
  git
RUN mkdir /build

# Download dependencies
ADD ./go.mod /build/
ADD ./go.sum /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go get

# Build
ADD . /build/
RUN CGO_ENABLED=0 \
  GOOS=linux \
  go build \
  -a \
  -installsuffix cgo \
  -ldflags '-extldflags "-static"' \
  -o farmstall .

# Copy over artifacts
FROM alpine
COPY --from=builder /build/farmstall /app/
COPY ./openapi.yaml /app/
WORKDIR /app
ENV PORT 80
EXPOSE 80
CMD ["./farmstall"]
