package spinner

import (
	"time"

	"github.com/briandowns/spinner"
)

// TODO: Definitely to refactor :)

// Spin is a global variable to avoid re-initializations
var Spin spinner.Spinner

// Init simply initializes the spinner with some default values
func Init() {
	Spin = *spinner.New(spinner.CharSets[11], 100*time.Millisecond, spinner.WithColor("magenta"), spinner.WithFinalMSG("✅\n"))
}

