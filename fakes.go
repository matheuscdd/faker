package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	brcnpj "github.com/brazilian-utils/go/cnpj"
	brcpf "github.com/brazilian-utils/go/cpf"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/fma0x/faker/lib/ptbr"
	"golang.design/x/clipboard"
)

// runCommand executa um comando externo
func runCommand(cmd string, args []string) error {
	c := exec.Command(cmd, args...)
	return c.Run()
}

var supportedTypes = []string{
	"cpf",
	"cnpj",
	"zip-code",
	"date",
	"full-name",
	"first-name",
	"number",
	"title",
	"company",
 	"email",
}

type DataGenerator struct {
	random   *rand.Rand
	surnames []string
}

func NewDataGenerator() *DataGenerator {
	return &DataGenerator{
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
		surnames: []string{
			"Almeida",
			"Alves",
			"Andrade",
			"Araujo",
			"Barbosa",
			"Barros",
			"Batista",
			"Campos",
			"Cardoso",
			"Carneiro",
			"Carvalho",
			"Castro",
			"Costa",
			"Claus",
			"Dias",
			"Duarte",
			"Freitas",
			"Ferreira",
			"Fernandes",
			"Garcia",
			"Gomes",
			"Goncalves",
			"Lima",
			"Lopes",
			"Machado",
			"Marques",
			"Martins",
			"Melo",
			"Mendes",
			"Monteiro",
			"Moraes",
			"Moreira",
			"Nascimento",
			"Neco",
			"Nogueira",
			"Novaes",
			"Oliveira",
			"Pereira",
			"Pinto",
			"Ramos",
			"Rezende",
			"Ribeiro",
			"Rocha",
			"Rodrigues",
			"Santana",
			"Santos",
			"Silva",
			"Soares",
			"Sousa",
			"Souza",
			"Teixeira",
			"Vieira",
		},
	}
}

func (generator *DataGenerator) CPF() string {
	return brcpf.Generate()
}

func (generator *DataGenerator) CNPJ() string {
	branch := generator.random.Intn(9999) + 1
	return brcnpj.Generate(branch)
}

func (generator *DataGenerator) Phone() string {
	num := generator.random.Intn(90000000000) + 10000000000
	return fmt.Sprintf("%12d", num)
}

func (generator *DataGenerator) ZipCode() string {
	zipCodes := []string{
		"01001-000",
		"01310-000",
		"01414-000",
		"02011-000",
		"02210-000",
		"03045-000",
		"04094-000",
		"05001-000",
		"05407-000",
		"06010-000",
		"07010-000",
		"08010-000",
		"09010-000",
		"11010-000",
		"12010-000",
		"13010-000",
		"14010-000",
		"15010-000",
		"16010-000",
		"17010-000",
	}
	return zipCodes[generator.random.Intn(len(zipCodes))]
}

func (generator *DataGenerator) Date() string {
	minimumDate := time.Now().AddDate(-80, 0, 0)
	maximumDate := time.Now().AddDate(-18, 0, 0)
	secondsRange := maximumDate.Unix() - minimumDate.Unix()
	if secondsRange <= 0 {
		return maximumDate.Format("02/01/2006")
	}

	randomOffset := generator.random.Int63n(secondsRange + 1)
	generatedDate := minimumDate.Add(time.Duration(randomOffset) * time.Second)
	return generatedDate.Format("02/01/2006")
}

func (generator *DataGenerator) FullName() string {
	firstName := generator.FirstName()
	surnameCount := 2
	if generator.random.Intn(4) == 0 {
		surnameCount = 1
	}

	usedSurnames := map[string]struct{}{}
	nameParts := []string{firstName}

	for len(nameParts)-1 < surnameCount {
		surname := generator.randomSurname()
		if _, exists := usedSurnames[surname]; exists {
			continue
		}

		usedSurnames[surname] = struct{}{}
		nameParts = append(nameParts, surname)
	}

	return strings.Join(nameParts, " ")
}

func (generator *DataGenerator) FirstName() string {
	return ptbr.RandomFirstName()
}

func (generator *DataGenerator) Number() string {
	return strconv.Itoa(generator.random.Intn(1000) + 1)
}

func (generator *DataGenerator) Title() string {
	return capitalize(gofakeit.Word())
}

func (generator *DataGenerator) Company() string {
	return gofakeit.Company()
}

func (generator *DataGenerator) Email() string {
	// Usa o gofakeit para gerar um email realista
	return gofakeit.Email()
}

func (generator *DataGenerator) Generate(kind string) (string, error) {
	normalizedKind := strings.TrimSpace(strings.ToLower(kind))
	types := map[string]func() string{
		"cpf":        generator.CPF,
		"cnpj":       generator.CNPJ,
		"zip-code":   generator.ZipCode,
		"date":       generator.Date,
		"full-name":  generator.FullName,
		"first-name": generator.FirstName,
		"number":     generator.Number,
		"title":      generator.Title,
		"company":    generator.Company,
 		"email":      generator.Email,
		"phone":      generator.Phone,
	}

	generatorFunction, exists := types[normalizedKind]
	if !exists {
		return "", fmt.Errorf("unknown type: %q", kind)
	}

	return normalizeSpace(generatorFunction()), nil
}

func (generator *DataGenerator) randomSurname() string {
	return generator.surnames[generator.random.Intn(len(generator.surnames))]
}

func capitalize(value string) string {
	if value == "" {
		return value
	}

	runes := []rune(value)
	firstLetter := strings.ToUpper(string(runes[0]))
	firstLetterRunes := []rune(firstLetter)
	if len(firstLetterRunes) > 0 {
		runes[0] = firstLetterRunes[0]
	}

	return string(runes)
}

func normalizeSpace(value string) string {
	return strings.Join(strings.Fields(value), " ")
}

func printUsage() {
	fmt.Fprintln(os.Stderr, "Usage: generator <type>")
	fmt.Fprintf(os.Stderr, "Types: %s\n", strings.Join(supportedTypes, ", "))
}

func main() {
	if err := clipboard.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize clipboard: %v\n", err)
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	generator := NewDataGenerator()
	kind := os.Args[1]
	value, err := generator.Generate(kind)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		printUsage()
		os.Exit(1)
	}

	clipboard.Write(clipboard.FmtText, []byte(value))
	fmt.Printf("Generated %s: %s (copied to clipboard)\n", kind, value)
}