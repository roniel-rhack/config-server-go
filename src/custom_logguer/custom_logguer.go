package clog

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"os"
	"time"
)

var LOG *log.Logger

func Initialize() {

	styles := log.DefaultStyles()
	styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().
		SetString("ERROR").
		//Padding(0, 1, 0, 1).
		//Background(lipgloss.Color("204")).
		Foreground(lipgloss.Color("204"))

	LOG = log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    false,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
		//Prefix:          "Config Server => ",
	})

	LOG.SetStyles(styles)
}

func Info(format string, args ...interface{}) {
	LOG.Infof(format, args...)
}

func Warn(format string, args ...interface{}) {
	LOG.Warnf(format, args...)
}

func Print(format string, args ...interface{}) {
	LOG.Printf(format, args...)
}

func Debug(format string, args ...interface{}) {
	LOG.Debugf(format, args...)
}

func Error(format string, args ...interface{}) {
	LOG.Errorf(format, args...)
}
