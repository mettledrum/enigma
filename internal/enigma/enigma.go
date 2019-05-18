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

// ABC is used for alphabetic indexing.
const ABC = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// shouldTurnover indicates if the rotor to the left should rotate.
func (r *rotor) shouldTurnover() bool {
	ltr := string(ABC[r.grundStellung])
	for _, n := range r.notches {
		if ltr == n {
			return true
		}
	}
	return false
}

// rotate advances the rotor one position.
func (r *rotor) rotate() {
	r.grundStellung = (r.grundStellung + 1) % len(r.wiring)
}

// unShift undoes the ring settings and rotor rotations for passing the index
// along to the next rotor or reflector.
func (r *rotor) unShift(idx int) int {
	idxWithoutRing := (idx + r.ringStellung + len(r.wiring)) % len(r.wiring)
	idxWithoutPosition := (idxWithoutRing - r.grundStellung + len(r.wiring)) % len(r.wiring)
	return idxWithoutPosition
}

// shift accounts for the ring position and rotor rotations.
func (r *rotor) shift(idx int) int {
	idxWithPosition := (idx + r.grundStellung + len(r.wiring)) % len(r.wiring)
	idxWithRing := (idxWithPosition - r.ringStellung + len(r.wiring)) % len(r.wiring)
	return idxWithRing
}

// getEncodedIdxIn runs the letter through the rotor L<-R.
func (r *rotor) getEncodedIdxIn(idx int) int {
	idx = r.shift(idx)

	letterThruWire := r.wiring[idx]
	idxThruWire := strings.Index(ABC, string(letterThruWire))

	return r.unShift(idxThruWire)
}

// getEncodedIdxOut runs the letter through the rotor L->R.
func (r *rotor) getEncodedIdxOut(idx int) int {
	idx = r.shift(idx)

	letterThruWire := strings.Split(ABC, "")[idx]
	idxThruWire := strings.Index(r.wiring, letterThruWire)

	return r.unShift(idxThruWire)
}

type rotor struct {
	grundStellung int
	ringStellung  int
	notches       []string
	wiring        string
}

// RotorPosition is used to setup an enigma's rotors
type RotorPosition struct {
	Walzenlage    string // wheel name
	GrundStellung int    // initial setting
	RingStellung  int    // setting that rotates
}

// Config is how any enigma impl. can be setup
type Config struct {
	PluboardWirings []string        // pairings of letters for the plugboard
	Reflector       string          // name of reflector
	RotorPositions  []RotorPosition // rotors and rotation settings
}

// EncodeString takes a [A-Z]* string and encodes it using the underlying enigma
// and writes to the writer.
func (e *Encoder) EncodeString(userInput string) error {
	if !regexp.MustCompile(`^[A-Z]*$`).MatchString(userInput) {
		return errors.New("can only encode [A-Z] text")
	}

	for _, letter := range userInput {
		encoded := e.eng.typeLetter(string(letter))
		if _, err := e.w.Write([]byte(encoded)); err != nil {
			return err
		}
	}
	return nil
}

// typeLetter encodes the letter passed and rotates the rotors's positions.
// plugboard -> rings -> reflector -> reverse rings -> plugboard
func (en *enigma) typeLetter(userLetter string) string {
	out := userLetter

	// rotate before encoding letter
	en.rotateRotors()

	// plugboard in
	if p, ok := en.plugboard[out]; ok {
		out = p
	}
	idx := strings.Index(ABC, out)

	// rotors from right to left
	for i := len(en.rotors) - 1; i >= 0; i-- {
		idx = en.rotors[i].getEncodedIdxIn(idx)
	}

	// reflector
	idx = en.reflector.getEncodedIdxOut(idx)

	// rotors from left to right
	for i := 0; i < len(en.rotors); i++ {
		idx = en.rotors[i].getEncodedIdxOut(idx)
	}

	// plugboard out
	out = strings.Split(ABC, "")[idx]
	if p, ok := en.plugboard[out]; ok {
		out = p
	}

	return out
}

// rotateRotors determines which rotors to rotate
func (en *enigma) rotateRotors() {
	for i := 0; i < len(en.rotors); i++ {
		switch i {
		case 0: // leftmost rotor only looks at neighbor to the right
			if en.rotors[i+1].shouldTurnover() {
				en.rotors[i].rotate()
			}
		case len(en.rotors) - 1: // rightmost rotor always turns
			en.rotors[i].rotate()
		case len(en.rotors) - 2: // middle rotor turns if self or right neighbor are on a notch
			if en.rotors[i+1].shouldTurnover() || en.rotors[i].shouldTurnover() {
				en.rotors[i].rotate()
			}
		default: // 4th rotor (on M4) doesn't rotate
		}
	}
}
