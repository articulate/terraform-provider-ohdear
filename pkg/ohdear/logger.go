package ohdear

import (
	"fmt"
	"log"
)

type TerraformLogger struct{}

func (l *TerraformLogger) log(level, message string) {
	log.Printf("[%s] %s\n", level, message)
}

func (l *TerraformLogger) Errorf(format string, v ...interface{}) {
	l.log("ERROR", fmt.Sprintf(format, v...))
}

func (l *TerraformLogger) Warnf(format string, v ...interface{}) {
	l.log("WARN", fmt.Sprintf(format, v...))
}

func (l *TerraformLogger) Debugf(format string, v ...interface{}) {
	l.log("DEBUG", fmt.Sprintf(format, v...))
}
