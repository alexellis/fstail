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

### Installation

```bash
go install github.com/alexellis/fstail
```

If you need Go on a Linux system:

```
curl -sLS https://get.arkade.dev | sudo sh
sudo arkade system install go
```

## License

Copyright Alex Ellis 2023

MIT license.
