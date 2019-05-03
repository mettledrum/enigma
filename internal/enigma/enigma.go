package enigma

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type Decoder struct {
	r io.Reader
	e enigma
}

type Encoder struct {
	w io.Writer
	e enigma
}

type enigma struct {
	rings     []ring
	reflector reflector
	plugboard map[string]string
}

type ring struct {
	position int
	rotor    rotor
}

type reflector struct {
	wiring []string
}

type rotor struct {
	wiring    []string
	turnovers []int
}

const abc = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Config struct {
	RotorNames       []string
	RingPositions    []int
	PluboardPairings []string
}

func NewEnigmaIEncoder(r io.Reader, cfg Config) (*Encoder, error) {
	e, err := newEnigmaI(cfg)
	if err != nil {
		return nil, err
	}

	return *Encoder{
		r: r,
		e: e,
	}, nil
}

func newEnigmaI(cfg Config) (*Encoder, error) {
	const numberOfRotors = 3
	if len(cfg.RingPositions) != 3 || len(cfg.RotorNames) != 3 {
		return nil, errors.New("enigma 1 only has 3 rotors")
	}

	// validate ring settings
	if ringPos1 < 0 || ringPos1 > 25 {
		return nil, errors.New("ring setting I must be [0,25]")
	}
	if ringPos2 < 0 || ringPos2 > 25 {
		return nil, errors.New("ring setting II must be [0,25]")
	}
	if ringPos3 < 0 || ringPos3 > 25 {
		return nil, errors.New("ring setting III must be [0,25]")
	}

	// validate plugboard wirings
	isUppercaseAlphaPair := regexp.MustCompile(`^[A-Z][A-Z]$`).MatchString
	for _, p := range plugboardPairs {
		if !isUppercaseAlphaPair(p) || p[0] == p[1] {
			return nil, fmt.Errorf("plugboard pair %s is invalid", p)
		}
		// TODO: also have to make sure that a plug isn't attempted more than once
		// ie Q cannot be plugged into A and also into B
	}

	availableReflectors := map[string]reflector{
		"A": {wiring: strings.Split("EJMZALYXVBWFCRQUONTSPIKHGD", "")},
		"B": {wiring: strings.Split("YRUHQSLDPXNGOKMIEBFZCWVJAT", "")},
		"C": {wiring: strings.Split("FVPJIAOYEDRZXWGCTKUQSBNMHL", "")},
	}

	availableRotors := map[string]rotor{
		"I": rotor{
			wiring:    strings.Split("EKMFLGDQVZNTOWYHXUSPAIBRCJ", ""),
			turnovers: []int{strings.Index(abc, "Q")},
		},
		"II": rotor{
			wiring:    strings.Split("AJDKSIRUXBLHWTMCQGZNPYFVOE", ""),
			turnovers: []int{strings.Index(abc, "E")},
		},
		"III": rotor{
			wiring:    strings.Split("BDFHJLCPRTXVZNYEIWGAKMUSQO", ""),
			turnovers: []int{strings.Index(abc, "V")},
		},
		"IV": rotor{
			wiring:    strings.Split("ESOVPZJAYQUIRHXLNFTGKDCMWB", ""),
			turnovers: []int{strings.Index(abc, "J")},
		},
		"V": rotor{
			wiring:    strings.Split("VZBRGITYUPSDNHLXAWMJQOFECK", ""),
			turnovers: []int{strings.Index(abc, "Z")},
		},
		"VI": rotor{
			wiring:    strings.Split("JPGVOUMFYQBENHZRDKASXLICTW", ""),
			turnovers: []int{strings.Index(abc, "Z"), strings.Index(abc, "M")},
		},
		"VII": {
			wiring:    strings.Split("NZJHGRCXMYSWBOUFAIVLPEKQDT", ""),
			turnovers: []int{strings.Index(abc, "Z"), strings.Index(abc, "M")},
		},
		"VIII": {
			wiring:    strings.Split("FKQHTLXOCBJSPDZRAMEWNIUYGV", ""),
			turnovers: []int{strings.Index(abc, "Z"), strings.Index(abc, "M")},
		},
	}

	// validate rotor names
	rot1, ok := availableRotors[rotorName1]
	if !ok {
		return nil, fmt.Errorf("rotor name: %s not found for position 1", rotorName1)
	}
	rot2, ok := availableRotors[rotorName2]
	if !ok {
		return nil, fmt.Errorf("rotor name: %s not found for position 2", rotorName2)
	}
	rot3, ok := availableRotors[rotorName3]
	if !ok {
		return nil, fmt.Errorf("rotor name: %s not found for position 3", rotorName3)
	}

	// validate reflector name
	ref, ok := availableReflectors[reflectorName]
	if !ok {
		return nil, fmt.Errorf("reflector name: %s not found", reflectorName)
	}

	var pb map[string]string
	for _, pp := range plugboardPairs {
		pb[string(pp[0])] = string(pp[1])
		pb[string(pp[1])] = string(pp[0])
	}

	return &Enigma{
		w: w,
		r: r,
		rings: []ring{
			{rotor: rot1, position: ringPos1},
			{rotor: rot2, position: ringPos2},
			{rotor: rot3, position: ringPos3},
		},
		reflector: ref,
		plugboard: pb,
	}, nil
}

func (e *Encoder) Encode(in string) error {
	if !regexp.MustCompile(`^[A-Z]*$`).MatchString(in) {
		return errors.New("cannot encode non-alpha text")
	}

	// plugboard -> rings -> reflector -> reverse rings -> plugboard
	// for _, t := range text {

	// 	// only run plugboard if it exists; some models don't have one
	// 	// plugIn := e.plugboard[string(t)]

	// 	for _, ring := range e.rings {

	// 	}

	// 	// plugOut := e.plugboard[t]
	// }
	return 0, nil
}
