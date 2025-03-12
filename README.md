## fstail - scan a directory for changed files and tail them

**What's this for?**

When you need to see the output from all changed files within a directory.

**Why doesn't `tail -f /var/logs/*` work?**

Unfortunately, `tail -f /logs/*` may not do what you want it to do. Bash will expand `*` to all existing files within `/logs/` and then show the extra lines added to each of them.

It also will not recurse down, any levels deeper than the current directory.

**How is fstail different then?**

`fstail` uses the [gopkg.in/fsnotify](https://pkg.go.dev/gopkg.in/fsnotify.v1@v1.4.7) to detect both new files, and existing files that are changed. It then starts concatenting their contents to the terminal.

I needed this for [actuated.dev](https://actuated.dev) which launches microVMs on servers for CI.

Each VM launched will create a different file at: `/var/log/actuated/GUID.txt`, and `tail -f *` would only find existing files. 

### Usage

Tail the current directory:

```
cd /var/log/actuated
fstail


[7cd0d139b24d9cd30e3ad9ce7afcfe5999d2bf20.txt] [   12.770522] bash[1191]: 2023-03-22 11:04:00Z: Running job: arkade-e2e (run-job)
[49c8f4be774730ff6e5070166fc34ac25dc0e320.txt] [   13.363398] bash[1183]: 2023-03-22 11:04:00Z: Running job: arkade-e2e (k3sup)
```

To turn off the file prefix, set `FS_PREFIX=0`.

```
[   12.770522] bash[1191]: 2023-03-22 11:04:00Z: Running job: arkade-e2e (run-job)
[   13.363398] bash[1183]: 2023-03-22 11:04:00Z: Running job: arkade-e2e (k3sup)
```

Tail files in a given directory:

```
fstail /var/log/nginx/
```

By default, the base filename is going to be printed as a prefix for each tailed file:

```bash
1e51959055fb132720d03584388b5ac738689798.txt | Booting Linux Kernel.. OK
13d461e989733fa0f75df5227debab4be3504726.txt | Shutting down in 30s... 
```

To suppress the prefix, run the command with `FS_PREFIX=0`.

### Tail container logs from Kubernetes

The kubelet with K3s will create log files for containers in `/var/log/containers` with long prefixes, the `FS_PREFIX=k8s` env-var can be used to redact the string.

Any text passed after the folder forms a grep for the filenames, so here, we are only looking at "queue-worker" containers.

Note that these files do not receive a "WRITE" fsnotify event, only an initial "CREATE".

```bash
FS_PREFIX=k8s sudo -E fstail /var/log/containers/ "queue-worker"

2025/03/12 07:39:23 Attaching to: queue-worker-6b579d58d7-krgx5_openfaas_queue-worker-08e45f83b4f6edbd9ae2182f6709ccfb28e3e5871b5cc0d303c20c4c7b2a3d56.log
2025/03/12 07:39:23 Attaching to: queue-worker-6b579d58d7-krgx5_openfaas_queue-worker-f03309cc46a9752a4366062da18872ee215cc1837dcb335fa0d339948f3d9609.log
queue-worker-6b579d58d7-krgx5| 2025-03-10T12:40:45.84048421Z stderr F 2025-03-10T12:40:45.838Z	info	jetstream-queue-worker/main.go:118	JetStream queue-worker	{"version": "0.3.46", "gitCommit": "8320975de6c98c8b6ef6781e2db610cc389c7e33"}
queue-worker-6b579d58d7-krgx5| 2025-03-10T12:40:45.840510746Z stderr F 2025-03-10T12:40:45.838Z	info	jetstream-queue-worker/main.go:120	Licensed to: Alex Ellis 
queue-worker-6b579d58d7-krgx5| 2025-03-10T12:40:45.840514335Z stderr F 
queue-worker-6b579d58d7-krgx5| 2025-03-10T12:40:45.840517017Z stderr F 2025-03-10T12:40:45.839Z	info	metrics/metrics.go:79	Starting metrics server on port 8081
```

### Installation

Download a release binary:

* [fstail releases](https://github.com/alexellis/fstail/releases/)

Or download via arkade:

```bash
curl -SLs https://get.arkade.dev | sh

arkade get fstail
```

Or install via Go to install from source:

```bash
go install github.com/alexellis/fstail
```

If you need Go on a Linux system:

```
curl -sLS https://get.arkade.dev | sudo sh
sudo arkade system install go
```

If you wish to build multi-arch binaries on your own machine:

```bash
make dist
```

## License

Copyright Alex Ellis 2023

MIT license.
