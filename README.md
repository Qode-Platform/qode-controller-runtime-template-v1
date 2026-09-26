# controller-runtime template

Provisioned from [`Qode-Platform/fleet-template-v1`](https://github.com/Qode-Platform/fleet-template-v1) - the fleet
lifecycle contract with a controller-runtime starter on top.

## Verified

Built and tested locally on Go 1.23.4 (toolchain auto-upgraded to 1.25):
`go build ./...` and `go test ./...` both pass.

## Fleet lifecycle

| step | command |
|---|---|
| install | `go mod download` |
| build | `go build -o ./.bin/app ./cmd/app` |
| start | `(none - not a service)` |

## Notes

- NOT RUNNABLE OUTSIDE A CLUSTER: the manager needs a kubeconfig or an in-cluster service account, so START_CMD is empty.
- Under k8s the manager serves /healthz and /readyz on $PORT - wire HEALTH_PATH back up if you deploy it as a workload.
