package health_worker

import (
	"fmt"
	"os" // Import this package
	"time"

	"github.com/pir5/health-worker/model"
	"github.com/sonod/go-workers"

	"github.com/facebookgo/pidfile"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/gommon/log"
	// Import godis package
)

const EnqueKey = "health_check"
const IntervalKey = "check_interval"

var CmdRegister = &Command{
	Run:       runRegister,
	UsageLine: "register",
	Short:     "Start Register Server",
	Long:      `The register retrieves health checks that are registered in the database and enqueues tasks to the health checker. `,
}

// runRegister executes sub command and return exit code.
func runRegister(cmdFlags *GlobalFlags, args []string) error {
	conf, err := initCommand(cmdFlags)
	if err != nil {
		return err
	}
	defer func() {
		if err := os.Remove(pidfile.GetPidfilePath()); err != nil {
			log.Fatalf("Error removing %s: %s", pidfile.GetPidfilePath(), err)
		}
	}()

	workers.Configure(getWorkerConfig(conf))
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		conf.DB.UserName,
		conf.DB.Password,
		conf.DB.Host,
		conf.DB.Port,
		conf.DB.DBName,
	))

	if err != nil {
		return err
	}

	healthCheck := model.NewHealthCheckModeler(db)
	for {
		healthChecks, err := healthCheck.FindBy(map[string]interface{}{
			IntervalKey: conf.PollInterval,
		})

		if err != nil {
			return err
		}

		for _, v := range healthChecks {
			workers.Enqueue(EnqueKey,
				"Add",
				&v,
			)
		}

		time.Sleep(time.Duration(conf.PollInterval) * time.Second)
	}
}
