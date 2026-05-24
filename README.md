# goRPC

A small Discord RPC (OpenASAR) thing written in Go!

### Install

You can download the repository to use it or download the compiled binaries in the [Releases tab](https://github.com/Lukiiy/goRPC/releases)!

### Usage

You can use `rpc [args]` or `go run rpc.go [args]`

```bash
rpc -client-id [clientid] -details "Playing Game" -state "Level 2"
```

```bash
rpc -client-id [clientid] -details "Playing Game" -large-image game_logo -large-text "Game Logo"
```

```bash
rpc -client-id [clientid] -details "Playing Game" -state "Ranked Match" -large-image game_logo -large-text "Game Logo" -small-image ranked_icon -small-text "Ranked logo"
```
