package main

import (
	"io"
	"net/http"
	"regexp"
	"strconv"
)

type WriteFlusher struct {
	writer  io.Writer
	flusher http.Flusher
}

func (wf WriteFlusher) Write(buffer []byte) (int, error) {
	n, e := wf.writer.Write(buffer)
	wf.flusher.Flush()
	return n, e
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
