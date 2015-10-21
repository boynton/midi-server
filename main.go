package main

import (
	"github.com/boynton/ell"
	midi "github.com/boynton/midi-ell"
	"os"
)

func main() {
	ell.SetFlags(true, false, false, false, false)
	ell.Init(new(midi.Extension))
	pianoPath := os.Getenv("GOPATH") + "/src/github.com/boynton/piano-server"
	ell.AddEllDirectory(pianoPath)
	ell.Run("piano-server")
}
