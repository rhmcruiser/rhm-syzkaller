// Copyright 2016 syzkaller project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

// Package log provides functionality similar to standard log package with some extensions:
//   - verbosity levels
//   - global verbosity setting that can be used by multiple packages
//   - ability to disable all output
//   - ability to cache recent output in memory
package log

import (
	"bytes"
	"flag"
	"fmt"
	golog "log"
	"strings"
	"sync"
	"time"
)

var (
	flagV        = flag.Int("vv", 0, "verbosity")
	mu           sync.Mutex
	cacheMem     int
	cacheMaxMem  int
	cachePos     int
	cacheEntries []string
	instanceName string
	prependTime  = true // for testing
)

// EnableCaching enables in memory caching of log output.
// Caches up to maxLines, but no more than maxMem bytes.
// Cached output can later be queried with CachedOutput.
func EnableLogCaching(maxLines, maxMem int) {
	mu.Lock()
	defer mu.Unlock()
	if cacheEntries != nil {
		Fatalf("log caching is already enabled")
	}
	if maxLines < 1 || maxMem < 1 {
		panic("invalid maxLines/maxMem")
	}
	cacheMaxMem = maxMem
	cacheEntries = make([]string, maxLines)
}

// Retrieves cached log output.
func CachedLogOutput() string {
	mu.Lock()
	defer mu.Unlock()
	buf := new(bytes.Buffer)
	for i := range cacheEntries {
		pos := (cachePos + i) % len(cacheEntries)
		if cacheEntries[pos] == "" {
			continue
		}
		buf.WriteString(cacheEntries[pos])
		buf.Write([]byte{'\n'})
	}
	return buf.String()
}

// If the name is set, it will be displayed for all logs.
func SetName(name string) {
	instanceName = name
}

// V reports whether verbosity at the call site is at least the requested level.
// See https://pkg.go.dev/github.com/golang/glog#V for details.
func V(level int) bool {
	return level <= *flagV
}

func Logf(v int, msg string, args ...interface{}) {
	writeMessage(v, "", msg, args...)
}

func Errorf(msg string, args ...interface{}) {
	writeMessage(0, "ERROR", msg, args...)
}

func Fatal(err error) {
	golog.Fatal("SYZFATAL: ", err)
}

func Fatalf(msg string, args ...interface{}) {
	golog.Fatalf("SYZFATAL: "+msg, args...)
}

func writeMessage(v int, severity, msg string, args ...interface{}) {
	var sb strings.Builder
	if severity != "" {
		fmt.Fprintf(&sb, "[%s] ", severity)
	}
	if instanceName != "" {
		fmt.Fprintf(&sb, "%s: ", instanceName)
	}
	fmt.Fprintf(&sb, msg, args...)
	writeRawMessage(v, sb.String())
}

func writeRawMessage(v int, msg string) {
	mu.Lock()
	if cacheEntries != nil && v <= 1 {
		cacheMem -= len(cacheEntries[cachePos])
		if cacheMem < 0 {
			panic("log cache size underflow")
		}
		timeStr := ""
		if prependTime {
			timeStr = time.Now().Format("2006/01/02 15:04:05 ")
		}
		cacheEntries[cachePos] = timeStr + msg
		cacheMem += len(cacheEntries[cachePos])
		cachePos++
		if cachePos == len(cacheEntries) {
			cachePos = 0
		}
		for i := 0; i < len(cacheEntries)-1 && cacheMem > cacheMaxMem; i++ {
			pos := (cachePos + i) % len(cacheEntries)
			cacheMem -= len(cacheEntries[pos])
			cacheEntries[pos] = ""
		}
		if cacheMem < 0 {
			panic("log cache size underflow")
		}
	}
	mu.Unlock()

	if V(v) {
		golog.Print(msg)
	}
}

type VerboseWriter int

func (w VerboseWriter) Write(data []byte) (int, error) {
	Logf(int(w), "%s", data)
	return len(data), nil
}
