package model

import (
	"context"
	"errors"
	"fmt"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/pir5/pdns-api/models"
	"github.com/pir5/pir5-go/dnsapi/operations"
)

type RecordModel interface {
	ChangeStateToEnable() error
	ChangeStateToDisable() error
	GetState() (bool, error)
}

func NewRecordModel(id int64) RecordModel {
	return &Record{
		Addr: "127.0.0.1",
		Port: 8080,
		ID:   id,
	}
}

type Record struct {
	Addr string
	Port int
	ID   int64
}

func (r *Record) ChangeStateToEnable() error {
	fmt.Printf("[DEBUG] change state to enable (id: %d)\n", r.ID)
	transport := httptransport.New(fmt.Sprintf("%s:%d", r.Addr, r.Port), "v1", nil)
	p := &operations.PutRecordsEnableIDParams{
		ID: r.ID,
		Record: &models.ModelRecord{
			Disabled: false,
		},
		Context: context.Background(),
	}
	client := operations.New(transport, strfmt.Default)
	res, err := client.PutRecordsEnableID(p)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("[DEBUG] updated record: id: %d, disabled: %v\n", res.Payload.ID, res.Payload.Disabled)
	return err
}

func (r *Record) ChangeStateToDisable() error {
	fmt.Printf("[DEBUG] change state to disable (id: %d)\n", r.ID)
	transport := httptransport.New(fmt.Sprintf("%s:%d", r.Addr, r.Port), "v1", nil)
	p := &operations.PutRecordsDisableIDParams{
		ID: r.ID,
		Record: &models.ModelRecord{
			Disabled: true,
		},
		Context: context.Background(),
	}
	client := operations.New(transport, strfmt.Default)
	res, err := client.PutRecordsDisableID(p)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("[DEBUG] updated record: id: %d, disabled: %v\n", res.Payload.ID, res.Payload.Disabled)
	return err
}

func (r *Record) GetState() (bool, error) {
	transport := httptransport.New(fmt.Sprintf("%s:%d", r.Addr, r.Port), "v1", nil)
	p := &operations.GetRecordsParams{
		ID: &r.ID,
		Context: context.Background(),
	}
	client := operations.New(transport, strfmt.Default)
	record, err := client.GetRecords(p)
	if err != nil {
		return false, err
	}

	if len(record.GetPayload()) > 1 {
		return false, errors.New("Found records same ID")
	}

	return !record.GetPayload()[0].Disabled, nil
}
