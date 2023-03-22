## fstail - scan a directory for changed files and tail them

Unfortunately, `tail -f /logs/*` may not do what you want it to do.

Bash will expand `*` to all existing files within `/logs/` and then show the extra lines added to each of them.

`fstail` uses the [gopkg.in/fsnotify](https://pkg.go.dev/gopkg.in/fsnotify.v1@v1.4.7) to detect both new files, and existing files that are changed. It then starts concatenting their contents to the terminal.

I needed this for [actuated.dev](https://actuated.dev) which launches microVMs on servers for CI.

Each VM launched will create a different file at: `/var/log/actuated/GUID.txt`, and `tail -f *` would only find existing files. 

### Usage

Tail the current directory:

```
cd /var/log/
fstail
```

Tail files in a given directory:

```
/var/log/nginx/
```

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
