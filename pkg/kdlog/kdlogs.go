package kdlog

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/spf13/pflag"
)

func Loop(f func(), intervalDuration time.Duration) {
	tickChan := time.NewTicker(intervalDuration)
	for {
		select {
		case <-tickChan.C:
			f()
		}
	}
}

const logFlushFreqFlagName = "log-flush-frequency"

var logFlushFreq = pflag.Duration(logFlushFreqFlagName, 5*time.Second, "Maximum number of seconds between log flushes")

var inited = false

// AddFlags registers this package's flags on arbitrary FlagSets, such that they point to the
// same value as the global flags.
func AddFlags(fs *pflag.FlagSet) {
	fs.AddFlag(pflag.Lookup(logFlushFreqFlagName))
}

// KlogWriter serves as a bridge between the standard log package and the glog package.
type KlogWriter struct{}

// Write implements the io.Writer interface.
func (writer KlogWriter) Write(data []byte) (n int, err error) {
	InfoDepth(1, string(data))
	return len(data), nil
}

// InitLogs initializes logs the way we want for kubernetes.
func InitLogs() {
	if !inited {
		InitFlags(flag.CommandLine)
		flag.Set("logtostderr", "false")
		inited = true
	}
	log.SetOutput(KlogWriter{})
	log.SetFlags(0)
	// The default glog flush interval is 5 seconds.
	go Loop(Flush, *logFlushFreq)
}

// FlushLogs flushes logs immediately.
func FlushLogs() {
	Flush()
}

// NewLogger creates a new log.Logger which sends logs to kdlog.Info.
func NewLogger(prefix string) *log.Logger {
	return log.New(KlogWriter{}, prefix, 0)
}

// GlogSetter is a setter to set glog level.
func GlogSetter(val string) (string, error) {
	var level Level
	if err := level.Set(val); err != nil {
		return "", fmt.Errorf("failed set kdlog.logging.verbosity %s: %v", val, err)
	}
	return fmt.Sprintf("successfully set kdlog.logging.verbosity to %s", val), nil
}
