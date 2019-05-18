package enigma_test

import (
	"bytes"
	"enigma/internal/enigma"
	"io/ioutil"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEnigma(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Enigma Suite")
}

var _ = Describe("enigma", func() {

	Context("M3", func() {

		var (
			validM3Cfg enigma.Config
			w          *bytes.Buffer
		)

		BeforeEach(func() {
			w = &bytes.Buffer{}
			validM3Cfg = enigma.Config{
				Reflector:       "UKW-B",
				PluboardWirings: []string{"AN", "QT"},
				RotorPositions: []enigma.RotorPosition{
					{Walzenlage: "IV", GrundStellung: 8, RingStellung: 3},
					{Walzenlage: "V", GrundStellung: 7, RingStellung: 7},
					{Walzenlage: "VI", GrundStellung: 20, RingStellung: 0},
				},
			}
		})

		Context("#NewEnigmaM3Encoder", func() {

			Context("valid config", func() {

				It("creates encoder", func() {
					_, err := enigma.NewEnigmaM3Encoder(w, validM3Cfg)
					Expect(err).To(BeNil())
				})
			})

			Context("invalid config", func() {

				It("can't repeat plugboard letters in different pairs", func() {
					validM3Cfg.PluboardWirings = []string{"AN", "NM"}
					_, err := enigma.NewEnigmaM3Encoder(w, validM3Cfg)
					Expect(err).NotTo(BeNil())
				})

				It("can't repeat plugboard letters in pair", func() {
					validM3Cfg.PluboardWirings = []string{"AA"}
					_, err := enigma.NewEnigmaM3Encoder(w, validM3Cfg)
					Expect(err).NotTo(BeNil())
				})

				It("plugboard letters cannot be empty", func() {
					validM3Cfg.PluboardWirings = []string{"AQ", ""}
					_, err := enigma.NewEnigmaM3Encoder(w, validM3Cfg)
					Expect(err).NotTo(BeNil())
				})

				It("plugboard letters must be pairs", func() {
					validM3Cfg.PluboardWirings = []string{"M"}
					_, err := enigma.NewEnigmaM3Encoder(w, validM3Cfg)
					Expect(err).NotTo(BeNil())
				})

				It("rotors' names must be valid", func() {
					validM3Cfg.RotorPositions = []enigma.RotorPosition{
						{Walzenlage: "some weird rotor name"},
						{Walzenlage: "I"},
						{Walzenlage: "II"},
					}
					_, err := enigma.NewEnigmaM3Encoder(w, validM3Cfg)
					Expect(err).NotTo(BeNil())
				})

				It("rotors cannot repeat", func() {
					validM3Cfg.RotorPositions = []enigma.RotorPosition{
						{Walzenlage: "II"},
						{Walzenlage: "I"},
						{Walzenlage: "II"},
					}
					_, err := enigma.NewEnigmaM3Encoder(w, validM3Cfg)
					Expect(err).NotTo(BeNil())
				})

				It("rotors cannot be blank", func() {
					validM3Cfg.RotorPositions = []enigma.RotorPosition{
						{Walzenlage: ""},
						{Walzenlage: "III"},
						{Walzenlage: "II"},
					}
					_, err := enigma.NewEnigmaM3Encoder(w, validM3Cfg)
					Expect(err).NotTo(BeNil())
				})

				It("must have 3 rotors", func() {
					validM3Cfg.RotorPositions = []enigma.RotorPosition{
						{Walzenlage: "III"},
						{Walzenlage: "II"},
					}
					_, err := enigma.NewEnigmaM3Encoder(w, validM3Cfg)
					Expect(err).NotTo(BeNil())
				})

				It("reflector cannot be blank", func() {
					validM3Cfg.Reflector = ""
					_, err := enigma.NewEnigmaM3Encoder(w, validM3Cfg)
					Expect(err).NotTo(BeNil())
				})

				It("reflector name must be valid", func() {
					validM3Cfg.Reflector = "some reflektorrr"
					_, err := enigma.NewEnigmaM3Encoder(w, validM3Cfg)
					Expect(err).NotTo(BeNil())
				})

				It("ringstellung must be [0", func() {
					validM3Cfg.RotorPositions[0].RingStellung = -5
					_, err := enigma.NewEnigmaM3Encoder(w, validM3Cfg)
					Expect(err).NotTo(BeNil())
				})

				It("ringstellung must be 25]", func() {
					validM3Cfg.RotorPositions[0].RingStellung = 26
					_, err := enigma.NewEnigmaM3Encoder(w, validM3Cfg)
					Expect(err).NotTo(BeNil())
				})

				It("grundstellung must be [0", func() {
					validM3Cfg.RotorPositions[1].GrundStellung = -50
					_, err := enigma.NewEnigmaM3Encoder(w, validM3Cfg)
					Expect(err).NotTo(BeNil())
				})

				It("grundstellung must be 25]", func() {
					validM3Cfg.RotorPositions[2].GrundStellung = 99
					_, err := enigma.NewEnigmaM3Encoder(w, validM3Cfg)
					Expect(err).NotTo(BeNil())
				})
			})
		})

		Context("#Encode", func() {

			It("encodes messages [1]", func() {
				c := enigma.Config{
					Reflector:       "UKW-C",
					PluboardWirings: []string{"AC", "XY"},
					RotorPositions: []enigma.RotorPosition{
						{Walzenlage: "I", GrundStellung: 0, RingStellung: 0},
						{Walzenlage: "II", GrundStellung: 0, RingStellung: 0},
						{Walzenlage: "III", GrundStellung: 0, RingStellung: 0},
					},
				}
				e, err := enigma.NewEnigmaM3Encoder(w, c)
				Expect(err).To(BeNil())

				Expect(e.EncodeString("ABCDEFGH")).To(BeNil())
				b, err := ioutil.ReadAll(w)
				Expect(err).To(BeNil())
				Expect(string(b)).To(Equal("BYBVVLXJ"))
			})

			It("encodes messages [2]", func() {
				c := enigma.Config{
					Reflector:       "UKW-B",
					PluboardWirings: []string{"AC", "XY", "BT"},
					RotorPositions: []enigma.RotorPosition{
						{Walzenlage: "II", GrundStellung: 4, RingStellung: 2},
						{Walzenlage: "IV", GrundStellung: 0, RingStellung: 25},
						{Walzenlage: "VI", GrundStellung: 11, RingStellung: 1},
					},
				}
				e, err := enigma.NewEnigmaM3Encoder(w, c)
				Expect(err).To(BeNil())

				Expect(e.EncodeString("QWERTYQWERTYQWERTYQWERTYQWERTYQWERTYQWERTYQWERTY")).To(BeNil())
				b, err := ioutil.ReadAll(w)
				Expect(err).To(BeNil())
				Expect(string(b)).To(Equal("VYNPLKEQMZHAGLIUSPXZDCZMDOAJOULIIUXRKQKNXBVICIZK"))
			})
		})
	})
})
