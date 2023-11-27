package main

import (
	"encoding/json"
)


type GrepResult struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

func (p *GrepResult) UnmarshalJSON(data []byte) error {
	var typ struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(data, &typ); err != nil {
		return err
	}

	p.Type = typ.Type

	switch typ.Type {
	case "begin":
		p.Data = new(GrepBegin)
	case "match", "context":
		p.Data = new(GrepMatch)
	case "end":
		p.Data = new(GrepEnd)
	case "summary":
		p.Data = new(GrepSummary)
	}

	type tmp GrepResult
	return json.Unmarshal(data, (*tmp)(p))
}

type GrepText struct {
	Text string `json:"text"`
}
type GrepSubmatch struct {
	Match GrepText    `json:"match"`
	Start json.Number `json:"start"`
	End   json.Number `json:"end"`
}
type GrepDuration struct {
	Secs  json.Number `json:"secs"`
	Nanos json.Number `json:"nanos"`
	Human string      `json:"human"`
}
type GrepStats struct {
	Elapsed           GrepDuration `json:"elapsed"`
	Searches          json.Number  `json:"searches"`
	SearchesWithMatch json.Number  `json:"searches_with_match"`
	BytesSearched     json.Number  `json:"bytes_searched"`
	BytesPrinted      json.Number  `json:"bytes_printed"`
	MatchedLines      json.Number  `json:"matched_lines"`
	Matches           json.Number  `json:"matches"`
}
type GrepBegin struct {
	Path GrepText `json:"path"`
}
type GrepMatch struct {
	Path           GrepText       `json:"path"`
	Lines          GrepText       `json:"lines"`
	LineNumber     json.Number    `json:"line_number"`
	AbsoluteOffset json.Number    `json:"absolute_offset"`
	Submatches     []GrepSubmatch `json:"submatches"`
}
type GrepEnd struct {
	Path  GrepText  `json:"path"`
	Stats GrepStats `json:"stats"`
}
type GrepSummary struct {
	Path         GrepText     `json:"path"`
	ElapsedTotal GrepDuration `json:"elapsed_total"`
	Stats        GrepStats    `json:"stats"`
}
