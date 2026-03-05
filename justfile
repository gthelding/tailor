# List available recipes
default:
    @just --list

# Alter tailor swatches
alter:
    @tailor alter

# Build tailor binary
build:
    @go build -ldflags "-s -w" -o tailor ./cmd/tailor

# Run tests
test:
    @go test ./...

# Check what tailor would change and measure
measure:
    @tailor baste
    @tailor measure
