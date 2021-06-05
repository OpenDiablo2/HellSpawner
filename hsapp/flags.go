package hsapp

import (
	"flag"
	"fmt"
	"os"

	"github.com/OpenDiablo2/HellSpawner/hsconfig"
)

// Flags specifies app flags
type Flags struct {
	optionalConfigPath *string
	bgColor            *string
	logFile            *string
}

// parse all of the command line args
func (a *App) parseArgs() {
	a.parseConfigArgs()
	a.parseLogFileArgs()
	a.parseBackgroundColorArgs()

	// help args need to be parsed last, so that all other args
	// can be parsed before a possible os.Exit() invoked by `-h` or `--help`
	//
	// otherwise, other flags will not be printed in usage string!
	a.parseHelpArgs()
}

func (a *App) parseHelpArgs() {
	const (
		short    = "h"
		long     = "help"
		fmtUsage = "usage: %s [<flags>]\n\nFlags:\n"
	)

	// we will use a single variable for both short and long flags
	var showHelp bool

	flag.BoolVar(&showHelp, long, false, "Show help")
	flag.BoolVar(&showHelp, short, false, "Show help (shorthand)")

	flag.Usage = func() {
		fmt.Printf(fmtUsage, os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if showHelp {
		flag.Usage()
		os.Exit(0) // this is dangerous, forces us to parse the help flags last
	}
}

func (a *App) parseConfigArgs() {
	const (
		name         = "config"
		defaultValue = ""
		fmtDesc      = "specify a custom config path.\nDefault is:\n\t%s"
	)

	desc := fmt.Sprintf(fmtDesc, hsconfig.GetConfigPath())
	a.Flags.optionalConfigPath = flag.String(name, defaultValue, desc)
}

func (a *App) parseBackgroundColorArgs() {
	const (
		name = "bgColor"
		desc = "custom background color."
	)

	defaultValue := fmt.Sprintf("0x%x", hsconfig.DefaultBGColor)
	a.Flags.bgColor = flag.String(name, defaultValue, desc)
}

func (a *App) parseLogFileArgs() {
	const (
		name = "log"
		desc = "path to the output log file."
	)

	a.Flags.logFile = flag.String(name, "", desc)
}
