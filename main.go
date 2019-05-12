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
		PluboardWirings: []string{"AN"},
		RotorPositions: []enigma.RotorPosition{
			{Walzenlage: "III", GrundStellung: 0, RingStellung: 3},
			{Walzenlage: "VI", GrundStellung: 0, RingStellung: 0},
			{Walzenlage: "VII", GrundStellung: 0, RingStellung: 0},
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
