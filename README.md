# Bloosom
poc implimentaion of a raft cluster leveraging the dragonboat package

### Usage
in 3 seporate shells run:
-  `go run . 1`
-  `go run . 2`
-  `go run . 3`

### Notes
- this was worrysome from their examples:
```
// https://github.com/golang/go/issues/17393
if runtime.GOOS == "darwin" {
    signal.Ignore(syscall.Signal(0xd))
}
```