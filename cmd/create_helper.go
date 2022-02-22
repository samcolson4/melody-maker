package cmd

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/go-audio/midi"
	"gitlab.com/gomidi/midi/reader"
)

type midiModel struct {
	Note     uint32
	Velocity uint32
}

type midiModels struct {
	midiData []midiModel
}

type printer struct{}

func (pr printer) noteOn(p *reader.Position, channel, key, vel uint8) {
	fmt.Printf("Track: %v Pos: %v NoteOn (ch %v: key %v vel: %v)\n", p.Track, p.AbsoluteTicks, channel, key, vel)
}

func (pr printer) noteOff(p *reader.Position, channel, key, vel uint8) {
	fmt.Printf("Track: %v Pos: %v NoteOff (ch %v: key %v)\n", p.Track, p.AbsoluteTicks, channel, key)
}

func random(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}

func midisToNotes(midiIn []int) []string {
	var notes []string

	for _, key := range midiIn {
		notes = append(notes, midi.NoteToName(key))
	}

	return notes
}

func containsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func containsInt32(s []int32, e int) bool {
	for _, a := range s {
		if a == int32(e) {
			return true
		}
	}
	return false
}

func setFlagValues(f flags) flags {
	if f.TwoOctaveLimit {
		f.HighNoteMidi = f.LowNoteMidi + int(24)
	}

	f.Scale = strings.ToLower(f.Scale)

	return f
}

func removeDuplicateInt(intSlice []int32) []int32 {
	allKeys := make(map[int32]bool)
	list := []int32{}
	for _, item := range intSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func setScale(f flags) []string {
	var scale []string

	if f.Scale == "major" {
		scale = scales[f.Scale]
	} else if f.Scale == "minor" {
		scale = scales[f.Scale]
	} else {
		scale = scales["major"]
	}

	return scale
}
