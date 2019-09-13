package health_worker

import (
	"fmt"
	"os" // Import this package

	goredis "github.com/go-redis/redis"
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
var failedCounter model.FailedCounterModel

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

	redisClient := goredis.NewClient(&goredis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Redis.Host, conf.Redis.Port),
		Password: conf.Redis.Password,
		DB:       conf.Redis.DB,
	})
	if _, err := redisClient.Ping().Result(); err != nil {
		return err
	}

	failedCounter = model.NewFailedCounter(redisClient)

	workers.Process(EnqueKey, do, conf.Concurrency)

	// test data
	t := map[string]string{
		"hoge": "example",
	}
	workers.Enqueue(EnqueKey,
		"Add",
		t,
	)
	// Blocks until process is told to exit via unix signal
	workers.Run()

	return nil
}

func do(msg *workers.Msg) {
	var checkError error
	h, err := model.NewHealthCheck(msg)

	if err != nil || h == nil {
		log.Error(errors.Wrap(err, "parse params failed"))
	}

	if h != nil {
		switch h.Type {
		case model.HealthCheckTypeTCP:
			if checkError = model.TCPCheck(&h.Params); checkError != nil {
				log.Error(errors.Wrap(checkError, "tcp checker failed"))
			}
		case model.HealthCheckTypeHTTP:
			if checkError = model.HTTPCheck(&h.Params, "http"); checkError != nil {
				log.Error(errors.Wrap(checkError, "http checker failed"))
			}
		case model.HealthCheckTypeHTTPS:
			if checkError = model.HTTPCheck(&h.Params, "https"); checkError != nil {
				log.Error(errors.Wrap(checkError, "https checker failed"))
			}
		}

		if err := afterCheck(h, failedCounter, (checkError == nil)); err != nil {
			log.Error(errors.Wrap(err, "after check process failed"))
		}
	}
}

func afterCheck(h *model.HealthCheck, counter model.FailedCounterModel, checkResult bool) error {
	var currentFailedCount int
	key := fmt.Sprintf("failed_counter_%d", h.ID)
	if !checkResult {
		c, err := counter.Increment(key)
		if err != nil {
			return err
		}
		currentFailedCount = c
	} else {
		err := counter.Reset(key)
		if err != nil {
			return err
		}
	}

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
