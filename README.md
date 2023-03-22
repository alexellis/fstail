## fstail - tail files in adirectory including new files

Unfortunately, `tail -f /logs/*` may not do what you want it to do.

Bash will expand `*` to all existing files within `/logs/` and then show the extra lines added to each of them.

**What if you want to include newly created files too?**

`fstail` uses the [fsnotify mechanism](https://docs.huihoo.com/doxygen/linux/kernel/3.7/include_2linux_2fsnotify_8h_source.html) to detect both new and existing files, then starts printing their contents as and when they are written to.

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
