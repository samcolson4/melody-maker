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
	Run: func(cmd *cobra.Command, args []string) {
		var f flags

		outputFolder, _ := cmd.Flags().GetString("output-folder")
		if outputFolder != "" {
			f.OutputFolder = outputFolder
		}

		lowNoteMidi, err := cmd.Flags().GetInt("low-note-midi")
		if lowNoteMidi != 0 && err == nil {
			f.LowNoteMidi = lowNoteMidi
		}

		highNoteMidi, err := cmd.Flags().GetInt("high-note-midi")
		if highNoteMidi != 0 && err == nil {
			f.HighNoteMidi = highNoteMidi
		}

		makeMidi(f)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.PersistentFlags().String("output-folder", "midi", "Folder to output to.")
	createCmd.PersistentFlags().Int("low-note-midi", 21, "The lowest note of the output, in midi format.")
	createCmd.PersistentFlags().Int("high-note-midi", 108, "The highest note of the output, in midi format.")
}

type flags struct {
	OutputFolder string
	LowNoteMidi  int
	HighNoteMidi int
}

type midiModel struct {
	Note     uint32
	Velocity uint32
}

type midiModels struct {
	midiData []midiModel
}

func makeMidi(f flags) {
	dir := f.OutputFolder
	files := []string{"a", "b", "c", "d", "e"}

	fmt.Printf("Writing to folder: '%s'.\n", dir)

	for _, file := range files {
		filename := fmt.Sprintf("%s.mid", file)
		outputPath := filepath.Join(dir, filename)

		err := writer.WriteSMF(outputPath, 2, func(wr *writer.SMF) error {
			wr.SetChannel(1) // sets the channel for the next messages

			noteNumber := random(3, 15)
			var d midiModels
			var noteValues []int
			sum := 0
			for i := 1; i < noteNumber; i++ {
				sum += i
				d.midiData = append(d.midiData, midiModel{Note: uint32(random(f.LowNoteMidi, f.HighNoteMidi)),
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

			fmt.Printf("Written file '%s': %s\n", filename, stringNotes)
			writer.EndOfTrack(wr)
			return nil
		})

		if err != nil {
			fmt.Printf("could not write SMF file %v\n", f)
			return
		}
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
