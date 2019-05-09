package enigma

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

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
	if err := validatePlugboardForEnigmaM3(cfg.PluboardWirings); err != nil {
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
		rt := rotorsForEnigmaM3()[r.Name]
		rt.position = r.Position
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
		"UKW-B": rotor{wiring: strings.Split("YRUHQSLDPXNGOKMIEBFZCWVJAT", "")},
		"UKW-C": rotor{wiring: strings.Split("FVPJIAOYEDRZXWGCTKUQSBNMHL", "")},
	}
}

func validatePlugboardForEnigmaM3(pbs []string) error {
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

// TODO: cannot reuse a rotor
func validateRotorsForEnigmaM3(rs []RotorPosition) error {
	const numberOfRotors = 3

	if len(rs) != numberOfRotors {
		return fmt.Errorf("Enigma M3 has %d rotors", numberOfRotors)
	}

	for _, r := range rs {
		if _, ok := rotorsForEnigmaM3()[r.Name]; !ok {
			return fmt.Errorf("Enigma M3 rotor %s not allowed", r.Name)
		}
		if r.Position < 0 || r.Position > 25 {
			return errors.New("rotor position must be in range [0,25]")
		}
	}

	return nil
}

func rotorsForEnigmaM3() map[string]rotor {
	return map[string]rotor{
		"I": rotor{
			wiring:    strings.Split("EKMFLGDQVZNTOWYHXUSPAIBRCJ", ""),
			turnovers: map[int]bool{strings.Index(ABC, "Q"): true},
		},
		"II": rotor{
			wiring:    strings.Split("AJDKSIRUXBLHWTMCQGZNPYFVOE", ""),
			turnovers: map[int]bool{strings.Index(ABC, "E"): true},
		},
		"III": rotor{
			wiring:    strings.Split("BDFHJLCPRTXVZNYEIWGAKMUSQO", ""),
			turnovers: map[int]bool{strings.Index(ABC, "V"): true},
		},
		"IV": rotor{
			wiring:    strings.Split("ESOVPZJAYQUIRHXLNFTGKDCMWB", ""),
			turnovers: map[int]bool{strings.Index(ABC, "J"): true},
		},
		"V": rotor{
			wiring:    strings.Split("VZBRGITYUPSDNHLXAWMJQOFECK", ""),
			turnovers: map[int]bool{strings.Index(ABC, "Z"): true},
		},
		"VI": rotor{
			wiring:    strings.Split("JPGVOUMFYQBENHZRDKASXLICTW", ""),
			turnovers: map[int]bool{strings.Index(ABC, "Z"): true, strings.Index(ABC, "M"): true},
		},
		"VII": {
			wiring:    strings.Split("NZJHGRCXMYSWBOUFAIVLPEKQDT", ""),
			turnovers: map[int]bool{strings.Index(ABC, "Z"): true, strings.Index(ABC, "M"): true},
		},
		"VIII": {
			wiring:    strings.Split("FKQHTLXOCBJSPDZRAMEWNIUYGV", ""),
			turnovers: map[int]bool{strings.Index(ABC, "Z"): true, strings.Index(ABC, "M"): true},
		},
	}
}