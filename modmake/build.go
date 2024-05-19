package main

import (
	. "github.com/saylorsolutions/modmake"
)

func main() {
	b := NewBuild()
	b.Test().Does(Go().TestAll())

	b.Execute()
}
