package enigma

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

func NewEnigmaIEncoder(w io.Writer, cfg Config) (*Encoder, error) {
	e, err := newEnigmaI(cfg)
	if err != nil {
		return nil, err
	}

	return &Encoder{
		w:   w,
		eng: e,
	}, nil
}

func newEnigmaI(cfg Config) (*enigma, error) {
	if err := validateReflectorForEnigmaI(cfg.Reflector); err != nil {
		return nil, err
	}
	if err := validateRotorsForEnigmaI(cfg.RotorPositions); err != nil {
		return nil, err
	}
	if err := validatePlugboardForEnigmaI(cfg.PluboardWirings); err != nil {
		return nil, err
	}

	// create plugboard 2-way mapping
	// ie []string{"AR"} -> map[string]string{"A":"R","R":"A"}
	var pb map[string]string
	for _, pp := range cfg.PluboardWirings {
		pb[string(pp[0])] = string(pp[1])
		pb[string(pp[1])] = string(pp[0])
	}

	// setup rings in order with positions
	// index 0 is the leftmost ring
	rts := []rotor{}
	for name, pos := range cfg.RotorPositions {
		rt := rotorsForEnigmaI()[name]
		rt.position = pos
		rts = append(rts, rt)
	}

	return &enigma{
		plugboard: pb,
		reflector: reflectorsForEnigmaI()[cfg.Reflector],
		rotors:    rts,
	}, nil
}

func validateReflectorForEnigmaI(ref string) error {
	if _, ok := reflectorsForEnigmaI()[ref]; !ok {
		return errors.New("invalid reflector name for Enigma I")
	}
	return nil
}

func reflectorsForEnigmaI() map[string]rotor {
	return map[string]rotor{
		"A": rotor{wiring: strings.Split("EJMZALYXVBWFCRQUONTSPIKHGD", "")},
		"B": rotor{wiring: strings.Split("YRUHQSLDPXNGOKMIEBFZCWVJAT", "")},
		"C": rotor{wiring: strings.Split("FVPJIAOYEDRZXWGCTKUQSBNMHL", "")},
	}
}

func validatePlugboardForEnigmaI(pbs []string) error {
	// validate plugboard wirings
	isUppercaseAlphaPair := regexp.MustCompile(`^[A-Z][A-Z]$`).MatchString
	for _, p := range pbs {
		if !isUppercaseAlphaPair(p) || p[0] == p[1] {
			return fmt.Errorf("plugboard pair %s is invalid", p)
		}
		// TODO: also have to make sure that a plug isn't attempted more than once
		// ie Q cannot be plugged into A and also into B
	}
	return nil
}

// TODO cannot reuse a rotor
func validateRotorsForEnigmaI(rs map[string]int) error {
	const numberOfRotors = 3

	if len(rs) != numberOfRotors {
		return fmt.Errorf("Enigma I has %d rotors", numberOfRotors)
	}

	for name, pos := range rs {
		if _, ok := rotorsForEnigmaI()[name]; !ok {
			return fmt.Errorf("Enigma I rotor %s not allowed", name)
		}
		if pos < 0 || pos > 25 {
			return errors.New("rotor position must be in range [0,25]")
		}
	}

	return nil
}

func rotorsForEnigmaI() map[string]rotor {
	return map[string]rotor{
		"I": rotor{
			wiring:    strings.Split("EKMFLGDQVZNTOWYHXUSPAIBRCJ", ""),
			turnovers: map[int]bool{strings.Index(abc, "Q"): true},
		},
		"II": rotor{
			wiring:    strings.Split("AJDKSIRUXBLHWTMCQGZNPYFVOE", ""),
			turnovers: map[int]bool{strings.Index(abc, "E"): true},
		},
		"III": rotor{
			wiring:    strings.Split("BDFHJLCPRTXVZNYEIWGAKMUSQO", ""),
			turnovers: map[int]bool{strings.Index(abc, "V"): true},
		},
		"IV": rotor{
			wiring:    strings.Split("ESOVPZJAYQUIRHXLNFTGKDCMWB", ""),
			turnovers: map[int]bool{strings.Index(abc, "J"): true},
		},
		"V": rotor{
			wiring:    strings.Split("VZBRGITYUPSDNHLXAWMJQOFECK", ""),
			turnovers: map[int]bool{strings.Index(abc, "Z"): true},
		},
		"VI": rotor{
			wiring:    strings.Split("JPGVOUMFYQBENHZRDKASXLICTW", ""),
			turnovers: map[int]bool{strings.Index(abc, "Z"): true, strings.Index(abc, "M"): true},
		},
		"VII": {
			wiring:    strings.Split("NZJHGRCXMYSWBOUFAIVLPEKQDT", ""),
			turnovers: map[int]bool{strings.Index(abc, "Z"): true, strings.Index(abc, "M"): true},
		},
		"VIII": {
			wiring:    strings.Split("FKQHTLXOCBJSPDZRAMEWNIUYGV", ""),
			turnovers: map[int]bool{strings.Index(abc, "Z"): true, strings.Index(abc, "M"): true},
		},
	}
}
