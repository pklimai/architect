// nolint: forbidigo
package logger

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var (
	infoLogLevel  = color.GreenString("[INFO] ")
	fatalLogLevel = color.RedString("[FATAL] ")
)

func Info(msg string) {
	fmt.Println(infoLogLevel, msg)
}

func Infof(format, msg string) {
	fmt.Println(infoLogLevel, fmt.Sprintf(format, msg))
}

func Fatal(msg string) {
	fmt.Println(fatalLogLevel, msg)
	os.Exit(1)
}

func Fatalf(format, msg string) {
	fmt.Println(fatalLogLevel, fmt.Sprintf(format, msg))
}

func FatalIfErr(err error) {
	if err != nil {
		Fatal(err.Error())
	}
}

func FatalfIfErr(err error, format string) {
	if err != nil {
		Fatal(fmt.Errorf(format, err).Error())
	}
}
