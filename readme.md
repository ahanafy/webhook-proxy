# webhook-proxy


## Example webhook
```
$ curl --request POST \
  --url http://localhost:8080/sleuth \
  --data '{"name":"podinfo","namespace":"test","phase":"Succeeded","metadata":{"some":"message"}}'
```

# Debugging
## Install the latest delve release:

```
$ go install github.com/go-delve/delve/cmd/dlv@latest
```
