package enigma

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type Decoder struct {
	eng enigma
	r   io.Reader
}

type Encoder struct {
	eng *enigma
	w   io.Writer
}

type enigma struct {
	plugboard map[string]string
	reflector rotor
	rotors    []rotor
}

// ABC is used for alphabetic indexing
const ABC = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// ShouldTurnover indicates if the rotor to the left should rotate
func (r *rotor) ShouldTurnover() bool {
	return r.turnovers[r.position]
}

// Rotate advances the rotor one position
func (r *rotor) Rotate() {
	r.position = (r.position + 1) % 26
}

func (r *rotor) GetEncodedLetterIn(idx int) int {
	// add position to that index
	idxWithPosition := (idx + r.position + 26) % 26

	// look it up
	letterThruWire := r.wiring[idxWithPosition]
	idxThruWire := strings.Index(ABC, letterThruWire)

	// remove positioning
	out := (idxThruWire - r.position + 26) % 26
	return out
}

func (r *rotor) GetEncodedLetterOut(idx int) int {
	// add position to that index
	idxWithPosition := (idx + r.position + 26) % 26

	// look it up
	letterThruWire := strings.Split(ABC, "")[idxWithPosition]
	idxThruWire := strings.Index(strings.Join(r.wiring, ""), letterThruWire)

	// remove position
	out := (idxThruWire - r.position + 26) % 26
	return out
}

type rotor struct {
	position  int
	turnovers map[int]bool
	wiring    []string
}

type RotorPosition struct {
	Name     string
	Position int
}

// Config is how any enigma impl. can be setup
type Config struct {
	PluboardWirings []string
	Reflector       string
	RotorPositions  []RotorPosition
}

// EncodeString takes a [A-Z]* string and encodes it using the underlying enigma
// and writes to the writer.
func (e *Encoder) EncodeString(userInput string) error {
	if !regexp.MustCompile(`^[A-Z]*$`).MatchString(userInput) {
		return errors.New("can only encode [A-Z] text")
	}

	for _, letter := range userInput {
		encoded := e.eng.Type(string(letter))
		if _, err := e.w.Write([]byte(encoded)); err != nil {
			return err
		}
	}
	return nil
}

// Type encodes the letter passed and rotates the rotors's positions
// plugboard -> rings -> reflector -> reverse rings -> plugboard
func (en *enigma) Type(userLetter string) string {
	// rotate before encoding letter
	en.rotateRotors()

	fmt.Printf("user in:\t%s\n", userLetter)
	encoded := userLetter

	// plugboard in
	if p, ok := en.plugboard[encoded]; ok {
		encoded = p
	}
	fmt.Printf("plugboard in:\t%s\n", encoded)

	idx := strings.Index(ABC, encoded)

	// rotors from right to left
	for i := len(en.rotors) - 1; i >= 0; i-- {
		idx = en.rotors[i].GetEncodedLetterIn(idx)
		fmt.Printf("rotor[%d] in:\t%d\t%+v\n", i, idx, en.rotors[i].wiring)
	}

	// reflector
	idx = en.reflector.GetEncodedLetterOut(idx)
	fmt.Printf("reflector:\t%d\t%+v\n", idx, en.reflector.wiring)

	// rotors from left to right
	for i := 0; i < len(en.rotors); i++ {
		idx = en.rotors[i].GetEncodedLetterOut(idx)
		fmt.Printf("rotor[%d] out:\t%d\t%+v\n", i, idx, en.rotors[i].wiring)
	}

	// plugboard out
	encoded = strings.Split(ABC, "")[idx]
	if p, ok := en.plugboard[encoded]; ok {
		encoded = p
	}
	fmt.Printf("plugboard out:\t%s\n\n", encoded)

	return encoded
}

// rotateRotors determines which rotors to rotate
func (en *enigma) rotateRotors() {
	for i := 0; i < len(en.rotors); i++ {
		switch i {
		case 0: // leftmost rotor only looks at neighbor to the right
			if en.rotors[i+1].ShouldTurnover() {
				en.rotors[i].Rotate()
			}
		case len(en.rotors) - 1: // rightmost rotor always turns
			en.rotors[i].Rotate()
		default: // middle rotor(s) turn if self or right neighbor are on a notch
			if en.rotors[i+1].ShouldTurnover() || en.rotors[i].ShouldTurnover() {
				en.rotors[i].Rotate()
			}
		}
	}

	// show rotor settings as letters
	ps := make([]string, len(en.rotors))
	for i, r := range en.rotors {
		ps[i] = string(ABC[r.position])
	}
	fmt.Printf("%v\n", ps)
}
