package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/jrallison/go-workers"
	"github.com/pkg/errors"
)

const (
	HealthCheckTypeHTTP = iota
	HealthCheckTypeHTTPS
	HealthCheckTypeTCP
)

type HealthCheckModel interface {
	FindBy(map[string]interface{}) (HealthChecks, error)
}

func NewHealthCheckModel(db *gorm.DB) *HealthCheck {
	return &HealthCheck{
		db: db,
	}
}

type HealthCheck struct {
	db              *gorm.DB
	ID              int
	Name            string
	Type            int
	CheckInterval   int
	Threshould      int
	Params          healthCheckParams `gorm:"type:json"`
	RoutingPolicies *RoutingPolicies
}

type HealthChecks []HealthCheck

type healthCheckParams struct {
	protocol   int
	Addr       string
	Port       int
	HostName   string
	Path       string
	SearchWord string
	Timeout    time.Duration
}

// http://qiita.com/roothybrid7/items/2db3ccbf46f2bdb9cd00

func (s *healthCheckParams) ToJSON() (string, error) {
	r, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(r), nil
}

// Value SqlDriver interface:https://golang.org/pkg/database/sql/driver/#Valuer
func (s *healthCheckParams) Value() (driver.Value, error) {
	return s.ToJSON()
}

// Scan SqlDriver interface:https://golang.org/pkg/database/sql/#Scanner
func (s *healthCheckParams) Scan(value interface{}) (err error) {
	switch v := value.(type) {
	case string:
		if err := json.Unmarshal([]byte(v), s); err != nil {
			err = errors.New("spec.Scan: unmarshal json")
		}
	case []uint8:
		b := make([]byte, len(v))
		for i, a := range v {
			b[i] = byte(a)
		}

		if err := json.Unmarshal(b, s); err != nil {
			err = errors.New("spec.Scan: unmarshal json")
		}
	default:
		err = errors.New("spec.Scan: invalid value")
	}
	return nil
}

func NewHealthCheck(message *workers.Msg) (*HealthCheck, error) {
	b, err := message.Args().Encode()
	if err != nil {
		return nil, errors.Wrap(err, "job message encode failed")
	}

	p := HealthCheck{}
	if err := json.Unmarshal(b, &p); err != nil {
		return nil, errors.Wrap(err, "job message unmarshal failed")
	}
	return &p, nil
}

func (h *HealthCheck) FindBy(params map[string]interface{}) (HealthChecks, error) {
	query := h.db
	for k, v := range params {
		query = query.Where(k+" in(?)", v)
	}

	hs := HealthChecks{}
	r := query.Find(&hs)
	if r.Error != nil {
		if r.RecordNotFound() {
			return nil, nil
		} else {
			return nil, r.Error
		}
	}

	return hs, nil
}
