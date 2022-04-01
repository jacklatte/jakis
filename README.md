# jakis

`jekis` is a simple redis implementation, written in `go`.

## How to use

1. Create server on port 8080

```
go run cmd/jakis-server/main.go
```

2. Connect to server

```
nc 127.0.0.1 8080
```

3. Type Command

```
> set 1 100
OK
> get 1
100
> exit
BYE
```
