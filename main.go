package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

type Mode int

const (
	Replace Mode = iota
	Insert
	Append
)

type Arguments struct {
	Text       string
	File       *os.File
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
		fmt.Printf("Line number not valid.\n")
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

	_ = Arguments{
		Text:       textToInsert,
		File:       fileHandler,
		LineNumber: lineNumber,
		Mode:       mode,
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
