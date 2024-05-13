package lib

type PropertiesPreferences struct {
	UnwrapScalar      bool
	KeyValueSeparator string
	UseArrayBrackets  bool
}

func NewDefaultPropertiesPreferences() PropertiesPreferences {
	return PropertiesPreferences{
		UnwrapScalar:      true,
		KeyValueSeparator: ":",
		UseArrayBrackets:  true,
	}
}

func (p *PropertiesPreferences) Copy() PropertiesPreferences {
	return PropertiesPreferences{
		UnwrapScalar:      p.UnwrapScalar,
		KeyValueSeparator: p.KeyValueSeparator,
		UseArrayBrackets:  p.UseArrayBrackets,
	}
}

var ConfiguredPropertiesPreferences = NewDefaultPropertiesPreferences()
