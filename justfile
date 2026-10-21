
name := "Rarrel"
cmd := "rrl"
build_dir := "build"
version := "v0.0.1-alpha.2"
version_path := "github.com/piyushdugawa/rarrel/internal/rrl/cmd.Version"

go_home := env("GO_HOME", "~/go")
go_bin := go_home / "bin"

fmt:
    go fmt ./...

lint:
    golangci-lint run

test:
    go test ./...

run *ARGS:
    go run ./cmd/rrl {{ARGS}}

build:
    CGO_ENABLED=0 go build -o {{build_dir}}/{{cmd}} -trimpath -ldflags="-s -w -X {{version_path}}={{version}}" ./cmd/rrl/

clean:
    rm -rf {{build_dir}}

install: build
    cp ./{{build_dir}}/{{cmd}} {{go_bin}}

size:
    du -sh ./{{build_dir}}/{{cmd}}
