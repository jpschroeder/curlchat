package main

import (
	"io"
	"net/http"
	"regexp"
	"strconv"
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

type UserAgent struct {
	isCurl bool
	major  int64
	minor  int64
}

func parseAgent(agent string) UserAgent {
	match := rex.FindStringSubmatch(agent)
	if match == nil {
		return UserAgent{false, 0, 0}
	}
	major, _ := strconv.ParseInt(match[1], 10, 32)
	minor, _ := strconv.ParseInt(match[2], 10, 32)
	return UserAgent{true, major, minor}
}

// is user agent a version of curl less than 7.68
func (a UserAgent) isOldCurl() bool {
	if !a.isCurl {
		return false
	}
	if a.major < 7 {
		return true
	}
	if a.major > 7 {
		return false
	}
	if a.minor < 68 {
		return true
	}
	return false
}
