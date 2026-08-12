# gio-term

A lightweight terminal emulator for wayland.

---

## TO DO

- Alt mode

---

## Prerequisites

* [Go](https://go.dev/doc/install) 1.26.5 or higher
* `make` utility (pre-installed on Linux/macOS, available via WSL or Chocolatey on Windows)
* vulkan-headers, pkg-config, libxkbcommon, libGL, libX11/Xcursor/Xfixes

---

## Quick Start

### Clone and Run

```bash
# Clone 
git clone https://github.com/ruzen42/gio-term # or git@github.com:ruzen42/gio-term for ssh
cd gio-term

# Build and run the project
make build
sudo make install # for cp /usr/local/bin 
gio-term
```

---

## Available Commands

All primary development tasks are managed via `make`:

| Target | Command | Description |
| --- | --- | --- |
| `all` | `make` or `make all` | Default target. Builds the project binary. |
| `build` | `make build` | Runs `go fmt`, creates `output/`, and compiles `output/main`. |
| `run` | `make run` | Compiles the binary and immediately executes it. |
| `fmt` | `make fmt` | Runs `go fmt` on the codebase to enforce standard formatting. |
| `install` | `make install` | Builds the binary and copies it to `/usr/local/bin/main`. |
| `clean` | `make clean` | Removes the `output/` build directory and compiled binary. |

---

## Build Customization

You can customize the build parameters directly in the `Makefile` or override them at runtime:

```makefile
BUILD_DIR=output
VERSION=0.1.0
NAME=main

```

### Overriding Variables at Runtime

```bash
# Build with a custom version tag
make build VERSION=1.2.3

# Build with a custom output binary name
make build NAME=myapp

```

---


---

## License

This template is licensed under the [Unlicense]([https://unlicense.org/]).
