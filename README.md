# Kubectl Trace Operator

This is a simple operator to wrap [kubectl-trace][0], such that a user may
create TraceJob CRDs against a Kubernetes API server. It maps the user-provided
CRD directly into the existing [tracejob type][1]. The main advantage of this
approach is it enables a dynamic, programmatic method to create tracejobs
without external user action. For example, detection of breach of some metric
could trigger deeper analysis via creating a TraceJob CRD. 

Since we are running in a controller context I've made a slight modification and
returned the ConfigMap and Job so the controller can apply ownership labels for
easy garbage collection when the TraceJob CRD is deleted.

## Quick Start

Build the docker images for the controller manager. In bash:

```bash
git clone https://github.com/alexeldeib/trace-operator
cd trace-operator
IMG=your_name/trace-operator:devtag make docker-build docker-push
```

Generate and install the CRD definitions and the manager deployment:
```bash
IMG=your_name/trace-operator:devtag make deploy
```

Apply a sample tracejob. Here we use `biolatency.bt` [(original source)][2]  from Brendan Gregg's BPF
Perf Tools book.

```bash
# ./config/samples/observe_v1alpha1_tracejob.yaml contains biolatency.bt as its program
kubectl apply -f https://raw.githubusercontent.com/alexeldeib/trace-operator/master/config/samples/observe_v1alpha1_tracejob.yaml
```

Noticed that the corresponding job and pod have been created:

```bash
$ kubectl get job
NAME                                                 COMPLETIONS   DURATION   AGE
kubectl-trace-16f11776-5b39-494a-8341-718850061edc   0/1           68s        68s
$ kubectl get pod
NAME                                                       READY   STATUS    RESTARTS   AGE
kubectl-trace-16f11776-5b39-494a-8341-718850061edc-tvqxc   1/1     Running   0          70s
```

We can end the job and log the data by attaching to the pod and sending ctrl-c (TODO(ace): more appropriate way to do this for controllers?):
```bash
kubectl attach -it kubectl-trace-16f11776-5b39-494a-8341-718850061edc-tvqxc
Defaulting container name to kubectl-trace-16f11776-5b39-494a-8341-718850061edc.
Use 'kubectl describe pod/kubectl-trace-16f11776-5b39-494a-8341-718850061edc-tvqxc -n default' to see all of the containers in this pod.
If you don't see a command prompt, try pressing enter.

^C
first SIGINT received, now if your program had maps and did not free them it should print them out



@usecs:
[32, 64)               5 |                                                    |
[64, 128)             88 |@@@@@@@                                             |
[128, 256)           336 |@@@@@@@@@@@@@@@@@@@@@@@@@@@@@                       |
[256, 512)           446 |@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@             |
[512, 1K)            224 |@@@@@@@@@@@@@@@@@@@                                 |
[1K, 2K)             593 |@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@|
[2K, 4K)             228 |@@@@@@@@@@@@@@@@@@@                                 |
[4K, 8K)               4 |                                                    |
[8K, 16K)             51 |@@@@                                                |
[16K, 32K)           124 |@@@@@@@@@@                                          |
[32K, 64K)            11 |                                                    |
```

## Cleanup

After running the quickstart, deleting all operator related objects can be done
with:
```bash
# From the root of this repo
kustomize build ./config/default | kubectl delete -f -
```

## Implementation Details

Since the operator shells out to the the tracejob package from kubectl-trace
under the hood, it has some nice interoperability with the CLI. For example, a
job created via CRD will have a field `.status.id` which can be fed back into
kubectl-trace to attach (although it seems there's a potential issue where
instead of printing the maps the process gets killed first):

```bash
$ kubectl get tracejob -o jsonpath="{.items[0].status.id}" 
8384546d-b1a4-4024-bc75-6c3bdb539646
$ kubectl trace attach 8384546d-b1a4-4024-bc75-6c3bdb539646
^C
first SIGINT received, now if your program had maps and did not free them it should print them out


@usecs:
[256, 512)             2 |@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@|
[512, 1K)              1 |@@@@@@@@@@@@@@@@@@@@@@@@@@                          |
[1K, 2K)               0 |                                                    |
[2K, 4K)               0 |                                                    |
[4K, 8K)               0 |                                                    |
[8K, 16K)              1 |@@@@@@@@@@@@@@@@@@@@@@@@@@                          |
```

[0]: https://github.com/iovisor/kubectl-trace
[1]: https://github.com/iovisor/kubectl-trace/blob/master/pkg/tracejob/job.go
[2]: https://github.com/brendangregg/bpf-perf-tools-book/blob/master/originals/Ch09_Disks/biolatency.bt