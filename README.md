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

## Rule: everything under BASE_PATH

**This repo has no application HTTP surface.** It is a controller-runtime
operator that reconciles cluster objects; it registers no routes, emits no
redirects and serves no assets, and `BASE_PATH` is unused - the fleet may still
inject it, and the manager ignores it.

The one exception is not an exception in practice: controller-runtime's manager
binds a probe server on `$PORT` serving `/healthz` and `/readyz`. Those paths
are owned by the library and are hit by the kubelet straight at the pod, not
through the fleet proxy, so they are correct unprefixed and must stay that way.

The rule applies to anything HTTP you add on top. The fleet serves apps behind a
proxy at `BASE_PATH=/direct/<agent>:<port>`, and that prefix is forwarded
**unchanged** - nginx does not strip it - so a handler registered at `/metrics`
or `/admin` is never reached through the proxy. If you add such a surface
(a metrics endpoint, an admission webhook's own UI, a debug page):

1. Add a `basePath()` helper that reads `BASE_PATH` and normalises it to `""`
   (standalone) or `/leading/no-trailing-slash`.
2. Mount every route on the group/router returned from it - never on the root
   router. Route patterns stay group-relative (`"/metrics"`, not the full
   prefixed path).
3. Prefix `basePath()` onto every absolute URL you hand a client: redirect
   targets, links and form actions in rendered HTML, and asset paths.

The gin, echo, fiber, chi and buffalo templates each carry a worked example of
this helper.
