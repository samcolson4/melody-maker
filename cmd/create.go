package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"gitlab.com/gomidi/midi/writer"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new melody.",
	Long:  "Generates a series of midi outputs and writes them to a file.",
	Run: func(cmd *cobra.Command, args []string) {
		f := assignFlags(cmd, flags{})
		makeMidi(f)
	},
}

// TODO: Change to imported list of random words
var fileNames = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "y", "x", "z"}

// TODO: Put vars into struct?
var keys = []string{"A", "A#", "B", "C", "C#", "D", "D#", "E", "F", "F#", "G", "G#"}
var noteNames = []string{"A", "A#", "B", "C", "C#", "D", "D#", "E", "F", "F#", "G", "G#"}
var scales = map[string][]string{
	"major": {"W", "W", "H", "W", "W", "W", "H"},
	"minor": {"W", "H", "W", "W", "H", "W", "W"},
}
var breaks = []bool{true, true, false, false, false}

func init() {
	rootCmd.AddCommand(createCmd)

	// Flags
	createCmd.PersistentFlags().StringP("output-folder", "o", "midi", "Folder to output to.")
	createCmd.PersistentFlags().IntSliceP("midi-note-range", "r", []int{21, 108}, "The lowest & highest midi notes that could be included.")
	createCmd.PersistentFlags().Bool("two-octave-limit", false, "Use with low-note-midi to have a two note melody from a starting point.")
	createCmd.PersistentFlags().IntP("file-number", "n", 6, "The number of files to generate. Max 23.")
	createCmd.PersistentFlags().StringP("key", "k", "C", "The key of the melody.")
	createCmd.PersistentFlags().StringP("scale", "s", "major", "The scale to use when generating a melody.")
	createCmd.PersistentFlags().IntSliceP("sequence-length-range", "l", []int{10, 20}, "Min & max number of notes in a sequence.")
	createCmd.PersistentFlags().IntSliceP("note-length", "d", []int{300, 500}, "Min & max length of each note.")
	createCmd.PersistentFlags().IntSliceP("velocity-range", "v", []int{15, 110}, "Min & max note velocity.")
	createCmd.PersistentFlags().StringP("instrument", "i", "piano", "See docs for full list.")
	createCmd.PersistentFlags().IntSlice("gap-bars-range", []int{0, 0}, "Min & max gap between notes (bars).")
	createCmd.PersistentFlags().IntSlice("gap-num-range", []int{0, 2}, "Min & max gap between notes (num of denoms).")
	createCmd.PersistentFlags().IntSlice("gap-denom-range", []int{0, 8}, "Min & max gap between notes (denoms).")
}

func makeMidi(f flags) {
	f = setFlagValues(f)
	files := fileNames[0:f.FileNumber]

	dir := f.OutputFolder
	fmt.Printf("Writing to folder: '%s'.\n", dir)

	allowedNotes := findAllowedNotes(f)

	// Write each file
	for _, file := range files {
		filename := fmt.Sprintf("%s.mid", file)
		outputPath := filepath.Join(dir, filename)

		err := writer.WriteSMF(outputPath, 2, func(wr *writer.SMF) error {
			wr.SetChannel(uint8(f.Instrument))

			numberOfNotes := random(f.NoteNumberRange[0]+1, f.NoteNumberRange[1]+1)
			var d midiModels
			var noteValues []int

			for i := 1; i < numberOfNotes; i++ {
				note := generateNote(f, allowedNotes)
				d.midiData = append(d.midiData, midiModel{Note: uint32(note),
					Velocity: uint32(random(f.VelocityMin, f.VelocityMax)),
				})
			}

			forwardOnLastLoop := false
			for i := 1; i < len(d.midiData); i++ {
				writer.NoteOn(wr, uint8(d.midiData[i].Note), uint8(d.midiData[i].Velocity))
				wr.SetDelta(uint32(random(f.MinNoteLength, f.MaxNoteLength)))
				writer.NoteOff(wr, uint8(d.midiData[i].Note))

				breakEl := random(0, len(breaks)-1)
				if breaks[breakEl] && !forwardOnLastLoop {
					writer.Forward(wr, uint32(random(f.GapBarsRange[0], f.GapBarsRange[1])), uint32(random(f.GapNumRange[0], f.GapNumRange[1])), uint32(random(f.GapDenomRange[0], f.GapDenomRange[1])))
					forwardOnLastLoop = true
				} else {
					forwardOnLastLoop = false
				}

				noteValues = append(noteValues, int(d.midiData[i].Note))
			}

			stringNotes := midisToNotes(noteValues)
			fmt.Printf("%s: %s\n", filename, stringNotes)

			writer.EndOfTrack(wr)

			return nil
		})

		if err != nil {
			fmt.Printf("Could not write SMF file %v\n", f)
			return
		}
	}
}

func findAllowedNotes(f flags) []int32 {
	var allowedNotes []int32
	var rootNoteMidi int

	scale := setScale(f)
	rootNote := f.Key

	for i, note := range noteNames {
		if note == rootNote {
			rootNoteMidi = i + 21
		}
	}

	// Add a skip value to the start of the scale before iterating
	scale = append([]string{"skip"}, scale...)

	for i, step := range scale {
		if i == 0 && step == "skip" {
			allowedNotes = append(allowedNotes, int32(rootNoteMidi))
		} else if i != 0 && step == "W" {
			allowedNotes = append(allowedNotes, allowedNotes[len(allowedNotes)-1]+2)
		} else if i != 0 && step == "H" {
			allowedNotes = append(allowedNotes, allowedNotes[len(allowedNotes)-1]+1)
		}
	}

	var i int
	for i < 10 {
		i += 1
		for _, note := range allowedNotes {
			allowedNotes = append(allowedNotes, note+12)
		}
	}

	allowedNotes = removeDuplicateInt(allowedNotes)

	return allowedNotes
}

func generateNote(f flags, allowedNotes []int32) int {
	var noteAllowed bool
	potentialNote := random(f.NoteRange[0], f.NoteRange[1])

	for !noteAllowed {
		if containsInt32(allowedNotes, potentialNote) {
			noteAllowed = true
		} else {
			potentialNote = random(f.NoteRange[0], f.NoteRange[1])
		}
	}

	return potentialNote
}
