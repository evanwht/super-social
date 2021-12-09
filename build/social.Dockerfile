FROM golang:1.17 as builder

WORKDIR /superhuman-social
COPY . .
ENV GOOS=linux
ENV GOARCH=amd64
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -v ./cmd/social

# Stage 3: Generate application image.
FROM gcr.io/google-appengine/debian9:latest

COPY --from=builder /superhuman-social/social ./

CMD ["./social"]
