package model

import (
	"context"
	"errors"
	"fmt"
	"github.com/pir5/pdns-api/models"
	"github.com/pir5/pir5-go/dnsapi/operations"
)

type RecordModel interface {
	ChangeStateToEnable() error
	ChangeStateToDisable() error
	GetState() (bool, error)
}

func NewRecordModel(id int64, client *operations.Client) RecordModel {
	return &Record{
		ID:     id,
		client: client,
	}
}

type Record struct {
	ID     int64
	client *operations.Client
}

func (r *Record) ChangeStateToEnable() error {
	fmt.Printf("[DEBUG] change state to enable (id: %d)\n", r.ID)
	p := &operations.PutRecordsEnableIDParams{
		ID: r.ID,
		Record: &models.ModelRecord{
			Disabled: false,
		},
		Context: context.Background(),
	}

	res, err := r.client.PutRecordsEnableID(p)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("[DEBUG] updated record: id: %d, disabled: %v\n", res.Payload.ID, res.Payload.Disabled)
	return err
}

func (r *Record) ChangeStateToDisable() error {
	fmt.Printf("[DEBUG] change state to disable (id: %d)\n", r.ID)
	p := &operations.PutRecordsDisableIDParams{
		ID: r.ID,
		Record: &models.ModelRecord{
			Disabled: true,
		},
		Context: context.Background(),
	}
	res, err := r.client.PutRecordsDisableID(p)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("[DEBUG] updated record: id: %d, disabled: %v\n", res.Payload.ID, res.Payload.Disabled)
	return err
}

func (r *Record) GetState() (bool, error) {
	p := &operations.GetRecordsParams{
		ID:      &r.ID,
		Context: context.Background(),
	}
	record, err := r.client.GetRecords(p)
	if err != nil {
		return false, err
	}

	if len(record.GetPayload()) > 1 {
		return false, errors.New("Found records same ID")
	}

	return !record.GetPayload()[0].Disabled, nil
}
