fmt:
    go fmt ./...

lint:
    golangci-lint run

test:
    go test ./...

run *ARGS:
    go run ./cmd/rrl {{ARGS}}

build:
    CGO_ENABLED=0 go build -o build/rrl -trimpath -ldflags="-s -w -X github.com/piyushdugawa/rarrel/internal/rrl/cmd.Version=0.0.1-alpha.1" ./cmd/rrl/

clean:
    rm -rf ./build

install:
    CGO_ENABLED=0 go build -o build/rrl -trimpath -ldflags="-s -w -X github.com/piyushdugawa/rarrel/internal/rrl/cmd.Version=0.0.1-alpha.1" ./cmd/rrl/
    cp ./build/rrl ~/go/bin/

size:
    du -sh ./build/rrl
