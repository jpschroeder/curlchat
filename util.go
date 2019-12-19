package main

import (
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type WriteFlusher interface {
	io.Writer
	http.Flusher
}

type WriteFlush struct {
	writer  io.Writer
	flusher http.Flusher
}

func (wf WriteFlush) Write(buffer []byte) (int, error) {
	return wf.writer.Write(buffer)
}

func (wf WriteFlush) Flush() {
	wf.flusher.Flush()
}

var rex = regexp.MustCompile(`curl\/(?P<major>\d+)\.(?P<minor>\d+)\.(?P<patch>\d+)`)

// Parse user agent for a version of curl less than 7.68
func isOldCurl(agent string) bool {
	match := rex.FindStringSubmatch(agent)
	if match == nil {
		return false
	}
	major, maerr := strconv.ParseInt(match[1], 10, 32)
	minor, mierr := strconv.ParseInt(match[2], 10, 32)
	if maerr != nil || mierr != nil {
		return false
	}
	if major < 7 {
		return true
	}
	if major > 7 {
		return false
	}
	if minor < 68 {
		return true
	}
	return false
}

func isTerminal(agent string) bool {
	return strings.HasPrefix(agent, "curl")
}
