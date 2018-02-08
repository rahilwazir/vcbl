# vcbl

> CLI tool to fetch definitions from Vocabulary.com

## Install

```
$ go get github.com/rahilwazir/vcbl
```

Make sure `$GOPATH/bin` is exported

## Usage

```
$ vcbl tool

A tool is an instrument that you use to help you accomplish some task. If you are going to build a bookcase, you'll need the proper tools, like a saw, a drill, and a tape measure.
1. an implement used in the practice of a vocationwork with a tool
2. drive
3. obscene terms for penis
```

## Flags

```
--desc value, -d value  Description type of the lookup word. Possible values are: short, long, both (default: "short")
--suggestions, -s       Shows suggestion for similar words
--play, -p              Play the word pronounciation with SoX cli. SoX must be installed
--verbose               Debug output
--help, -h              show help
--version, -v           print the version
```

Run `vcbl --help` for a list of further options

## License

MIT © [Rahil Wazir](https://github.com/rahilwazir)
