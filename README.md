# melody-maker
Generate fragments of midi data to help inspire melodies.

Built using the excellent [gomidi](https://github.com/gomidi/midi) package.

# Usage
Clone the repo to your local system, ensure you have the dependencies installed, then:
```
make buildAndRun
```

When more features have been added, a binary will be provided as part of a release.

### To do
- [X] Build initial random midi output

Flags:
- [ ] Velocity range
- [ ] Octave range

Keys & scales:
- [ ] Add key selection (so users can generate a melody to go with an existing chord sequence)
- [ ] Add scale selection
- [ ] `--random` flag, to allow non-scales etc

Tests:
- [ ] Add some...
