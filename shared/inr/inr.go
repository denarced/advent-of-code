package inr

import (
	"bytes"
	"io"
	"os"
	"strings"
)

type readOpt struct {
	includeSpace bool
	includeEmpty bool
}

type Option func(opt *readOpt)

func NoTrim() Option {
	return func(opt *readOpt) {
		opt.includeSpace = true
	}
}

func IncludeEmpty() Option {
	return func(opt *readOpt) {
		opt.includeEmpty = true
	}
}

func ReadPath(filep string, options ...Option) (lines []string, err error) {
	var f *os.File
	f, err = os.Open(filep)
	if err != nil {
		return
	}
	defer f.Close()

	var b []byte
	b, err = io.ReadAll(f)
	if err != nil {
		return
	}

	opt := new(readOpt)
	for _, each := range options {
		each(opt)
	}

	for _, each := range bytes.Split(b, []byte("\n")) {
		line := string(each)
		if !opt.includeSpace {
			line = strings.TrimSpace(line)
		}
		if !opt.includeEmpty && line == "" {
			continue
		}
		lines = append(lines, line)
	}
	return
}
