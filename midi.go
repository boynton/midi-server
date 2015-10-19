/*
Copyright 2014 Lee Boynton

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"github.com/boynton/ell"
	"github.com/rakyll/portmidi"
	"sync"
	"time"
)

var inputKey = ell.Intern("input:")
var outputKey = ell.Intern("output:")
var bufsizeKey = ell.Intern("bufsize:")

func initMidi() error {
   ell.DefineFunctionKeyArgs("midi-open", midiOpen, ell.NullType,
		[]*ell.LOB{ell.StringType, ell.StringType, ell.NumberType},
		[]*ell.LOB{ell.EmptyString, ell.EmptyString, ell.NewInt(1024)},
		[]*ell.LOB{inputKey, outputKey, bufsizeKey})
	ell.DefineFunction("midi-write", midiWrite, ell.NullType, ell.NumberType, ell.NumberType, ell.NumberType)
	ell.DefineFunction("midi-close", midiClose, ell.NullType)
	ell.DefineFunction("midi-listen", midiListen, ell.ChannelType)
	return nil
}

var midiOpened = false
var midiInDevice string
var midiOutDevice string
var midiBufsize int64
var midiBaseTime float64

var midiOut *portmidi.Stream
var midiIn *portmidi.Stream
var midiChannel chan portmidi.Event
var midiMutex = &sync.Mutex{}

func findMidiInputDevice(name string) portmidi.DeviceId {
	devcount := portmidi.CountDevices()
	for i := 0; i < devcount; i++ {
		id := portmidi.DeviceId(i)
		info := portmidi.GetDeviceInfo(id)
		if info.IsInputAvailable {
			if info.Name == name {
				return id
			}
		}
	}
	return portmidi.DeviceId(-1)
}

func findMidiOutputDevice(name string) (portmidi.DeviceId, string) {
	devcount := portmidi.CountDevices()
	for i := 0; i < devcount; i++ {
		id := portmidi.DeviceId(i)
		info := portmidi.GetDeviceInfo(id)
		if info.IsOutputAvailable {
			if info.Name == name {
				return id, info.Name
			}
		}
	}
	id := portmidi.GetDefaultOutputDeviceId()
	info := portmidi.GetDeviceInfo(id)
	return id, info.Name
}

func midiOpen(argv []*ell.LOB) (*ell.LOB, error) {
//	defaultInput := "USB Oxygen 8 v2"
//	defaultOutput := "IAC Driver Bus 1"
	latency := int64(10)
	if !midiOpened {
		err := portmidi.Initialize()
		if err != nil {
			return nil, err
		}
		midiOpened = true
		midiInDevice = ell.StringValue(argv[0])
		midiOutDevice = ell.StringValue(argv[1])
		midiBufsize = ell.Int64Value(argv[2])

		outdev, outname := findMidiOutputDevice(midiOutDevice)
		out, err := portmidi.NewOutputStream(outdev, midiBufsize, latency)
		if err != nil {
			return nil, err
		}
		midiOut = out
		midiOutDevice = outname
		if midiInDevice != "" {
			indev := findMidiInputDevice(midiInDevice)
			if indev >= 0 {
				in, err := portmidi.NewInputStream(indev, midiBufsize)
				if err != nil {
					return nil, err
				}
				midiIn = in
			}
		}
		midiBaseTime = ell.Now()
		
	}
	result := ell.MakeStruct(4)
	if midiInDevice != "" {
		ell.Put(result, inputKey, ell.NewString(midiInDevice))
	}
	if midiOutDevice != "" {
		ell.Put(result, outputKey, ell.NewString(midiOutDevice))
	}
	ell.Put(result, bufsizeKey, ell.NewInt64(midiBufsize))
	return result, nil
}

func midiAllNotesOff() {
	midiOut.WriteShort(0xB0, 0x7B, 0x00)
}

func midiClose(argv []*ell.LOB) (*ell.LOB, error) {
	midiMutex.Lock()
	if midiOut != nil {
		midiAllNotesOff()
		midiOut.Close()
		midiOut = nil
	}
	midiMutex.Unlock()
	return ell.Null, nil
}

// (midi-write 144 60 80) -> middle C note on
// (midi-write 128 60 0) -> middle C note off
func midiWrite(argv []*ell.LOB) (*ell.LOB, error) {
	status := ell.Int64Value(argv[0])
	data1 := ell.Int64Value(argv[1])
	data2 := ell.Int64Value(argv[2])
	midiMutex.Lock()
	var err error
	if midiOut != nil {
		err = midiOut.WriteShort(status, data1, data2)
	}
	midiMutex.Unlock()
	return ell.Null, err
}

func midiListen(argv []*ell.LOB) (*ell.LOB, error) {
	ch := ell.Null
	midiMutex.Lock()
	if midiIn != nil {
		ch = ell.NewChannel(int(midiBufsize), "midi")
		go func(s *portmidi.Stream, ch *ell.LOB) {
			for {
				time.Sleep(10 * time.Millisecond)
				events, err := s.Read(1024)
				if err != nil {
					continue
				}
				channel := ell.ChannelHandle(ch)
				if channel != nil {
					for _, ev := range events {
						ts := (float64(ev.Timestamp) / 1000) + midiBaseTime
						st := ev.Status
						d1 := ev.Data1
						d2 := ev.Data2
						channel <- ell.List(ell.NewFloat64(ts), ell.NewInt64(st), ell.NewInt64(d1), ell.NewInt64(d2))
					}
				}
			}
		}(midiIn, ch)
	}
	midiMutex.Unlock()
	return ch, nil
}

