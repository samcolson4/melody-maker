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

		lowNoteMidi, err := cmd.Flags().GetInt("low-note-midi")
		if err == nil {
			f.LowNoteMidi = lowNoteMidi
		}

		highNoteMidi, err := cmd.Flags().GetInt("high-note-midi")
		if err == nil {
			f.HighNoteMidi = highNoteMidi
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

		maxNotes, _ := cmd.Flags().GetInt("max-notes")
		if err == nil {
			f.MaxNotes = maxNotes
		}

		minNotes, _ := cmd.Flags().GetInt("min-notes")
		if err == nil {
			f.MinNotes = minNotes
		}

		if f.MinNotes > f.MaxNotes {
			log.Fatalf("Minimum notes cannot be larger than maximum notes.")
		}

		makeMidi(f)
	},
}

var fileNames = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "y", "x", "z"}

func init() {
	rootCmd.AddCommand(createCmd)

	// Flags
	createCmd.PersistentFlags().String("output-folder", "midi", "Folder to output to.")
	createCmd.PersistentFlags().Int("low-note-midi", 21, "The lowest note of the output, in midi format.")
	createCmd.PersistentFlags().Int("high-note-midi", 108, "The highest note of the output, in midi format.")
	createCmd.PersistentFlags().Bool("two-octave-limit", false, "Use with low-note-midi to have a two note melody from a starting point. Overwrites `--high-note-midi`.")
	createCmd.PersistentFlags().Int("file-number", 6, "The number of files to generate. Max 23.")
	createCmd.PersistentFlags().String("key", "C", "The key of the melody.")
	createCmd.PersistentFlags().String("scale", "major", "The scale to use when generating a melody.")
	// TODO Set logic default instrument (0 is piano etc)

	createCmd.PersistentFlags().Int("max-notes", 10, "The maximum number of notes in a sequence.")
	createCmd.PersistentFlags().Int("min-notes", 4, "The minimum number of notes in a sequence.")
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
			wr.SetChannel(0)

			numberOfNotes := random(f.MinNotes+1, f.MaxNotes+1)
			var d midiModels
			var noteValues []int

			sum := 0
			for i := 1; i < numberOfNotes; i++ {
				sum += i

				note := generateNote(f, allowedNotes)

				d.midiData = append(d.midiData, midiModel{Note: uint32(note),
					Velocity: uint32(random(40, 100)),
				})
			}

			for _, midi := range d.midiData {
				writer.NoteOn(wr, uint8(midi.Note), uint8(midi.Velocity))
				wr.SetDelta(uint32(random(300, 1500)))
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
