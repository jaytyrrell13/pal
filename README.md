# Pal

Pal is a command line tool to generate an alias to `cd` into each of your projects. An alias can also be generated to open the directory with your editor's CLI command, e.g. nvim, code, subl.

## Installation

### Homebrew

```shell
brew install jaytyrrell13/tap/pal
```

### Manual Installation

Download the latest release archive from the [releases](https://github.com/jaytyrrell13/pal/releases) page.

## Usage

To get started, you may execute Pal's `install` command to create the config file in `~/.config/pal/config.json`.

```shell
pal install
```

The `create` command will ask, "What type of alias would you like to create?," with three options, "Parent", "Directory", and "Action". Use "Parent" to create aliases for the child directories under a specific directory. Use "Directory" to create an alias for a single directory. Use "Action" to create an alias for a command or script (ex. alias "ll" for "ls -lah").

```shell
pal create
```

The `update` command allows updating the name or command for one or more aliases.

```shell
pal update
```

The `remove` command allows removing one or more aliases.

```shell
pal remove
```

The `list` command will print out all the aliases currently in `~/.config/pal/aliases`.

```shell
pal list
```

The `config` command allows listing or updating config values.

```shell
pal config list
```

```shell
pal config update
```

## Support

If you'd like to support the development of `pal`, you can [buy me a coffee](https://www.buymeacoffee.com/jaytyrrell).
