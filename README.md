# mgnl-workspace-copy

[![ci](https://github.com/quangthe/mgnl-workspace-copy/actions/workflows/build-docker.yaml/badge.svg)](https://github.com/quangthe/mgnl-workspace-copy/actions/workflows/build-docker.yaml)
[![Docker Stars](https://img.shields.io/docker/stars/pcloud/mgnl-workspace-copy.svg?style=flat)](https://hub.docker.com/r/pcloud/mgnl-workspace-copy/)
[![Docker Pulls](https://img.shields.io/docker/pulls/pcloud/mgnl-workspace-copy.svg?style=flat)](https://hub.docker.com/r/pcloud/mgnl-workspace-copy/)

A small utility to copy Magnolia workspace data from Postgres DB dump.

Supported Postgres version of DB dump file: `11`, `12`

> **Attention**: The data in the current workspaces will be deleted (truncated) when doing workspace copy. 
> Ensure you have all the backups of the running Postgres DB before performing the operation.  

Tool version and usage:
```shell
docker run --rm -it  pcloud/mgnl-workspace-copy -v
docker run --rm -it  pcloud/mgnl-workspace-copy copy --help
```

## Task: Copy workspace content from Postgres DB dump

Create K8S job to run workspace content copy

`copy-job.yaml`

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  generateName: workspace-copy-
spec:
  template:
    metadata:
      labels:
        app: workspace-copy
    spec:
      volumes:
        - name: data
          persistentVolumeClaim:
            # the volume contains Postgres DB dump
            claimName: db-dump-pvc
      restartPolicy: Never
      containers:
        - name: copy
          image: pcloud/mgnl-workspace-copy
          imagePullPolicy: IfNotPresent
          command:
            - /app
            - copy
            - --pgversion
            - "12"
            # CHANGEME: the k8s service dns name of the running Postgres DB
            - --host
            - author-db.dev
            # CHANGEME: the Postgres database name of the running Postgres DB
            - --dbname
            - author
            # CHANGEME: path to the Postgres DB dump file
            - --dump-file
            - /db/author.dump
            # CHANGEME: list of workspaces to copy content
            - --mgnl-workspaces
            - "campaigns,category,dam,messages"
            # copy datastore (true/false)
            - --copy-datastore=true
          volumeMounts:
            - name: data
              mountPath: /db
          resources:
            requests:
              memory: 64Mi
              cpu: 100m
```

Run the job in the namespace contains the `db-dump-pvc` volume to perform workspace data copy

```shell
kubectl -n <NAMESPACE> create copy-job.yaml
```