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

type HealthCheckModeler interface {
	FindBy(map[string]interface{}) (HealthChecks, error)
	UpdateByID(string, *HealthCheck) (bool, error)
	DeleteByID(string) (bool, error)
	Create(*HealthCheck) error
}

func NewHealthCheckModeler(db *gorm.DB) *HealthCheckModel {
	return &HealthCheckModel{
		db: db,
	}
}

type HealthCheckModel struct {
	db *gorm.DB
}
type HealthCheck struct {
	ID              int                `json:"id"`
	Name            string             `json:"name"`
	Type            int                `json:"type"`
	CheckInterval   int                `json:"check_interval"`
	Threshould      int                `json:"threshould"`
	Params          *HealthCheckParams `json:"params" gorm:"type:json"`
	RoutingPolicies *RoutingPolicies   `json:"routing_polilcies"`
}

type HealthChecks []HealthCheck

type HealthCheckParams struct {
	protocol   int
	Addr       string
	Port       int
	HostName   string
	Path       string
	SearchWord string
	Timeout    time.Duration `swaggertype:"integer"`
}

// http://qiita.com/roothybrid7/items/2db3ccbf46f2bdb9cd00

func (s *HealthCheckParams) ToJSON() (string, error) {
	r, err := json.Marshal(s)
	if err != nil {
		return "", err
	}
	return string(r), nil
}

// Value SqlDriver interface:https://golang.org/pkg/database/sql/driver/#Valuer
func (s *HealthCheckParams) Value() (driver.Value, error) {
	return s.ToJSON()
}

// Scan SqlDriver interface:https://golang.org/pkg/database/sql/#Scanner
func (s *HealthCheckParams) Scan(value interface{}) (err error) {
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

func (h *HealthCheckModel) TableName() string {
	return "health_checks"
}
func (h *HealthCheckModel) FindBy(params map[string]interface{}) (HealthChecks, error) {
	query := h.db.New()
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

func (d *HealthCheckModel) UpdateByID(id string, newHealthCheck *HealthCheck) (bool, error) {
	hc := HealthCheck{}
	r := d.db.Where("id = ?", id).Take(&hc)
	if r.Error != nil {
		if r.RecordNotFound() {
			return false, nil
		} else {
			return false, r.Error
		}
	}

	r = d.db.Model(&hc).Updates(&newHealthCheck)
	if r.Error != nil {
		return false, r.Error
	}
	return true, nil
}
func (d *HealthCheckModel) DeleteByID(id string) (bool, error) {
	hc := HealthCheck{}
	r := d.db.Where("id = ?", id).Take(&hc)
	if r.Error != nil {
		if r.RecordNotFound() {
			return false, nil
		} else {
			return false, r.Error
		}
	}

	r = d.db.Delete(&hc)
	if r.Error != nil {
		return false, r.Error
	}
	return true, nil
}

func (d *HealthCheckModel) Create(newHealthCheck *HealthCheck) error {
	if err := d.db.Create(newHealthCheck).Error; err != nil {
		return err
	}
	return nil
}
