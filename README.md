# go pointer analysis bug

Steps to reproduce:

1. Install dependencies.

    In the directory where the repository is cloned, run:

    ```bash
    mkdir gopath
    env GOPATH=$PWD/gopath GO111MODULE=off go get github.com/docker/docker \
        golang.org/x/tools/go/{pointer,packages}
    ```

2. Run the program that calls the pointer analysis.

    ```bash
    env GOPATH=$PWD/gopath GO111MODULE=off go run main.go
    ```
