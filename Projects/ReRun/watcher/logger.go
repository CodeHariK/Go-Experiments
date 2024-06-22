package watcher

import (
	"os"
)

type stdOutSave struct {
	savedOutput []byte

	fn func(string)
}

func (so *stdOutSave) Write(p []byte) (n int, err error) {
	so.savedOutput = append(so.savedOutput, p...)

	so.fn(string(p))

	return os.Stdout.Write(p)
}

type stdErrSave struct {
	savedOutput []byte

	fn func(string)
}

func (so *stdErrSave) Write(p []byte) (n int, err error) {
	so.savedOutput = append(so.savedOutput, p...)

	so.fn(string(p))

	return os.Stderr.Write(p)
}
