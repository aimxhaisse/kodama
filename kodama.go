package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"io"
	_ "github.com/aimxhaisse/kodama/cr2"
	"github.com/aimxhaisse/kodama/filters"
	"log"
	"os"
	"runtime"
	"strings"
)

var input_file = flag.String("infile", "", "input file")

// GetImage returns the image pointed by path
func GetImage(p string) (*filters.FilterImage, error) {
	file, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	img64 := image.NewRGBA64(img.Bounds())
	draw.Draw(img64, img64.Bounds(), img, image.Point{0, 0}, draw.Over)
	return &filters.FilterImage{img64}, nil
}

// PutImage write the image to path
func PutImage(image *filters.FilterImage, path string) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	err = jpeg.Encode(file, image.Image, nil)
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
	Id           int
}

// Instruction 
type Instruction struct {
	Argv      []string
	Operation filters.Filter
	Parent    *Step
	Id        int
}

// NewInstruction creates a new instruction for the given parameters
func NewInstruction(s *Step, tokens []string, id int) (*Instruction, error) {
	res := Instruction{}

	res.Argv = tokens
	res.Parent = s
	res.Id = id

	op := tokens[0]

	var err error

	switch op {

	case "blur":
		res.Operation, err = filters.NewBlur(tokens)
		if err != nil {
			s.Parent.Error(fmt.Sprintf("can't create blur: %s", err.Error()))
		}

	case "vblur":
		res.Operation, err = filters.NewVBlur(tokens)
		if err != nil {
			return nil, s.Parent.Error(fmt.Sprintf("can't create vblur: %s", err.Error()))
		}

	case "hblur":
		res.Operation, err = filters.NewHBlur(tokens)
		if err != nil {
			return nil, s.Parent.Error(fmt.Sprintf("can't create hblur: %s", err.Error()))
		}

	case "brightness":
		res.Operation, err = filters.NewBrightness(tokens)
		if err != nil {
			return nil, s.Parent.Error(fmt.Sprintf("can't create brightness: %s", err.Error()))
		}

	case "darkness":
		res.Operation, err = filters.NewDarkness(tokens)
		if err != nil {
			return nil, s.Parent.Error(fmt.Sprintf("can't create darkness: %s", err.Error()))
		}

	case "saturation":
		res.Operation, err = filters.NewSaturation(tokens)
		if err != nil {
			return nil, s.Parent.Error(fmt.Sprintf("can't create saturation: %s", err.Error()))
		}

	case "resize":
		res.Operation, err = filters.NewResize(tokens)
		if err != nil {
			return nil, s.Parent.Error(fmt.Sprintf("can't create resize: %s", err.Error()))
		}

	case "merge":
		res.Operation, err = filters.NewMerge(tokens)
		if err != nil {
			return nil, s.Parent.Error(fmt.Sprintf("can't create merge: %s", err.Error()))
		}

	default:
		return nil, s.Parent.Error(fmt.Sprintf("unknown operation: %s", op))
	}

	return &res, nil
}

// NewStep creates a step for the given script
func NewStep(s *Script, tokens []string, id int) (*Step, error) {
	res := Step{}

	if len(tokens) != 4 || tokens[2] != "as" || tokens[0] != "with" {
		return nil, s.Error("syntax error, expected syntax: with <input> as <output>")
	}

	res.Parent = s
	res.Input = tokens[1]
	res.Output = tokens[3]
	res.Id = id

	return &res, nil
}

// NewScript creates, parses and returns a Script ready to be exercuted
func NewScript(file *os.File) (*Script, error) {
	reader := bufio.NewReader(file)

	res := Script{}

	var expect_step bool = true
	var current_step *Step = nil

	for {
		line, err := reader.ReadString('\n')

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
			new_step, err := NewStep(&res, tokens, len(res.Steps)+1)
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
				new_instr, err := NewInstruction(current_step, tokens, len(current_step.Instructions)+1)
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
		fmt.Printf("step %d/%d (<- %s)\n", cur_step.Id, len(s.Steps), cur_step.Input)
		img, err := GetImage(cur_step.Input)
		if err != nil {
			return errors.New(fmt.Sprintf("can't open input %s: %s", cur_step.Input, err.Error()))
		}

		for j := 0; j < len(cur_step.Instructions); j++ {
			cur_instr := cur_step.Instructions[j]
			fmt.Printf("\tinstruction %d/%d (%s)... ", cur_instr.Id, len(cur_step.Instructions), cur_instr.Argv[0])
			op := cur_instr.Operation
			err = op.Process(img)
			if err != nil {
				return errors.New(fmt.Sprintf("can't process operation %s: %s", cur_instr.Argv[0], err.Error()))
			}
			fmt.Printf("done\n")
		}
		PutImage(img, cur_step.Output)
		fmt.Printf("done (-> %s)\n", cur_step.Output)
	}
	return nil
}

func main() {
	flag.Parse()

	var in *os.File

	if len(*input_file) == 0 {
		in = os.Stdin
	} else {
		var err error
		in, err = os.Open(*input_file)
		if err != nil {
			log.Fatalf(err.Error())
		}
		defer in.Close()
	}

	s, e := NewScript(in)
	if e != nil {
		log.Fatalf(e.Error())
	}

	runtime.GOMAXPROCS(4)

	e = s.Execute()
	if e != nil {
		log.Fatalf(e.Error())
	}
}
