package health_worker

import (
	"fmt"
	"os" // Import this package

	"github.com/pir5/health-worker/model"
	"github.com/pkg/errors"

	"github.com/facebookgo/pidfile"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/jrallison/go-workers"
	log "github.com/sirupsen/logrus"
	// Import godis package
)

var CmdWorker = &Command{
	Run:       runWorker,
	UsageLine: "worker",
	Short:     "Start Worker Server",
	Long: `
Start Worker Server
	`,
}
var globalDB *gorm.DB

// runWorker executes sub command and return exit code.
func runWorker(cmdFlags *GlobalFlags, args []string) error {
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
	globalDB = db

	workers.Process("health_check", do, conf.Concurrency)

	t := map[string]string{
		"hoge": "example",
	}
	workers.Enqueue("health_check",
		"Add",
		t,
	)
	// Blocks until process is told to exit via unix signal
	workers.Run()

	return nil
}

func do(msg *workers.Msg) {
	h, err := model.NewHealthCheck(msg)
	if err != nil || h == nil {
		log.Error(errors.Wrap(err, "parse params failed"))
	}

	if h != nil {
		switch h.Type {
		case model.HealthCheckTypeTCP:
			if err := model.TCPCheck(&h.Params); err != nil {
				log.Error(errors.Wrap(err, "tcp checker failed"))
			}
		case model.HealthCheckTypeHTTP:
			if err := model.HTTPCheck(&h.Params, "http"); err != nil {
				log.Error(errors.Wrap(err, "http checker failed"))
			}
		case model.HealthCheckTypeHTTPS:
			if err := model.HTTPCheck(&h.Params, "https"); err != nil {
				log.Error(errors.Wrap(err, "https checker failed"))
			}
		}

		if err := afterCheck(h, (err == nil)); err != nil {
			log.Error(errors.Wrap(err, "after check process failed"))
		}
	}
}

func afterCheck(h *model.HealthCheck, checkResult bool) error {
	// Todo: make redis counterr
	currentFailedCount := 100
	if currentFailedCount < h.Threshould {
		checkResult = true
	}

	rm := model.NewRoutingPolicyModel(globalDB)
	rs, err := rm.FindBy(map[string]interface{}{
		"health_check_id": h.ID,
	})

	if err != nil {
		return err
	}

	for _, r := range rs {
		err := r.ChangeState(checkResult)
		if err != nil {
			return err
		}
	}
	return nil
}
