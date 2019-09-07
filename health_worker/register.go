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
			"check_interval": conf.PollInterval,
		})
		if err != nil {
			return err
		}

		t := ""
		for _, v := range healthChecks {
			switch v.Type {
			case model.HealthCheckTypeTCP:
				t = "tcp_check"
			case model.HealthCheckTypeHTTP:
				t = "http_check"
			case model.HealthCheckTypeHTTPS:
				t = "https_check"
			}

			workers.Enqueue(t,
				"Add",
				&v.Params,
			)
		}
		time.Sleep(time.Duration(conf.PollInterval) * time.Second)
	}
}
