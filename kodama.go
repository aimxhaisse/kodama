package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"kodama/filters"
	"log"
	"os"
	"runtime"
	"fmt"
	"strings"
	"strconv"
)

// # TODO:
// - common interface to report errors (line, ...)
// - input/output
// - resizer
// - doc/with many samples
// - histo

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
func GetImage(path string) (*image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	image, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	file.Close()
	return &image, nil
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
	Operation filters.Filter
	Parent    *Step
}

// NewInstruction creates a new instruction for the given parameters
func NewInstruction(s *Step, tokens []string) (*Instruction, error) {
	res := Instruction{}

	res.Argv = tokens
	res.Parent = s

	op := tokens[0]

	switch op {

	case "blur":
		if len(tokens) != 2 {
			return nil, s.Parent.Error("invalid syntax for blur, expected usage: blur <radius>")
		}
		r, err := strconv.Atoi(tokens[1])
		if err != nil {
			return nil, s.Parent.Error(fmt.Sprintf("invalid parameter for blur: %s", err.Error()))
		}
		res.Operation, err = filters.NewBlur(r)
		if err != nil {
			return nil, s.Parent.Error(fmt.Sprintf("can't create blur: %s", err.Error()))
		}

	case "vblur":
		if len(tokens) != 2 {
			return nil, s.Parent.Error("invalid syntax for vblur, expected usage: vblur <strength>")
		}
		r, err := strconv.Atoi(tokens[1])
		if err != nil {
			return nil, s.Parent.Error(fmt.Sprintf("invalid parameter for vblur: %s", err.Error()))
		}
		res.Operation, err = filters.NewVBlur(r)
		if err != nil {
			return nil, s.Parent.Error(fmt.Sprintf("can't create vblur: %s", err.Error()))
		}

	case "hblur":
		if len(tokens) != 2 {
			return nil, s.Parent.Error("invalid syntax for hblur, expected usage: hblur <strength>")
		}
		r, err := strconv.Atoi(tokens[1])
		if err != nil {
			return nil, s.Parent.Error(fmt.Sprintf("invalid parameter for hblur: %s", err.Error()))
		}
		res.Operation, err = filters.NewHBlur(r)
		if err != nil {
			return nil, s.Parent.Error(fmt.Sprintf("can't create hblur: %s", err.Error()))
		}

	case "brightness":
		if len(tokens) != 2 {
			return nil, s.Parent.Error("invalid syntax for brightness, expected usage: brightness <strength>")
		}
		r, err := strconv.Atoi(tokens[1])
		if err != nil {
			return nil, s.Parent.Error(fmt.Sprintf("invalid parameter for brightness: %s", err.Error()))
		}
		res.Operation, err = filters.NewBrightness(r)
		if err != nil {
			return nil, s.Parent.Error(fmt.Sprintf("can't create brightness: %s", err.Error()))
		}

	default:
		return nil, s.Parent.Error(fmt.Sprintf("unknown operation: %s", op))
	}

	return &res, nil
}

// NewStep creates a step for the given script
func NewStep(s *Script, tokens []string) (*Step, error) {
	res := Step{}

	if len(tokens) != 4 || tokens[2] != "as" || tokens[0] != "with" {
		return nil, s.Error("syntax error, expected syntax: with <input> as <output>")
	}

	res.Parent = s
	res.Input = tokens[1]
	res.Output = tokens[3]

	return &res, nil
}

// NewScript creates, parses and returns a Script ready to be exercuted
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

		res.CurrentLine++
		tokens := strings.Split(strings.TrimSpace(strings.Trim(line, "\n")), " ")
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

// Error returns a new error with extra information about the context
func (s *Script) Error(e string) error {
	return errors.New(fmt.Sprintf("error on line %d: %s", s.CurrentLine, e))
}

// Execute executes the script
func (s *Script) Execute() error {
	for i := 0; i < len(s.Steps); i++ {
		cur_step := s.Steps[i]

		img, err := GetImage(cur_step.Input)
		if err != nil {
			return errors.New(fmt.Sprintf("can't open input %s: %s", cur_step.Input, err.Error()))
		}

		for j := 0; j < len(cur_step.Instructions); j++ {
			cur_instr := cur_step.Instructions[j]
			op := cur_instr.Operation
			if op.IsScalable() {
				*img = ProcessInParallel(*img, 4, op)
			} else {
				*img = ProcessInParallel(*img, 1, op)
			}
		}
		PutImage(*img, cur_step.Output)
	}
	return nil
}

func main() {
	flag.Parse()

	if len(*input_file) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	s, e := NewScript(*input_file)
	if e != nil {
		log.Fatalf(e.Error())
	}

	runtime.GOMAXPROCS(4)

	e = s.Execute()
	if e != nil {
		log.Fatalf(e.Error())
	}

	/* input := GetImage("sample.jpg") */
	/* blur := filters.NewHBlur(20) */
	/* output := ProcessInParallel(input, 4, blur) */
	/* PutImage(output, "processed.jpg") */
}
