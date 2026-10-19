fmt:
    go fmt ./...

lint:
    golangci-lint run

test:
    go test ./...

run:
    go run ./cmd/rrl

build:
    CGO_ENABLED=0 go build -o build/rrl -trimpath -ldflags="-s -w" ./cmd/rrl/

clean:
    rm -rf ./build

install:
    CGO_ENABLED=0 go build -o build/rrl -trimpath -ldflags="-s -w" ./cmd/rrl/
    cp ./build/rrl ~/go/bin/

size:
    du -sh ./build/rrl
