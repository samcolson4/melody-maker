# melody-maker
Generate fragments of midi data to help inspire melodies.

Built using the excellent [gomidi](https://github.com/gomidi/midi) and [go-audio/midi](https://github.com/go-audio/midi) packages.

## Usage
When more features have been added, a binary will be provided as part of a release. For now, building locally (`make build`) then examining the help options is the best approach.

`make buildRunDelete` is a helper used during development to create a binary with the latest local changes, generate output and then delete the binary without the manual steps.

### Example output
```
$ melody-maker create

Written file 'a.mid': [F#5 B4 D#1]
Written file 'b.mid': [E6 B3 G1 F#4 B4 D#2 A#5 F5 D1]
Written file 'c.mid': [G1 F#1 F#3 A#1 E4 A3 G3]
Written file 'd.mid': [C4 G#4 D#4 C3 D4 G#1 G#4 G4 G1 A4]
Written file 'e.mid': [C#6 C#3 B5 B0 C#6 F2 C3 C1 C1 C4]
```

## To do
- [X] Build initial random midi output
- [ ] Expand available scales (minor).
- [ ] Add `save` command to move specified files to a different folder, avoiding overwriting.

Flags:
- [X] Output folder
- [X] Octave range
- [X] Two octave limit
- [ ] Velocity range

Keys & scales:
- [X] Add key selection (so users can generate a melody to go with an existing chord sequence)
- [X] Add scale selection
- [ ] `--random` flag, to allow non-scales etc

Error reporting:
- [X] When entering an invalid key
- [ ] When entering an invalid scale

Tests:
- [ ] Add some...
