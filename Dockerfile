FROM golang:1.24 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/mitchellh/gox@latest \
    && gox -osarch="linux/mips linux/mipsle linux/arm64" -output="dist/kvasx-{{.Arch}}" ./cmd/kvasx

FROM scratch AS artifacts
COPY --from=build /app/dist /dist
