package lib

import (
	clog "configTest/custom_logguer"
	"container/list"
	"errors"
	"fmt"
	"io"
)

// A yaml expression evaluator that runs the expression multiple times for each given yaml document.
// Uses less memory than loading all documents and running the expression once, but this cannot process
// cross document expressions.
type StreamEvaluator interface {
	EvaluateAndReturnMap(filename string, reader io.Reader, printer Parser, decoder Decoder) (uint, map[string]string, error)
	EvaluateFilesAndReturnMap(filenames []string, printer Parser, decoder Decoder) (map[string]string, error)
}

type streamEvaluator struct {
	fileIndex int
}

func NewStreamEvaluator() StreamEvaluator {
	return &streamEvaluator{}
}

func (s *streamEvaluator) EvaluateFilesAndReturnMap(filenames []string, printer Parser, decoder Decoder) (map[string]string, error) {
	results := make(map[string]string)

	for _, filename := range filenames {
		clog.Info("Reading file: %s", filename)
		reader, err := readStream(filename)
		if err != nil {
			return results, err
		}

		_, res, err := s.EvaluateAndReturnMap(filename, reader, printer, decoder)
		if err != nil {
			reader.Close()
			return results, err
		}

		for key, value := range res {
			results[key] = value
		}

		reader.Close()
	}

	return results, nil
}

func (s *streamEvaluator) EvaluateAndReturnMap(filename string, reader io.Reader, printer Parser, decoder Decoder) (uint, map[string]string, error) {

	var currentIndex uint
	err := decoder.Init(reader)
	if err != nil {
		return 0, nil, err
	}

	results := make(map[string]string)

	for {
		candidateNode, errorReading := decoder.Decode()

		if errors.Is(errorReading, io.EOF) {
			s.fileIndex++
			return currentIndex, results, nil
		} else if errorReading != nil {
			return currentIndex, results, fmt.Errorf("bad file '%v': %w", filename, errorReading)
		}
		candidateNode.document = currentIndex
		candidateNode.filename = filename
		candidateNode.fileIndex = s.fileIndex

		inputList := list.New()
		inputList.PushBack(candidateNode)

		result := Context{MatchingNodes: inputList}
		resultsToMap, mapErr := printer.ResultsToMap(result.MatchingNodes)

		if mapErr != nil {
			return currentIndex, results, mapErr
		}

		for key, value := range resultsToMap {
			results[key] = value
		}
		currentIndex++
	}
}
