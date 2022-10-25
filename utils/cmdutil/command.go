package cmdutil

import (
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

// Duplicate logic in function bodies is intentional.
// This is to keep frame count the same (can't call one from the other), so that the caller may be logged.

// SudoExec executes external command with sudo privileges and logs output on the debug level.
// It returns a combined stderr and stdout output and exit code in case of an error.
func SudoExec(args ...string) error {
	out, err := exec.Command("sudo", args...).CombinedOutput()
	logSkipFrame := log.With().CallerWithSkipFrameCount(3).Logger()
	(&logSkipFrame).Debug().Msgf("%q output:\n%s", strings.Join(args, " "), out)
	return errors.Wrapf(err, "%q: %v output: %s", strings.Join(args, " "), err, out)
}

// Exec executes external command and logs output on the debug level.
// It returns a combined stderr and stdout output and exit code in case of an error.
func Exec(args ...string) error {
	out, err := exec.Command(args[0], args[1:]...).CombinedOutput()
	logSkipFrame := log.With().CallerWithSkipFrameCount(3).Logger()
	(&logSkipFrame).Debug().Msgf("%q output:\n%s", strings.Join(args, " "), out)
	return errors.Wrapf(err, "%q: %v output: %s", strings.Join(args, " "), err, out)
}

// ExecOutput executes external command and logs output on the debug level.
// It returns a combined stderr and stdout output and exit code in case of an error.
func ExecOutput(args ...string) (output string, err error) {
	out, err := exec.Command(args[0], args[1:]...).CombinedOutput()
	logSkipFrame := log.With().CallerWithSkipFrameCount(3).Logger()
	(&logSkipFrame).Debug().Msgf("%q output:\n%s", strings.Join(args, " "), out)
	if err != nil {
		return string(out), errors.Errorf("%q: %v output: %s", strings.Join(args, " "), err, out)
	}
	return string(out), nil
}
