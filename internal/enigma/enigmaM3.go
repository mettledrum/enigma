package enigma

import (
	"fmt"
	"io"
	"regexp"

	"github.com/pkg/errors"
)

// NewEnigmaM3Encoder returns a new encoder.
func NewEnigmaM3Encoder(w io.Writer, cfg Config) (*Encoder, error) {
	e, err := newEnigmaM3(cfg)
	if err != nil {
		return nil, err
	}

	return &Encoder{
		w:   w,
		eng: e,
	}, nil
}

func newEnigmaM3(cfg Config) (*enigma, error) {
	if err := validateReflectorForEnigmaM3(cfg.Reflector); err != nil {
		return nil, err
	}
	if err := validateRotorsForEnigmaM3(cfg.RotorPositions); err != nil {
		return nil, err
	}
	if err := validatePlugboard(cfg.PluboardWirings); err != nil {
		return nil, err
	}

	// create plugboard 2-way mapping
	// ie []string{"AR"} -> map[string]string{"A":"R","R":"A"}
	pb := map[string]string{}
	for _, pp := range cfg.PluboardWirings {
		pb[string(pp[0])] = string(pp[1])
		pb[string(pp[1])] = string(pp[0])
	}

	// setup rings in order with positions
	// index 0 is the leftmost ring
	rts := []rotor{}
	for _, r := range cfg.RotorPositions {
		rt := rotorsForEnigmaM3()[r.Walzenlage]
		rt.grundStellung = r.GrundStellung
		rt.ringStellung = r.RingStellung
		rts = append(rts, rt)
	}

	return &enigma{
		plugboard: pb,
		reflector: reflectorsForEnigmaM3()[cfg.Reflector],
		rotors:    rts,
	}, nil
}

func validateReflectorForEnigmaM3(ref string) error {
	if _, ok := reflectorsForEnigmaM3()[ref]; !ok {
		return errors.New("invalid reflector name for Enigma M3")
	}
	return nil
}

func reflectorsForEnigmaM3() map[string]rotor {
	return map[string]rotor{
		"UKW-B": rotor{wiring: "YRUHQSLDPXNGOKMIEBFZCWVJAT"},
		"UKW-C": rotor{wiring: "FVPJIAOYEDRZXWGCTKUQSBNMHL"},
	}
}

func checkReuse(k string, m map[string]bool) error {
	if m[k] {
		return fmt.Errorf("may not reuse %s", k)
	}
	m[k] = true
	return nil
}

func validatePlugboard(pbs []string) error {
	lettersUsed := map[string]bool{}
	isUppercaseAlphaPair := regexp.MustCompile(`^[A-Z][A-Z]$`).MatchString

	for _, p := range pbs {
		if !isUppercaseAlphaPair(p) {
			return fmt.Errorf("plugboard pair %s must be two capitalized, differing letters", p)
		}
		firstLetter := string(p[0])
		secLetter := string(p[1])

		if err := checkReuse(firstLetter, lettersUsed); err != nil {
			return errors.Wrap(err, "plugboard letters must be unique")
		}
		if err := checkReuse(secLetter, lettersUsed); err != nil {
			return errors.Wrap(err, "plugboard letters must be unique")
		}
	}

	return nil
}

func validateRotorsForEnigmaM3(rs []RotorPosition) error {
	const numberOfRotors = 3

	if len(rs) != numberOfRotors {
		return fmt.Errorf("Enigma M3 has %d rotors", numberOfRotors)
	}

	rotorsUsed := map[string]bool{}
	for _, r := range rs {
		if _, ok := rotorsForEnigmaM3()[r.Walzenlage]; !ok {
			return fmt.Errorf("Enigma M3 rotor %s not allowed", r.Walzenlage)
		}
		if r.GrundStellung < 0 || r.GrundStellung > 25 {
			return errors.New("rotor grundstellung must be in range [0,25]")
		}
		if r.RingStellung < 0 || r.RingStellung > 25 {
			return errors.New("rotor ringstellung must be in range [0,25]")
		}

		if err := checkReuse(r.Walzenlage, rotorsUsed); err != nil {
			return errors.Wrap(err, "rotors must be different")
		}
	}

	return nil
}

func rotorsForEnigmaM3() map[string]rotor {
	return map[string]rotor{
		"I": rotor{
			wiring:  "EKMFLGDQVZNTOWYHXUSPAIBRCJ",
			notches: []string{"Q"},
		},
		"II": rotor{
			wiring:  "AJDKSIRUXBLHWTMCQGZNPYFVOE",
			notches: []string{"E"},
		},
		"III": rotor{
			wiring:  "BDFHJLCPRTXVZNYEIWGAKMUSQO",
			notches: []string{"V"},
		},
		"IV": rotor{
			wiring:  "ESOVPZJAYQUIRHXLNFTGKDCMWB",
			notches: []string{"J"},
		},
		"V": rotor{
			wiring:  "VZBRGITYUPSDNHLXAWMJQOFECK",
			notches: []string{"Z"},
		},
		"VI": rotor{
			wiring:  "JPGVOUMFYQBENHZRDKASXLICTW",
			notches: []string{"Z", "M"},
		},
		"VII": {
			wiring:  "NZJHGRCXMYSWBOUFAIVLPEKQDT",
			notches: []string{"Z", "M"},
		},
		"VIII": {
			wiring:  "FKQHTLXOCBJSPDZRAMEWNIUYGV",
			notches: []string{"Z", "M"},
		},
	}
}
