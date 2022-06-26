# VMC (and OSC) parser for Go

This is a parser for the **V**irtual **M**otion **C**apture messages and, as the protocol uses OSC
as underlying format, a parser for **O**pen **S**ound **C**ontrol as well.

## Goals

The project was written to help my friend [@negue](https://github.com/negue) who wanted to simply
inspect and filter VMC messages before forwarding them to an app that actually handles the messages.
Existing solutions didn't really solve his problem, as they all bother with fully transforming and
handling the messages, which is what you need if you want to write your own VMC server.

Therefore, the project puts an especially high focus on only parsing incoming data, and avoiding
extra allocations where possible. Sadly, due to limitations of the Go language itself, it's not
possible to make it fully zero copy.

## Usage

Add the latest version to your project:

```sh
go get github.com/dnaka91/go-vmc
```

For further usage details and examples, please check out the
[API docs](https://pkg.go.dev/github.com/dnaka91/go-vmc).

## License

This project is licensed under [MIT License](LICENSE) (or <http://opensource.org/licenses/MIT>).
