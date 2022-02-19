package cmd

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/go-audio/midi"
	"github.com/spf13/cobra"
	"gitlab.com/gomidi/midi/reader"
	"gitlab.com/gomidi/midi/writer"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new melody.",
	Long:  "Generates a new series of midi outputs and writes to a specified file.",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.Flags().String("folder", "", "Output folder for the midi.")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	makeMidi()
}

type midiModel struct {
	Note     uint32
	Velocity uint32
}

type midiModels struct {
	midiData []midiModel
}

func makeMidi() {
	dir := "midi"
	f := filepath.Join(dir, "new.mid")

	err := writer.WriteSMF(f, 2, func(wr *writer.SMF) error {

		wr.SetChannel(1) // sets the channel for the next messages

		noteNumber := random(4, 15)
		var d midiModels
		var noteValues []int
		sum := 0
		for i := 1; i < noteNumber; i++ {
			sum += i
			d.midiData = append(d.midiData, midiModel{Note: uint32(random(35, 100)),
				Velocity: uint32(random(40, 100)),
			})
		}

		for _, midi := range d.midiData {
			writer.NoteOn(wr, uint8(midi.Note), uint8(midi.Velocity))
			wr.SetDelta(uint32(random(300, 1500)))
			writer.NoteOff(wr, uint8(midi.Note))

			noteValues = append(noteValues, int(midi.Note))
		}

		stringNotes := midiToNote(noteValues)

		fmt.Printf("Written: %s\n", stringNotes)
		writer.EndOfTrack(wr)
		return nil
	})

	if err != nil {
		fmt.Printf("could not write SMF file %v\n", f)
		return
	}
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

func midiToNote(midiIn []int) []string {
	var notes []string

	for _, key := range midiIn {
		notes = append(notes, midi.NoteToName(key))
	}

	return notes
}
