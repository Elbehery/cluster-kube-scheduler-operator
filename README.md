# Kubernetes Scheduler operator

The Kubernetes Scheduler operator manages and updates the [Kubernetes Scheduler](https://github.com/kubernetes/kubernetes) deployed on top of
[OpenShift](https://openshift.io). The operator is based on OpenShift [library-go](https://github.com/openshift/library-go) framework and it
is installed via [Cluster Version Operator](https://github.com/openshift/cluster-version-operator) (CVO).

It contains the following components:

* Operator
* Bootstrap manifest renderer
* Installer based on static pods
* Configuration observer

By default, the operator exposes [Prometheus](https://prometheus.io) metrics via `metrics` service.
The metrics are collected from following components:

* Kubernetes Scheduler operator


## Configuration

The configuration for the Kubernetes Scheduler is the result of merging:

* a [default config](https://github.com/openshift/cluster-kube-scheduler-operator/blob/master/bindata/v4.1.0/kube-scheduler/defaultconfig-postbootstrap.yaml)
* an observed config (compare observed values above) from the spec `schedulers.config.openshift.io`.

All of these are sparse configurations, i.e. unvalidated json snippets which are merged in order to form a valid configuration at the end.


## Debugging

Operator also expose events that can help debugging issues. To get operator events, run following command:

```
$ oc get events -n  openshift-cluster-kube-scheduler-operator
```

This operator is configured via [`KubeScheduler`](https://github.com/openshift/api/blob/master/operator/v1/types_scheduler.go#L12) custom resource:

```
$ oc describe kubescheduler
```
```yaml
apiVersion: operator.openshift.io/v1
kind: KubeScheduler
metadata:
  name: cluster
spec:
  managementState: Managed
  ...
```
```
$ oc explain kubescheduler
```
to learn more about the resource itself.

The current operator status is reported using the `ClusterOperator` resource. To get the current status you can run follow command:

```
$ oc get clusteroperator/kube-scheduler
```

## Developing and debugging the operator

In the running cluster [cluster-version-operator](https://github.com/openshift/cluster-version-operator/) is responsible
for maintaining functioning and non-altered elements.  In that case to be able to use custom operator image one has to
perform one of these operations:

1. Set your operator in umanaged state, see [here](https://github.com/openshift/cluster-version-operator/blob/master/docs/dev/clusterversion.md) for details, in short:

```
oc patch clusterversion/version --type='merge' -p "$(cat <<- EOF
spec:
  overrides:
  - group: apps/v1
    kind: Deployment
    name: kube-scheduler-operator
    namespace: openshift-kube-scheduler-operator
    unmanaged: true
EOF
)"
```

2. Scale down cluster-version-operator:

```
oc scale --replicas=0 deploy/cluster-version-operator -n openshift-cluster-version
```

IMPORTANT: This apprach disables cluster-version-operator completly, whereas previous only tells it to not manage a kube-scheduler-operator!

After doing this you can now change the image of the operator to the desired one:

```
oc patch pod/kube-scheduler-operator-<rand_digits> -n openshift-kube-scheduler-operator -p '{"spec":{"containers":[{"name":"kube-scheduler-operator","image":"<user>/cluster-kube-scheduler-operator"}]}}'
```

## Developing and debugging the bootkube bootstrap phase

The operator image version used by the [installer](https://github.com/openshift/installer/blob/master/pkg/asset/ignition/bootstrap/) bootstrap phase can be overridden by creating a custom origin-release image pointing to the developer's operator `:latest` image:

```
$ IMAGE_ORG=<user> make images
$ docker push <user>/origin-cluster-kube-scheduler-operator

$ cd ../cluster-kube-apiserver-operator
$ IMAGES=cluster-kube-scheduler-operator IMAGE_ORG=<user> make origin-release
$ docker push <user>/origin-release:latest

$ cd ../installer
$ OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE=docker.io/<user>/origin-release:latest bin/openshift-install cluster ...
```
