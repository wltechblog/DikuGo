package storage

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Parser is a utility for parsing DikuMUD data files
type Parser struct {
	scanner  *bufio.Scanner
	line     string
	lineNum  int
	file     *os.File
	filename string
	backedUp bool
}

// NewParser creates a new parser for the given file
func NewParser(filename string) (*Parser, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}

	scanner := bufio.NewScanner(file)
	return &Parser{
		scanner:  scanner,
		lineNum:  0,
		file:     file,
		filename: filename,
		backedUp: false,
	}, nil
}

// NextLine reads the next line from the file
func (p *Parser) NextLine() bool {
	// If we've backed up, just return true and reset the flag
	if p.backedUp {
		p.backedUp = false
		return true
	}

	// Otherwise, read the next line
	if !p.scanner.Scan() {
		return false
	}
	p.line = p.scanner.Text()
	p.lineNum++
	return true
}

// BackUp marks the parser as backed up, so the next call to NextLine will
// return the current line again
func (p *Parser) BackUp() {
	p.backedUp = true
}

// Line returns the current line
func (p *Parser) Line() string {
	return p.line
}

// LineNum returns the current line number
func (p *Parser) LineNum() int {
	return p.lineNum
}

// ReadString reads a string until the given delimiter
func (p *Parser) ReadString(delimiter string) (string, error) {
	var sb strings.Builder
	for p.NextLine() {
		if p.Line() == delimiter {
			return sb.String(), nil
		}
		if sb.Len() > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(p.Line())
	}
	return "", fmt.Errorf("unexpected end of file while reading string")
}

// ReadInt reads an integer from the current line
func (p *Parser) ReadInt() (int, error) {
	return strconv.Atoi(strings.TrimSpace(p.Line()))
}
