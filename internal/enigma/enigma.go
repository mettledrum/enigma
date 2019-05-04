package enigma

import (
	"errors"
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

const abc = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// ShouldTurnover indicates if the rotor to the left should rotate
func (r *rotor) ShouldTurnover() bool {
	return r.turnovers[r.position]
}

// Rotate advances the rotor one position
func (r *rotor) Rotate() {
	r.position = (r.position + 1) % len(r.wiring)
}

// GetEncodedLetter gets the new letter passed through the rotor's wiring
func (r *rotor) GetEncodedLetter(in string) string {
	i := strings.Index(abc, in)
	return r.wiring[i]
}

type rotor struct {
	position  int
	turnovers map[int]bool
	wiring    []string
}

type Config struct {
	PluboardWirings []string
	Reflector       string
	RotorPositions  map[string]int
}

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

// TODO advance based on turnovers
// plugboard -> rings -> reflector -> reverse rings -> plugboard
func (en *enigma) Type(userInput string) string {
	encoded := userInput

	// plugboard
	if p, ok := en.plugboard[encoded]; ok {
		encoded = p
	}

	// rotors from right to left
	for i := len(en.rotors) - 1; i >= 0; i-- {
		encoded = en.rotors[i].GetEncodedLetter(encoded)
	}

	// reflector
	encoded = en.reflector.GetEncodedLetter(encoded)

	// rotors from left to right
	for i := 0; i < len(en.rotors); i++ {
		encoded = en.rotors[i].GetEncodedLetter(encoded)
	}

	// final plugboard
	if p, ok := en.plugboard[encoded]; ok {
		encoded = p
	}

	en.rotateRotors()

	return encoded
}

// rotateRotors looks at the rotor turnover settings
func (en *enigma) rotateRotors() {
	// look at leftmost rotors to see if their neighbors to the right rotate them
	for i := 0; i < len(en.rotors)-1; i++ {
		if en.rotors[i+1].ShouldTurnover() {
			en.rotors[i].Rotate()
		}
	}
	// always rotate the rightmost rotor
	en.rotors[len(en.rotors)-1].Rotate()
}
