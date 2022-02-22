package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
	"gitlab.com/gomidi/midi/writer"
)

type flags struct {
	OutputFolder   string
	LowNoteMidi    int
	HighNoteMidi   int
	TwoOctaveLimit bool
	FileNumber     int
	Key            string
	Scale          string
	MaxNotes       int
	MinNotes       int
	MaxNoteLength  int
	MinNoteLength  int
	VelocityMax    int
	VelocityMin    int
	Instrument     int
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new melody.",
	Long:  "Generates a series of midi outputs and writes them to a file.",
	Run: func(cmd *cobra.Command, args []string) {
		var f flags

		outputFolder, _ := cmd.Flags().GetString("output-folder")
		if outputFolder != "" {
			f.OutputFolder = outputFolder
		}

		noteRange, err := cmd.Flags().GetIntSlice("midi-note-range")
		if err == nil {
			f.LowNoteMidi = noteRange[0]
			f.HighNoteMidi = noteRange[1]
		}

		twoOctaveLimit, err := cmd.Flags().GetBool("two-octave-limit")
		if err == nil {
			f.TwoOctaveLimit = twoOctaveLimit
		}

		fileNumber, err := cmd.Flags().GetInt("file-number")
		if err == nil {
			f.FileNumber = fileNumber
		}

		key, _ := cmd.Flags().GetString("key")
		if err == nil {
			if containsString(keys, key) {
				f.Key = key
			} else {
				log.Fatalf("Key '%s' does not exist", key)
			}
		}

		scale, _ := cmd.Flags().GetString("scale")
		if err == nil {
			f.Scale = scale
		}

		sequenceLengthRange, err := cmd.Flags().GetIntSlice("sequence-length-range")
		if err == nil {
			f.MinNotes = sequenceLengthRange[0]
			f.MaxNotes = sequenceLengthRange[1]
		}

		if f.MinNotes > f.MaxNotes {
			log.Fatalf("Minimum number of notes cannot be larger than the maximum number.")
		}

		noteLength, err := cmd.Flags().GetIntSlice("note-length")
		if err == nil {
			f.MinNoteLength = noteLength[0]
			f.MaxNoteLength = noteLength[1]
		}

		velocityRange, err := cmd.Flags().GetIntSlice("velocity-range")
		if err == nil {
			f.VelocityMin = velocityRange[0]
			f.VelocityMax = velocityRange[1]
		}

		instrument, err := cmd.Flags().GetString("instrument")
		if err == nil {
			switch instrument {
			case "piano":
				f.Instrument = 0
			case "synth":
				f.Instrument = 1
			case "bass":
				f.Instrument = 2
			case "pluck-synth":
				f.Instrument = 3
			case "strings":
				f.Instrument = 4
			case "session-strings":
				f.Instrument = 5
			case "brass":
				f.Instrument = 6
			case "trumpet":
				f.Instrument = 7
			case "edrums":
				f.Instrument = 8
			case "drums":
				f.Instrument = 9
			case "organ":
				f.Instrument = 10
			case "e-piano":
				f.Instrument = 11
			case "synth-strings":
				f.Instrument = 12
			case "analog-synth":
				f.Instrument = 13
			case "synth-brass":
				f.Instrument = 14
			case "sculpture-synth":
				f.Instrument = 15
			}
		}

		makeMidi(f)
	},
}

var fileNames = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "y", "x", "z"}

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

			numberOfNotes := random(f.MinNotes+1, f.MaxNotes+1)
			var d midiModels
			var noteValues []int

			sum := 0
			for i := 1; i < numberOfNotes; i++ {
				sum += i

				note := generateNote(f, allowedNotes)

				d.midiData = append(d.midiData, midiModel{Note: uint32(note),
					Velocity: uint32(random(f.VelocityMin, f.VelocityMax)),
				})
			}

			for _, midi := range d.midiData {
				writer.NoteOn(wr, uint8(midi.Note), uint8(midi.Velocity))
				wr.SetDelta(uint32(random(f.MinNoteLength, f.MaxNoteLength)))
				writer.NoteOff(wr, uint8(midi.Note))

				noteValues = append(noteValues, int(midi.Note))
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

	// Add all octaves of each note
	// TODO: Fix - sometimes hangs.
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

var keys = []string{"A", "A#", "B", "C", "C#", "D", "D#", "E", "F", "F#", "G", "G#"}
var noteNames = []string{"A", "A#", "B", "C", "C#", "D", "D#", "E", "F", "F#", "G", "G#"}
var scales = map[string][]string{
	"major": {"W", "W", "H", "W", "W", "W", "H"},
	"minor": {"W", "H", "W", "W", "H", "W", "W"},
}

func generateNote(f flags, allowedNotes []int32) int {
	var noteAllowed bool
	potentialNote := random(f.LowNoteMidi, f.HighNoteMidi)

	for !noteAllowed {
		if containsInt32(allowedNotes, potentialNote) {
			noteAllowed = true
		} else {
			potentialNote = random(f.LowNoteMidi, f.HighNoteMidi)
		}
	}

	return potentialNote
}
