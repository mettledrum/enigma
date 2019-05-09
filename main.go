package main

import (
	"bytes"
	"enigma/internal/enigma"
	"fmt"
	"io/ioutil"
)

func main() {
	c := enigma.Config{
		Reflector:       "UKW-B",
		PluboardWirings: []string{},
		RotorPositions: []enigma.RotorPosition{
			{Name: "I", Position: 0},
			{Name: "II", Position: 0},
			{Name: "III", Position: 0},
			// {Name: "II", Position: strings.Index(enigma.ABC, "D")},
			// {Name: "III", Position: strings.Index(enigma.ABC, "T")}, // TODO make helper for idx instead of ABC const
		},
	}
	w := &bytes.Buffer{}
	e, err := enigma.NewEnigmaM3Encoder(w, c)
	if err != nil {
		panic(err)
	}
	err = e.EncodeString("AAAAA")
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(w)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
