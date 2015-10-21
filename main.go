package main

import (
	"github.com/boynton/ell"
	"os"
)

type MidiEnabled struct {
}

func (ps *MidiEnabled) Init() error {
	pianoPath := os.Getenv("GOPATH") + "/src/github.com/boynton/piano-server"
	ell.AddEllDirectory(pianoPath)
	return initMidi()
}

func (ps *MidiEnabled) Cleanup() {
	midiClose(nil)
}

func main() {
   ell.SetFlags(true, false, false, false, false)
	ell.Init(new(MidiEnabled))
	ell.Run("piano")
}
