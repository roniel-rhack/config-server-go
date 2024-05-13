package lib

import (
	"fmt"
	"strings"
)

type EncoderFactoryFunction func() Encoder
type DecoderFactoryFunction func() Decoder

type Format struct {
	FormalName     string
	Names          []string
	EncoderFactory EncoderFactoryFunction
	DecoderFactory DecoderFactoryFunction
}

var YamlFormat = &Format{"yaml", []string{"y", "yml"},
	func() Encoder { return nil },
	func() Decoder { return NewYamlDecoder(ConfiguredYamlPreferences) },
}

var PropertiesFormat = &Format{"props", []string{"p", "properties"},
	func() Encoder { return NewPropertiesEncoder(ConfiguredPropertiesPreferences) },
	func() Decoder { return nil },
}

var Formats = []*Format{
	YamlFormat,
	PropertiesFormat,
}

func (f *Format) MatchesName(name string) bool {
	if f.FormalName == name {
		return true
	}
	for _, n := range f.Names {
		if n == name {
			return true
		}
	}
	return false
}

func (f *Format) GetConfiguredEncoder() Encoder {
	return f.EncoderFactory()
}

func FormatFromString(format string) (*Format, error) {
	if format != "" {
		for _, printerFormat := range Formats {
			if printerFormat.MatchesName(format) {
				return printerFormat, nil
			}
		}
	}
	return nil, fmt.Errorf("unknown format '%v' please use [%v]", format, GetAvailableOutputFormatString())
}

func GetAvailableOutputFormats() []*Format {
	var formats = []*Format{}
	for _, printerFormat := range Formats {
		if printerFormat.EncoderFactory != nil {
			formats = append(formats, printerFormat)
		}
	}
	return formats
}

func GetAvailableOutputFormatString() string {
	var formats = []string{}
	for _, printerFormat := range GetAvailableOutputFormats() {

		if printerFormat.FormalName != "" {
			formats = append(formats, printerFormat.FormalName)
		}
		if len(printerFormat.Names) >= 1 {
			formats = append(formats, printerFormat.Names[0])
		}
	}
	return strings.Join(formats, "|")
}
