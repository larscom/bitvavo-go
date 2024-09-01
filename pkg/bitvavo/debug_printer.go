package bitvavo

import (
	"fmt"
	"log"
)

type DebugPrinter interface {
	Println(value ...any)
}

type DefaultDebugPrinter struct{}

func NewDefaultDebugPrinter() *DefaultDebugPrinter {
	return &DefaultDebugPrinter{}
}

func (l *DefaultDebugPrinter) Println(value ...any) {
	log.Println("[DEBUG bitvavo-go]", fmt.Sprint(value...))
}

func debug(p DebugPrinter, value ...any) {
	if p != nil {
		p.Println(value...)
	}
}
