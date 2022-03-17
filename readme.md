# webhook-proxy


## Example webhook
```
$ curl --request POST \
  --url http://localhost:8080/sleuth \
  --data '{"name":"podinfo","namespace":"test","phase":"Succeeded","metadata":{"k8sDeployName":"podinfo", "sleuthDeployName": "podinfo"}}'
```

| Name                                | Description                            |
| ----------------------------------- | -------------------------------------- |
| `name`                              | name of the deployment                 |
| `namespace`                         | namespace the deployment is running in |
| `metadata` (optional)               | can include custom key/value pairs     |
| `metadata.k8sDeployName` (optional) | alternative name for lookup in kubernetes api     |
| `metadata.sleuthDeployName` (optional) | alternative name to map to Sleuth.io `Code Deployments`     |

# Debugging
## Install the latest delve release:

```
$ go install github.com/go-delve/delve/cmd/dlv@latest
```
