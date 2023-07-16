alias t := test
alias b := build

# Shows help
default:
    @just --list --justfile {{ justfile() }}

# Builds pcd binary
build:
    go build -o pcd -ldflags "-s" cmd/pcd/main.go

# Run the tests
test:
    go test -v -race ./...
