# gitkeep
CLI utility to add and remove .gitkeep files to directories.

## run

    go run ./cmd/gitkeep/main.go /path/to/directory
    
    gitkeep /path/to/directory (default ".")


## maintenance

    go mod edit -go=1.23

Or manually change this line in go.mod:

    go 1.23

### upgrade dependencies

Update all the modules to their latest versions:

    go get -u ./...

Or to upgrade only direct dependencies:

    go get -u

Then clean up unused or redundant entries:

    go mod tidy

    go clean -modcache

## install

Assuming **GOPATH** and **GOBIN** is set.

    go install ./cmd/gitkeep

