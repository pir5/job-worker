package health_worker

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/facebookgo/pidfile"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
)

type GlobalFlags struct {
	ConfPath *string
	PidPath  *string
	LogPath  *string
}

// A Command is an implementation of a pdns-api command
type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(cmdFlags *GlobalFlags, args []string) error

	// UsageLine is the one-line usage message.
	// The first word in the line is taken to be the command name.
	UsageLine string

	// Short is the short description shown in the 'pdns-api help' output.
	Short string

	// Long is the long message shown in the 'pdns-api help <this-command>' output.
	Long string

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet
}

// Name returns the command's name: the first word in the usage line.
func (c *Command) Name() string {
	name := c.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.UsageLine)
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(c.Long))
	os.Exit(2)
}

func initCommand(cmdFlags *GlobalFlags) (*Config, error) {
	conf, err := NewConfig(*cmdFlags.ConfPath)
	if err != nil {
		return nil, err
	}

	log.SetPrefix("health-worker")
	log.SetOutput(os.Stdout)
	pidfile.SetPidfilePath(*cmdFlags.PidPath)
	if *cmdFlags.LogPath != "" {
		f, err := os.OpenFile(*cmdFlags.LogPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return nil, errors.New("error opening file :" + err.Error())
		}
		log.SetOutput(f)
		log.SetLevel(log.WARN)
	} else {
		log.SetLevel(log.DEBUG)
	}

	if err := pidfile.Write(); err != nil {
		return nil, err
	}
	return &conf, nil
}

func getWorkerConfig(conf *Config) map[string]string {
	wc := map[string]string{
		"server":   fmt.Sprintf("%s:%d", conf.Redis.Host, conf.Redis.Port),
		"database": fmt.Sprintf("%d", conf.Redis.DB),
		// number of connections to keep open with redis
		"pool": fmt.Sprintf("%d", conf.Redis.PoolSize),
		// unique process id for this instance of workers (for proper recovery of inprogress jobs on crash)
		"process":       fmt.Sprintf("%d", conf.WorkerID),
		"poll_interval": fmt.Sprintf("%d", conf.PollInterval),
	}

	if conf.Redis.Password != "" {
		wc["password"] = conf.Redis.Password
	}

	return wc
}
