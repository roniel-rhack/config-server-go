package lib

import (
	"bufio"
	clog "configTest/custom_logguer"
	"io"
	"os"
	"strings"
)

type readCloser struct {
	io.Reader
	closer io.Closer
}

func (rc *readCloser) Close() error {
	if rc.closer != nil {
		return rc.closer.Close()
	}
	return nil
}

func readStream(filename string) (*readCloser, error) {
	if filename == "-" {
		return &readCloser{Reader: bufio.NewReader(os.Stdin)}, nil
	}
	// ignore CWE-22 gosec issue - that's more targeted for http based apps that run in a public directory,
	// and ensuring that it's not possible to give a path to a file outside that directory.
	file, err := os.Open(filename) // #nosec
	if err != nil {
		return nil, err
	}
	return &readCloser{Reader: bufio.NewReader(file), closer: file}, nil
}

func writeString(writer io.Writer, txt string) error {
	_, errorWriting := writer.Write([]byte(txt))
	return errorWriting
}

func safelyCloseFile(file *os.File) {
	err := file.Close()
	if err != nil {
		clog.Error("Error closing file: %s", err.Error())
	}
}

func isTruthyNode(node *CandidateNode) bool {
	if node == nil {
		return false
	}
	if node.Tag == "!!null" {
		return false
	}
	if node.Kind == ScalarNode && node.Tag == "!!bool" {
		// yes/y/true/on
		return strings.EqualFold(node.Value, "y") ||
			strings.EqualFold(node.Value, "yes") ||
			strings.EqualFold(node.Value, "on") ||
			strings.EqualFold(node.Value, "true")

	}
	return true
}
