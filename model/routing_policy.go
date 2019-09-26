package model

import (
	"github.com/jinzhu/gorm"
	"github.com/pir5/pir5-go/dnsapi/operations"
)

const (
	RoutingPolicyTypeSimple = iota
	RoutingPolicyTypeFailOverPrimary
	RoutingPolicyTypeFailOverSecondly
	RoutingPolicyTypeDetach
)

func NewRoutingPolicyModel(db *gorm.DB, client *operations.Client) *RoutingPolicy {
	return &RoutingPolicy{
		db:     db,
		Client: client,
	}
}

type RoutingPolicy struct {
	db            *gorm.DB
	Client        *operations.Client
	ID            int
	RecordID      int
	HealthCheckID int
	Type          int
}

type RoutingPolicies []RoutingPolicy

type RoutingPolicyModel interface {
	FindBy(map[string]interface{}) (RoutingPolicies, error)
	UpdateByID(string, *RoutingPolicy) (bool, error)
	DeleteByID(string) (bool, error)
	Create(policy *RoutingPolicy) error
	ChangeState(bool) error
}

func (r *RoutingPolicy) ChangeState(checkResult bool) error {
	// get state of records
	record := NewRecordModel(int64(r.RecordID), r.Client)
	currentState, err := record.GetState()
	if err != nil {
		return err
	}

	switch r.Type {
	case RoutingPolicyTypeFailOverPrimary, RoutingPolicyTypeDetach:
		if currentState && !checkResult {
			// change state to disable if currentState is enable and checkResult is failed
			return record.ChangeStateToDisable()
		} else if !currentState && checkResult {
			// change state to enable if currentState is disable and checkResult is success
			return record.ChangeStateToEnable()
		}
	case RoutingPolicyTypeFailOverSecondly:
		if checkResult && currentState {
			// change state to disable if currentState is enable and checkResult is success
			return record.ChangeStateToDisable()
		} else if !checkResult && !currentState {
			return record.ChangeStateToEnable()
		}
	}

	return nil
}

func (h *RoutingPolicy) FindBy(params map[string]interface{}) (RoutingPolicies, error) {
	query := h.db.New()
	for k, v := range params {
		query = query.Where(k+" in(?)", v)
	}

	hs := RoutingPolicies{}
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

func (d *RoutingPolicy) UpdateByID(id string, newRoutingPolicy *RoutingPolicy) (bool, error) {
	r := d.db.New().Where("id = ?", id).Take(&d)
	if r.Error != nil {
		if r.RecordNotFound() {
			return false, nil
		} else {
			return false, r.Error
		}
	}

	r = d.db.Model(&d).Updates(&newRoutingPolicy)
	if r.Error != nil {
		return false, r.Error
	}
	return true, nil
}

func (d *RoutingPolicy) DeleteByID(id string) (bool, error) {
	r := d.db.New().Where("id = ?", id).Take(&d)
	if r.Error != nil {
		if r.RecordNotFound() {
			return false, nil
		} else {
			return false, r.Error
		}
	}

	r = d.db.Delete(d)
	if r.Error != nil {
		return false, r.Error
	}
	return true, nil
}

func (d *RoutingPolicy) Create(newRoutingPolicy *RoutingPolicy) error {
	if err := d.db.New().Create(newRoutingPolicy).Error; err != nil {
		return err
	}
	return nil
}
