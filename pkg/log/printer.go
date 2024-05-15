package log

import (
	"fmt"
	"github.com/fatih/color"
	"os"
)

var (
	displayErrors bool
)

func DisplayErrors(enabled bool) {
	displayErrors = enabled
}

func Write(contents string, c color.Attribute) {
	color.Set(c)
	fmt.Println(contents)
	color.Unset()
}

func Fatal(contents string) {
	Write(contents, color.FgRed)
	os.Exit(1)
}

func Error(err error, contents string) {
	Write(contents, color.FgRed)

	if displayErrors {
		fmt.Println(err)
	}

}
