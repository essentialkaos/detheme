<p align="center"><a href="#readme"><img src=".github/images/card.svg"/></a></p>

<p align="center">
  <a href="https://kaos.sh/r/detheme"><img src="https://kaos.sh/r/detheme.svg" alt="GoReportCard" /></a>
  <a href="https://kaos.sh/l/detheme"><img src="https://kaos.sh/l/ca59d01e7d47014dbf4a.svg" alt="Code Climate Maintainability" /></a>
  <a href="https://kaos.sh/b/detheme"><img src="https://kaos.sh/b/b1fa2a1a-3bb3-431c-85c7-6f52cf53cd7d.svg" alt="Codebeat badge" /></a>
  <br/>
  <a href="https://kaos.sh/w/detheme/ci"><img src="https://kaos.sh/w/detheme/ci-push.svg" alt="GitHub Actions CI Status" /></a>
  <a href="https://kaos.sh/w/detheme/codeql"><img src="https://kaos.sh/w/detheme/codeql.svg" alt="GitHub Actions CodeQL Status" /></a>
  <a href="#license"><img src=".github/images/license.svg"/></a>
</p>

<p align="center"><a href="#installation">Installation</a> • <a href="#command-line-completion">Command-line completion</a> • <a href="#man-documentation">Man documentation</a> • <a href="#usage">Usage</a> • <a href="#ci-status">CI Status</a> • <a href="#contributing">Contributing</a> • <a href="#license">License</a></p>

<br/>

`detheme` is SublimeText color theme downgrader for converting `*.sublime-color-scheme` files to `*.tmTheme`.

**Known limitations:**

- HWB colors not supported;
- `blend()` adjuster not supported;
- `blenda()` adjuster not supported;
- `saturation()` adjuster not supported;
- `lightness()` adjuster not supported;
- `min-contrast()` adjuster not supported.

### Installation

#### From source

To build the `detheme` from scratch, make sure you have a working Go 1.21+ workspace (_[instructions](https://go.dev/doc/install)_), then:

```
go install github.com/essentialkaos/detheme@latest
```

#### Container Image

The latest version of `detheme` also available as container image on [GitHub Container Registry](https://kaos.sh/p/detheme) and [Docker Hub](https://kaos.sh/d/detheme):

```bash
podman run --rm -it ghcr.io/essentialkaos/detheme:latest
# or
docker run --rm -it ghcr.io/essentialkaos/detheme:latest
```

#### Prebuilt binaries

You can download prebuilt binaries for Linux and macOS from [EK Apps Repository](https://apps.kaos.st/detheme/latest):

```bash
bash <(curl -fsSL https://apps.kaos.st/get) detheme
```

### Command-line completion

You can generate completion for `bash`, `zsh` or `fish` shell.

Bash:
```bash
sudo detheme --completion=bash 1> /etc/bash_completion.d/detheme
```

ZSH:
```bash
sudo detheme --completion=zsh 1> /usr/share/zsh/site-functions/detheme
```

Fish:
```bash
sudo detheme --completion=fish 1> /usr/share/fish/vendor_completions.d/detheme.fish
```

### Man documentation

You can generate man page using next command:

```bash
detheme --generate-man | sudo gzip > /usr/share/man/man1/detheme.1.gz
```

### Usage

<p align="center"><img src=".github/images/usage.svg"/></p>

### CI Status

| Branch | Status |
|--------|----------|
| `master` | [![CI](https://kaos.sh/w/detheme/ci-push.svg?branch=master)](https://kaos.sh/w/detheme/ci-push?query=branch:master) |
| `develop` | [![CI](https://kaos.sh/w/detheme/ci-push.svg?branch=develop)](https://kaos.sh/w/detheme/ci-push?query=branch:develop) |

### Contributing

Before contributing to this project please read our [Contributing Guidelines](https://github.com/essentialkaos/contributing-guidelines#contributing-guidelines).

### License

[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0)

<p align="center"><a href="https://essentialkaos.com"><img src="https://gh.kaos.st/ekgh.svg"/></a></p>
