Conveyor usage example
======================

The simple `logtest.py` app just writes to 3 different log files (all having .log suffix) and at regular intervals,
it rotates the files (renames them to `<filename>.<number>`).

Docker image `wigwam/logtest` was created from this app, also `wigwam/conveyor` is an image for the conveyor app itself.

The `pod-manifest.yaml` demonstrates the behaviour: the pod will have two containers: `main` and `conveyor`, `main` writing
the log files and `conveyor` reading them and printing to stdout (with file name prefix).

After applying to your cluster, you can watch the stdout logs using:
```
kubectl logs -f logconveyor-example -c conveyor
```
