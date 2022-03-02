package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

type flags struct {
	OutputFolder    string
	NoteRange       []int
	TwoOctaveLimit  bool
	FileNumber      int
	Key             string
	Scale           string
	NoteNumberRange []int
	MaxNoteLength   int
	MinNoteLength   int
	VelocityMax     int
	VelocityMin     int
	Instrument      int
	GapBarsRange    []int
	GapNumRange     []int
	GapDenomRange   []int
	Ascending       bool
	Descending      bool
	IntervalRange   []int
}

func assignFlags(cmd *cobra.Command, f flags) flags {
	outputFolder, _ := cmd.Flags().GetString("output-folder")
	if outputFolder != "" {
		f.OutputFolder = outputFolder
	}

	noteRange, err := cmd.Flags().GetIntSlice("midi-note-range")
	if err == nil {
		f.NoteRange = []int{noteRange[0], noteRange[1]}
	}

	twoOctaveLimit, err := cmd.Flags().GetBool("two-octave-limit")
	if err == nil {
		f.TwoOctaveLimit = twoOctaveLimit
	}

	ascending, err := cmd.Flags().GetBool("ascending")
	if err == nil {
		f.Ascending = ascending
	}

	descending, err := cmd.Flags().GetBool("descending")
	if err == nil {
		f.Descending = descending
	}

	if f.Ascending && f.Descending {
		log.Fatalf("Cannot set both ascending and descending flags at the same time.")
	}

	intervalRange, err := cmd.Flags().GetIntSlice("interval")
	if err == nil {
		f.IntervalRange = intervalRange
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
		f.NoteNumberRange = []int{sequenceLengthRange[0], sequenceLengthRange[1]}

		if f.NoteNumberRange[0] > f.NoteNumberRange[1] {
			swap := f.NoteNumberRange[0]
			f.NoteNumberRange[0] = f.NoteNumberRange[1]
			f.NoteNumberRange[1] = swap
		}
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

	gapBarsRange, err := cmd.Flags().GetIntSlice("gap-bars-range")
	if err == nil {
		f.GapBarsRange = []int{gapBarsRange[0], gapBarsRange[1]}
	}

	gapNumRange, err := cmd.Flags().GetIntSlice("gap-num-range")
	if err == nil {
		f.GapNumRange = []int{gapNumRange[0], gapNumRange[1]}
	}

	gapDenomRange, err := cmd.Flags().GetIntSlice("gap-denom-range")
	if err == nil {
		f.GapDenomRange = []int{gapDenomRange[0], gapDenomRange[1]}
	}

	return f
}
