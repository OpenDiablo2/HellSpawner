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
	bgColor            *uint
}

func (a *App) parseArgs() {
	configDescr := fmt.Sprintf("specify a custom config path. Default is: %s", hsconfig.GetConfigPath())
	a.Flags.optionalConfigPath = flag.String("config", "", configDescr)
	a.Flags.bgColor = flag.Uint("bgColor", hsconfig.DefaultBGColor, "custom background color")
	showHelp := flag.Bool("h", false, "Show help")

	flag.Usage = func() {
		fmt.Printf("usage: %s [<flags>]\n\nFlags:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}
}
