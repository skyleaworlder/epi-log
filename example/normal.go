package main

import (
	"github.com/skyleaworlder/epilog"
)

func main() {
	go epilog.Use()
	defer epilog.End()

	epilog.Infoln("Hello! This is Infoln.")
	epilog.Warningln("WTF? Warning! Warning!")
	epilog.Infoln("Aha, another Infoln.")
	epilog.Debugln("Debugln, but this will be ignored...")
	epilog.Infoln("A third Infoln~")

	// time.Sleep(time.Second / 2)
}
