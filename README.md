# gbx, the Global Blackbox CLI

This project hosts the code for gbx, the [Global Blackbox](https://globalblackbox.io) command line interface (CLI).

# Installation

## Homebrew

gbx can be easily installed using homebrew:
```bash
brew install globalblackbox/tap/gbx
```

## Others

Alternatively, you can simply download the binary from the releases page and move it the most appropriate directory:
```bash
VERSION="0.1.5"
ARCH="amd64"
OS="linux"
wget https://github.com/globalblackbox/gbx/releases/download/v$VERSION/gbx_$VERSION\_$OS\_$ARCH.tar.gz -qO- | \
tar xvfz - gbx
mv ./gbx /usr/local/bin
```

and it should now be available:
```bash
gbx --help
GBX allows you to sign up, manage your account,
and interact with Global Blackbox services through a command-line interface.

Usage:
  gbx [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  logs        Retrieve and download logs from Global Blackbox
  sign-up     Sign up for a Global Blackbox account

Flags:
  -h, --help   help for gbx

Use "gbx [command] --help" for more information about a command.
```

# Features

gbx currently supports:

- Sign-up to Global Blackbox
- List probe failure log files per region, target domain and date
- Download log files for inspection

# Documentation

Full documentation for Global Blackbox can be found [here](https://globalblackbox.io/docs)
