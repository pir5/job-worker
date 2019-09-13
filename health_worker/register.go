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
	log "github.com/sirupsen/logrus"
	// Import godis package
)

const EnqueKey = "health_check"
const IntervalKey = "check_interval"

var CmdRegister = &Command{
	Run:       runRegister,
	UsageLine: "register",
	Short:     "Start Register Server",
	Long: `
Start Register Server
	`,
}

// runRegister executes sub command and return exit code.
func runRegister(cmdFlags *GlobalFlags, args []string) error {
	conf, err := setupWorkerComand(cmdFlags)
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

	healthCheck := model.NewHealthCheckModel(db)
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
				&v.Params,
			)
		}
		time.Sleep(time.Duration(conf.PollInterval) * time.Second)
	}
}
