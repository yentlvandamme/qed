package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Mode int

const (
	Replace Mode = iota
	Insert
	Append
)

type Arguments struct {
	Text       string
	File       io.Reader
	LineNumber int
	Mode       Mode
}

func main() {
	argsWithoutProg := os.Args[1:]
	argsLen := len(argsWithoutProg)

	if argsLen < 3 {
		fmt.Printf("Not enough arguments were passed. Expected at least 3 arguments but got %d.\n", argsLen)
		os.Exit(1)
	}

	// TODO: Make sure the file does have the required amount of lines to edit. If the user tries to edit the 10th line and there's only 5 present, exit the program.
	textToInsert := os.Args[1]
	fileName := os.Args[2]
	lineNumber, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Printf("Invalid line number.\n")
		os.Exit(1)
	}

	fileHandler, err := GetFile(fileName)
	if err != nil {
		fmt.Printf("Invalid file path.\n")
		os.Exit(1)
	}

	mode := Replace
	if argsLen > 3 {
		parsedMode, err := ParseMode(argsWithoutProg[4])
		if err != nil {
			fmt.Printf("Invalid mode value: %s.\n", argsWithoutProg[4])
			os.Exit(1)
		}

		mode = parsedMode
	}

	args := Arguments{
		Text:       textToInsert,
		File:       fileHandler,
		LineNumber: lineNumber,
		Mode:       mode,
	}

	// Make sure to close old file first.
	tempFile, err := CreateNewTargetFile(fileHandler.Name())
	if err != nil {
		fmt.Printf("Could not create temporary file %s\n", tempFile.Name())
		os.Exit(1)
	}

	err = UpdateContent(args, tempFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func GetFile(path string) (*os.File, error) {
	if filepath.IsAbs(path) {
		return os.OpenFile(path, syscall.O_RDWR, 0644)
	}

	base, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return os.OpenFile(filepath.Join(base, path), syscall.O_RDWR, 0644)
}

func ParseMode(mode string) (Mode, error) {
	if strings.EqualFold(mode, "replace") {
		return Replace, nil
	}
	if strings.EqualFold(mode, "append") {
		return Append, nil
	}
	if strings.EqualFold(mode, "insert") {
		return Insert, nil
	}

	return Append, fmt.Errorf("Could not parse mode: %s into supported mode.", mode)
}

// TODO: Fix trailing \n character. Whenever we execute an action, we append a \n at the end of the file.
func UpdateContent(args Arguments, target io.Writer) error {
	bufReader := bufio.NewReader(args.File)

	newLineBytes := []byte("\n")
	var currentLine int
	currentLine = 1

	for {
		line, _, err := bufReader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		if currentLine == args.LineNumber {
			if args.Mode == Append {
				target.Write(line)
				target.Write(newLineBytes)
				currentLine++

				target.Write([]byte(args.Text))
				target.Write(newLineBytes)
				currentLine++

				continue
			}
			if args.Mode == Insert {
				target.Write([]byte(args.Text))
				target.Write(newLineBytes)
				currentLine++

				target.Write(line)
				target.Write(newLineBytes)
				currentLine++

				continue
			}

			if args.Mode == Replace {
				target.Write([]byte(args.Text))
				target.Write(newLineBytes)
				currentLine++

				continue
			}
		}

		target.Write(line)
		target.Write(newLineBytes)
		currentLine++
	}

	return nil
}

func CreateNewTargetFile(baseName string) (*os.File, error) {
	currentDate := time.Now().Unix()
	name := baseName + "_" + strconv.Itoa(int(currentDate))
	return os.CreateTemp("", name)
}

func RemoveAndReplaceFile(old *os.File, new *os.File) error {
	return nil
}
