# Pal

Pal is a command line tool to generate an alias to `cd` into each of your projects. An alias can also be generated to open the directory with your editor's cli command, e.g. nvim, code, subl.

## Installation

### Homebrew

```bash
brew install jaytyrrell13/tap/pal
```

### Manual Installation

Download the latest release archive from the [releases](https://github.com/jaytyrrell13/pal/releases) page.

## Usage

To get started, you may execute Pal's `install` command. This will ask for the path to your projects and your editor's cli tool e.g. nvim, code, subl. These settings will be saved in `~/.config/pal/config.json`.

```bash
pal install [--path | -p] [--editorCmd | -e]
```

The `make` command will go through each directory of your projects and ask for the alias you want to use. This will generate a `~/.pal` file that will need to be sourced from your `.zshrc` or `.bashrc` file.

```bash
pal make
```

The `add` command can be used if you want an alias for a directory outside of your projects. For example, a directory of notes.

```bash
pal add [--path | -p] [--name | -n]
```

The `list` command will print out all of the aliases currently in `~/.pal`.

```bash
pal list
```

The `clean` command will delete your `~/.pal` file in your home directory.

```bash
pal clean
```

The `refresh` command will delete your `~/.pal` file and then run the `make` command.

```bash
pal refresh
```

## Support

If you'd like to support the development of `pal`, you can [buy me a coffee](https://www.buymeacoffee.com/jaytyrrell).
