package lib

import (
	"bufio"
	clog "configTest/custom_logguer"
	"io"
	"os"
	"strings"
)

func readStream(filename string) (io.Reader, error) {
	var reader *bufio.Reader
	if filename == "-" {
		reader = bufio.NewReader(os.Stdin)
	} else {
		// ignore CWE-22 gosec issue - that's more targeted for http based apps that run in a public directory,
		// and ensuring that it's not possible to give a path to a file outside that directory.
		file, err := os.Open(filename) // #nosec
		if err != nil {
			return nil, err
		}
		reader = bufio.NewReader(file)
	}
	return reader, nil

}

func writeString(writer io.Writer, txt string) error {
	_, errorWriting := writer.Write([]byte(txt))
	return errorWriting
}

func safelyCloseFile(file *os.File) {
	err := file.Close()
	if err != nil {
		clog.Error("Error closing file!")
		clog.Error(err.Error())
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
