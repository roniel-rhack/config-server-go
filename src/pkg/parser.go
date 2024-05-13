package lib

import (
	"bytes"
	clog "configTest/custom_logguer"
	"container/list"
)

type Parser interface {
	ResultsToMap(matchingNodes *list.List) (map[string]string, error)
}

type resultsParser struct {
	encoder           Encoder
	firstTimePrinting bool
	previousDocIndex  uint
	previousFileIndex int
	printedMatches    bool
}

func NewParser(encoder Encoder) Parser {
	return &resultsParser{
		encoder:           encoder,
		firstTimePrinting: true,
	}
}

func (p *resultsParser) ResultsToMap(matchingNodes *list.List) (map[string]string, error) {
	results := make(map[string]string)

	if matchingNodes.Len() == 0 {
		clog.Debug("no matching results, nothing to print")
		return results, nil
	}

	if !p.encoder.CanHandleAliases() {
		context := Context{MatchingNodes: matchingNodes}

		matchingNodes = context.MatchingNodes
	}

	if p.firstTimePrinting {
		node := matchingNodes.Front().Value.(*CandidateNode)
		p.previousDocIndex = node.GetDocument()
		p.previousFileIndex = node.GetFileIndex()
		p.firstTimePrinting = false
	}

	for el := matchingNodes.Front(); el != nil; el = el.Next() {
		mappedDoc := el.Value.(*CandidateNode)

		// Create a buffer to capture the output of the PrintLeadingContent method
		leadingContentBuffer := &bytes.Buffer{}
		if err := p.encoder.PrintLeadingContent(leadingContentBuffer, mappedDoc.LeadingContent); err != nil {
			return nil, err
		}

		toMap, err := p.encoder.EncodeToMap(mappedDoc)
		if err != nil {
			return nil, err
		}

		for k, v := range toMap {
			results[k] = v
		}
	}

	return results, nil
}
