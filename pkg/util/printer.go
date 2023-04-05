package util

import (
	"fmt"
	"io"
)

type Printer struct {
	w io.Writer
}

func NewPrinter(w io.Writer) *Printer {
	return &Printer{w: w}
}

func (p *Printer) print(str string) {
	_, _ = p.w.Write([]byte(str))
}

func (p *Printer) Print(a ...any) {
	p.print(fmt.Sprint(a...))
}

func (p *Printer) Printf(format string, a ...any) {
	p.Print(fmt.Sprintf(format, a...))
}

func (p *Printer) Println(a ...any) {
	p.Print(fmt.Sprintln(a...))
}
