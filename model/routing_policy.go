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

func NewRoutingPolicyModeler(db *gorm.DB, client *operations.Client) *RoutingPolicyModel {
	return &RoutingPolicyModel{
		db:     db,
		Client: client,
	}
}

type RoutingPolicyModel struct {
	db     *gorm.DB
	Client *operations.Client
}

type RoutingPolicy struct {
	ID            int
	RecordID      int
	HealthCheckID int
	Type          int
}

type RoutingPolicies []RoutingPolicy

type RoutingPolicyModeler interface {
	FindBy(map[string]interface{}) (RoutingPolicies, error)
	UpdateByID(string, *RoutingPolicy) (bool, error)
	DeleteByID(string) (bool, error)
	Create(*RoutingPolicy) error
	ChangeState(*RoutingPolicy, bool) error
}

func (rm *RoutingPolicyModel) ChangeState(r *RoutingPolicy, checkResult bool) error {
	// get state of records
	record := NewRecordModel(int64(r.RecordID), rm.Client)
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

func (h *RoutingPolicyModel) FindBy(params map[string]interface{}) (RoutingPolicies, error) {
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

func (d *RoutingPolicyModel) UpdateByID(id string, newRoutingPolicy *RoutingPolicy) (bool, error) {
	rp := RoutingPolicyModel{}
	r := d.db.Where("id = ?", id).Take(&rp)
	if r.Error != nil {
		if r.RecordNotFound() {
			return false, nil
		} else {
			return false, r.Error
		}
	}

	r = d.db.Model(&rp).Updates(&newRoutingPolicy)
	if r.Error != nil {
		return false, r.Error
	}
	return true, nil
}

func (d *RoutingPolicyModel) DeleteByID(id string) (bool, error) {
	rp := RoutingPolicy{}
	r := d.db.Where("id = ?", id).Take(&rp)
	if r.Error != nil {
		if r.RecordNotFound() {
			return false, nil
		} else {
			return false, r.Error
		}
	}

	r = d.db.Delete(&rp)
	if r.Error != nil {
		return false, r.Error
	}
	return true, nil
}

func (d *RoutingPolicyModel) Create(newRoutingPolicy *RoutingPolicy) error {
	if err := d.db.Create(newRoutingPolicy).Error; err != nil {
		return err
	}
	return nil
}
