# contest-util

This is a tool for programming contest.

## Feature

* output a new file added methods under `cutil` directory

  * to do this, you need to write tokens around the methods (see `example/reader.go`)

* if you're using mac, also copies to clipboard


## How to use

1. create `cutil` directory
  ```sh
  mkdir cutil
  ```
2. create and write libraries under `cutil` directory (see `example/reader.go`)
3. create main.go
  ```sh
  touch main.go
  ```
4. write code for contest (see `main.go.example`)
5. build `export/main.go`
  ```sh
  go build ./export/main.go
  ```
6. exec main
  ```sh
  ./main
  ```
7. you can see a new file under `output` directory
