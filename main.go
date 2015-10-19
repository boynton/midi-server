package main

import (
	"github.com/boynton/ell"
)

type MidiEnabled struct {
}

func (ps *MidiEnabled) Init() error {
	return initMidi()
}

func (ps *MidiEnabled) Cleanup() {
	midiClose(nil)
}

func main() {
	ell.Main(new(MidiEnabled))
}
