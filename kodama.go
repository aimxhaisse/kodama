package main

import (
	"bufio"
	"errors"
	"flag"
	"image"
	"image/jpeg"
	"io"
	"kodama/filters"
	"log"
	"os"
	"runtime"
	"strings"
)

var input_file = flag.String("infile", "", "input file")

// Split a job and send chunks to several workers
func ProcessInParallel(in image.Image, jobs int, worker filters.Filter) image.Image {
	bounds := in.Bounds()
	x_unit := (bounds.Max.X - bounds.Min.X) / jobs
	ch := make(chan bool)
	out := image.NewRGBA(bounds)
	for i := 0; i < jobs; i++ {
		min := image.Point{i * x_unit, bounds.Min.Y}
		max := image.Point{(i + 1) * x_unit, bounds.Max.Y}
		area := image.Rectangle{min, max}
		log.Printf("first job going from [%d,%d] to [%d,%d]", area.Min.X, area.Min.Y, area.Max.X, area.Max.Y)
		go func(ch chan bool) {
			worker.Process(in, out, area)
			ch <- true
		}(ch)
	}
	done := 0
	// Wait for workers to complete
	for done != jobs {
		select {
		case <-ch:
			done++
		}
	}
	return out.SubImage(bounds)
}

// GetImage returns the image pointed by path
func GetImage(path string) image.Image {
	file, err := os.Open("sample.jpg")
	if err != nil {
		log.Fatal(err)
	}
	image, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
	return image
}

// PutImage write the image to path
func PutImage(image image.Image, path string) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	err = jpeg.Encode(file, image, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// Script contains the state of a script as well as its operations
type Script struct {
	Steps       []*Step
	CurrentLine int
}

// Step contains the instructions to perform
type Step struct {
	Instructions []*Instruction
	Parent       *Script
	Input        string
	Output       string
}

// Instruction 
type Instruction struct {
	Argv      []string
	Operation *filters.Filter
	Parent    *Step
}

func NewInstruction(s *Step, tokens []string) (*Instruction, error) {
	res := Instruction{}

	res.Argv = tokens
	res.Parent = s

	return &res, nil
}

func NewStep(s *Script, tokens []string) (*Step, error) {
	res := Step{}

	if len(tokens) != 4 || tokens[2] != "as" || tokens[0] != "with" {
		return nil, errors.New("syntax error, expected syntax: with <input> as <output>")
	}

	res.Parent = s
	res.Input = tokens[1]
	res.Output = tokens[3]

	return &res, nil
}

func NewScript(path string) (*Script, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := bufio.NewReader(file)

	res := Script{}

	var line string
	var expect_step bool = true
	var current_step *Step = nil

	for {
		line, err = reader.ReadString('\n')

		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		log.Printf("processing %s", line)

		res.CurrentLine++
		tokens := strings.Split(strings.Trim(line, "\n"), " ")
		if len(tokens) == 0 || len(tokens[0]) == 0 || (len(tokens[0]) > 0 && tokens[0][0] == '#') {
			continue
		}
		if expect_step {
			new_step, err := NewStep(&res, tokens)
			if err != nil {
				return nil, err
			}
			res.Steps = append(res.Steps, new_step)
			current_step = new_step
			expect_step = false
		} else {
			if tokens[0] == "done" {
				current_step = nil
				expect_step = true
			} else {
				new_instr, err := NewInstruction(current_step, tokens)
				if err != nil {
					return nil, err
				}
				current_step.Instructions = append(current_step.Instructions, new_instr)
			}
		}
	}

	return &res, nil
}

func main() {
	flag.Parse()

	if len(*input_file) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	_, e := NewScript(*input_file)
	if e != nil {
		log.Fatalf(e.Error())
	}

	/* runtime.GOMAXPROCS(4) */
	/* input := GetImage("sample.jpg") */
	/* blur := filters.NewHBlur(20) */
	/* output := ProcessInParallel(input, 4, blur) */
	/* PutImage(output, "processed.jpg") */
}
