package cmd

import (
	"fmt"
	"path/filepath"

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
		fmt.Println("Generating melody...")
	},
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
}

func makeMidi() {
	dir := "./midi"
	f := filepath.Join(dir, "smf-test.mid")

	err := writer.WriteSMF(f, 2, func(wr *writer.SMF) error {

		wr.SetChannel(11) // sets the channel for the next messages
		writer.NoteOn(wr, 60, 100)
		wr.SetDelta(300)
		writer.NoteOff(wr, 120)

		writer.NoteOn(wr, 100, 100)
		wr.SetDelta(300)
		writer.NoteOff(wr, 120)

		// wr.SetDelta(240)
		// writer.NoteOn(wr, 125, 50)
		// wr.SetDelta(20)
		// writer.NoteOff(wr, 125)
		// writer.EndOfTrack(wr)

		// wr.SetChannel(2)
		// writer.NoteOn(wr, 120, 50)
		// wr.SetDelta(60)
		// writer.NoteOff(wr, 120)
		// writer.EndOfTrack(wr)
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
