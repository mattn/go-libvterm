# go-libvterm

Go binding to libvterm

## Usage

```go
vt := vterm.New(25, 80)
defer vt.Close()

vt.SetUTF8(true)

screen := vt.ObtainScreen()
screen.Reset(true)

_, err := vt.Write([]byte("\033[31mHello \033[32mGolang\033[0m"))
if err != nil {
	log.Fatal(err)
}
screen.Flush()

cell, err := screen.GetCellAt(0, 0)
```

## Requirements

libvterm: must be installed with pkg-config

## Installation

```
$ go get github.com/mattn/go-libvterm
```

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a. mattn)
