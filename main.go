package main

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"unicode"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stderr)
}

var (
	left  = []byte{'{', '{'}
	right = []byte{'}', '}'}
)

const (
	StateOut = iota
	StateIn
)

var (
	state = StateOut
)

var (
	split bufio.SplitFunc = func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF {
			return 0, nil, nil
		}

		l := len(data)
		switch state {
		case StateOut:
			n := bytes.Index(data, left)
			if n < 0 {
				return l, data, nil
			}

			return n + 2, data[:n+2], nil
		case StateIn:
			n := bytes.Index(data, right)
			if n < 0 {
				return l, data, errors.New("unclosed brackets")
			}

			return n, data[:n], nil
		default:
			return 0, nil, errors.New("invalid state")
		}
	}
)

func handle(token []byte, w *bufio.Writer) error {
	trailing := false

	for i, b := range token {
		if trailing && unicode.IsSpace(rune(b)) {
			continue
		}

		if b == '\n' {
			token[i] = ' '
			trailing = true
		} else {
			trailing = false
		}

		err := w.WriteByte(token[i])
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	s.Split(split)

	w := bufio.NewWriter(os.Stdout)

	for s.Scan() {
		token := s.Bytes()
		switch state {
		case StateOut:
			_, err := w.Write(token)
			if err != nil {
				log.Fatal(err)
			}

			state = StateIn
		case StateIn:
			err := handle(token, w)
			if err != nil {
				log.Fatal(err)
			}

			state = StateOut
		default:
			err := errors.New("invalid state")
			log.Fatal(err)
		}
	}

	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	err := w.Flush()
	if err != nil {
		log.Fatal(err)
	}
}
