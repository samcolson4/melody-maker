# melody-maker
Generate fragments of midi data to help inspire melodies.

Built using the excellent [gomidi](https://github.com/gomidi/midi) and [go-audio/midi](https://github.com/go-audio/midi) packages.

## Usage
When more features have been added, a binary will be provided as part of a release. For now, building locally (`make build`) then examining the help options is the best approach.

`make buildRunDelete` is a helper used during development to create a binary with the latest local changes, generate output and then delete the binary without the manual steps.

### Example
```
$ melody-maker create \
  --note-length=750,1500 \
  --sequence-length-range=10,30 \
  --midi-note-range=50,66 \
  --key G \
  --scale minor
```

```
Writing to folder: 'midi'.
a.mid: [D2 F3 A2 F3 A#2 D#3 A2 D2 A2 C3 F3 G2 C3 F2 G2 F2 D#3 A#2 A2 D#2 C3 F3 G2 F3 F3 A#2 D2 C3]
b.mid: [G2 G2 C3 F3 D#2 D#2 G2 F2 F2 A#2 D#3 G2]
c.mid: [D#2 D3 D3 G2 D2 D3 C3 D3 D3 G2 C3 D2 F2 F3 G2 D#2 C3 G2 C3 D#2]
d.mid: [D#2 D2 A#2 D3 A#2 D#3 F2 C3 F3 C3 A#2 F2 C3 F2 F3 F2 A#2 A#2 F3 D2 D#2 D#3 C3 C3 F3 G2 A2 F2 A2]
e.mid: [F2 G2 D3 F2 A#2 C3 G2 G2 C3 F3 G2 F2 F2 F2 A2 D3 A2 D#3 D#2 D2 D2 G2 C3 F2]
f.mid: [A2 A2 C3 D2 F2 G2 D3 D#2 F3 A#2 D#3 D2 F2 D2 F2 A2 F3 A#2 A#2 D#3]
```

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

## To do
- [X] Build initial random midi output
- [X] Expand available scales (minor).
- [ ] Add `save` command to move specified files to a different folder, avoiding overwriting.

Flags:
- [X] Output folder
- [X] Octave range
- [X] Two octave limit
- [X] Note length
- [X] Velocity range
- [ ] Ascending / descending patterns

Keys & scales:
- [X] Add key selection (so users can generate a melody to go with an existing chord sequence)
- [X] Add scale selection
- [ ] `--random` flag, to allow non-scales etc

Error reporting:
- [X] When entering an invalid key
- [ ] When entering an invalid scale

Tests:
- [ ] Add some...
