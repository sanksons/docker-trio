package logger

import (
	"strings"
)

//struct for storing FormatType configuration
type LogFormatter struct {
	IsJson bool
	//Add other type checks here
}

//Populate the format config attributes
func (n *LogFormatter) Initialise(formatType string) {
	n.IsJson = strings.EqualFold(conf.FormatType, "json")
}
