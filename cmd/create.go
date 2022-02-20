package cmd

import (
	"fmt"
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
		if lowNoteMidi != 0 && err == nil {
			f.LowNoteMidi = lowNoteMidi
		}

		highNoteMidi, err := cmd.Flags().GetInt("high-note-midi")
		f.HighNoteMidi = highNoteMidi

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
			f.Key = key
		}

		scale, _ := cmd.Flags().GetString("scale")
		if err == nil {
			f.Scale = scale
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
	// TODO Set channel
	// TODO # notes in sequence
}

func makeMidi(f flags) {
	f = setFlagValues(f)
	files := fileNames[0:f.FileNumber]

	dir := f.OutputFolder
	fmt.Printf("Writing to folder: '%s'.\n", dir)

	allowedNotes := findAllowedNotes(f)

	// Write to each file
	for _, file := range files {
		filename := fmt.Sprintf("%s.mid", file)
		outputPath := filepath.Join(dir, filename)

		err := writer.WriteSMF(outputPath, 2, func(wr *writer.SMF) error {
			wr.SetChannel(1) // sets the channel for the next messages

			numberOfNotes := random(3, 15)
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
	var scale []string

	if f.Scale == "Major" {
		scale = majorScale
	} else {
		scale = majorScale
	}

	rootNote := f.Key
	var rootNoteMidi int

	for i, note := range noteNames {
		if note == rootNote {
			rootNoteMidi = i + 21
		}
	}

	for i, step := range scale {
		if i == 0 {
			allowedNotes = append(allowedNotes, int32(rootNoteMidi))
		} else if i != 0 && step == "W" {
			allowedNotes = append(allowedNotes, allowedNotes[len(allowedNotes)-1]+2)
		} else if i != 0 && step == "H" {
			allowedNotes = append(allowedNotes, allowedNotes[len(allowedNotes)-1]+1)
		}
	}

	for i, note := range allowedNotes {
		for i < 40 {
			i += 1
			allowedNotes = append(allowedNotes, note+12)
		}
	}

	allowedNotes = removeDuplicateInt(allowedNotes)

	return allowedNotes
}

var noteNames = []string{"A", "A#", "B", "C", "C#", "D", "D#", "E", "F", "F#", "G", "G#"}
var majorScale = []string{"W", "W", "W", "H", "W", "W", "W", "H"}

func generateNote(f flags, allowedNotes []int32) int {
	var noteIsAllowed bool

	randomEl := random(0, len(allowedNotes)-1)
	potentialNote := random(f.LowNoteMidi, f.HighNoteMidi)

	for noteIsAllowed == false {
		if potentialNote == int(allowedNotes[randomEl]) {
			noteIsAllowed = true
		} else {
			potentialNote = random(f.LowNoteMidi, f.HighNoteMidi)
			noteIsAllowed = false
		}
	}

	return potentialNote
}
