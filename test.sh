# This file tests basic commands (help, create -h).
# A more complete testing suite is to be developed in Go.

# Help output
testHelpOutput=$(./melody-maker -h)
SUB='Create random'
if [[ "$testHelpOutput" == *"$SUB"* ]]; then
  echo "PASS"
fi

# Create help output
testCreateHelpOutput=$(./melody-maker create -h)
SUB='--file-number int'
if [[ "$testCreateHelpOutput" == *"$SUB"* ]]; then
  echo "PASS"
fi
