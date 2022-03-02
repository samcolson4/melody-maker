# melody-maker
Generate fragments of midi data to help inspire melodies.

Built using the excellent [gomidi](https://github.com/gomidi/midi) and [go-audio/midi](https://github.com/go-audio/midi) packages.

## Usage
When more features have been added, a binary will be provided as part of a release. For now, building locally (`make build`) then examining the help options is the best approach.

`make buildRunDelete` is a helper used during development to create a binary with the latest local changes, generate output and then delete the binary without the manual steps.

### Example configurations
Short notes, set to a synth when imported directly to Logic Pro, with minimal breaks between notes.

```
$ melody-maker create \
  --note-length=200,300 \
  --sequence-length-range=10,30 \
  --midi-note-range=50,66 \
  --key G \
  --scale minor \
  --velocity-range 55,95 \
  --instrument analog-synth
```

Medium-length notes, soft key presses.
```
$ melody-maker create \
  --note-length=750,1500 \
  --sequence-length-range=10,30 \
  --midi-note-range=50,66 \
  --key G \
  --scale minor \
  --velocity-range 35,75
```

Longer notes, with significant rests in between, defaulting to the key of C Major.
```
$ melody-maker create \
  --note-length=2500,3500 \
  --sequence-length-range=10,30 \
  --midi-note-range=50,66 \
  --velocity-range 50,100 \
  --gap-denom-range 0,5 \
  --gap-num-range 0,1
```

One-note rhythmic output. Useful for drum adornments.
```
melody-maker create \
  --note-length=80,130 \
  --sequence-length-range=80,80 \
  --midi-note-range=50,50 \
  --velocity-range 35,75 \
  --gap-denom-range 5,8 \
  --gap-num-range 0,1
```

### Rest configuration
The spaces between notes played, rests, can be configured through flags:

| Flag                | Example               | Description                                                                                                                                                                                                                  |
|---------------------|-----------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `--gap-bars-range`  | --gap-bars-range 0,4  | Sets the min & max distance (in bars) between each note in the sequence. This example could result in rests 4 bars long.                                                                                                     |
| `--gap-denom-range` | --gap-denom-range 0,8 | Sets the min & max fraction of a note to be rested for.  In this example, the rests could be between a 0th and an 8th of a note.  N.B - Setting this also requires the `--gap-num-range` flag to be set to come into effect. |
| `--gap-num-range`   | --gap-num-range 1,5   | Sets the min & max number of note denoms to use as a rest. In this example, each note in the sequence will have an at least 1 denom rest after it.                                                                           |


### Instruments
This flag sets the channel value. In Logic Pro, these values dictate what instrument is played with loading in the midi file. The instruments available are:
- piano
- synth
- bass
- pluck-synth
- strings
- session-strings
- brass
- trumpet
- edrums
- drums
- organ
- e-piano
- synth-strings
- analog-synth
- synth-brass
- sculpture-synth

## To do / features
- [X] Build initial random midi output
- [X] Expand available scales (minor).
- [ ] Add `save` command to move specified files to a different folder, avoiding overwriting.

Flags:
- [X] Output folder
- [X] Octave range
- [X] Two octave limit
- [X] Note length
- [X] Velocity range
- [X] Range of gap between notes
- [X] Ascending & descending patterns
- [X] Note interval range

Keys, scales & modes:
- [X] Add key selection (so users can generate a melody to go with an existing chord sequence)
- [X] Add scale selection
- [ ] Add modes.
- [ ] `--random` flag, to allow non-scales etc

Error reporting:
- [X] When entering an invalid key
- [X] When entering an invalid scale

Tests:
- [X] Basic bash tests to check build compilation.
- [ ] Go testing for the library.
