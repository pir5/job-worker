package model

import (
	"github.com/jinzhu/gorm"
)

const (
	RoutingPolicyTypeSimple = iota
	RoutingPolicyTypeFailOverPrimary
	RoutingPolicyTypeFailOverSecondly
	RoutingPolicyTypeDetach
)

func NewRoutingPolicyModel(db *gorm.DB) *RoutingPolicy {
	return &RoutingPolicy{
		db: db,
	}
}

type RoutingPolicy struct {
	db            *gorm.DB
	ID            int
	RecordID      int
	HealthCheckID int
	Type          int
}

type RoutingPolicies []RoutingPolicy

type RoutingPolicyModel interface {
	FindBy(map[string]interface{}) (RoutingPolicies, error)
	ChangeState(bool) error
}

func (r *RoutingPolicy) ChangeState(checkResult bool) error {
	// get state of records
	currentState := true
	switch r.HealthCheckID {
	case RoutingPolicyTypeFailOverPrimary, RoutingPolicyTypeDetach:
		// change state to disable if currentState is enable and checkResult is failed
		if currentState && !checkResult {
			// change state to enable if currentState is disable and checkResult is success
		} else if !currentState && checkResult {
		}
	case RoutingPolicyTypeFailOverSecondly:
		// change state to disable if currentState is enable and checkResult is success
		if checkResult && currentState {
			// change state to enable if currentState is disable and checkResult is failed
		} else if !checkResult && !currentState {
		}
	}

	return nil
}

func (h *RoutingPolicy) FindBy(params map[string]interface{}) (RoutingPolicies, error) {
	query := h.db
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
	r := d.db.Where("id = ?", id).Take(&d)
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
	r := d.db.Where("id = ?", id).Take(&d)
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
	if err := d.db.Create(newRoutingPolicy).Error; err != nil {
		return err
	}
	return nil
}
