package model

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
)

func TestRoutingPolicy_FindBy(t *testing.T) {
	type fields struct {
		ID            int
		RecordID      int
		HealthCheckID int
		Type          int
	}
	type args struct {
		params map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    RoutingPolicies
		retRows *sqlmock.Rows
		wantErr bool
		retErr  error
	}{
		{
			name:   "ok",
			fields: fields{},
			args: args{
				params: map[string]interface{}{
					"id": 10,
				},
			},
			retRows: sqlmock.NewRows([]string{
				"id",
				"record_id",
				"health_check_id",
				"type",
			}).
				AddRow(
					1,
					2,
					3,
					4,
				),
			want: RoutingPolicies{
				RoutingPolicy{
					ID:            1,
					RecordID:      2,
					HealthCheckID: 3,
					Type:          4,
				},
			},
		},
		{
			name:   "notfound",
			fields: fields{},
			args: args{
				params: map[string]interface{}{
					"id": 1,
				},
			},
			retErr: gorm.ErrRecordNotFound,
			want:   nil,
		},
		{
			name:   "other error",
			fields: fields{},
			args: args{
				params: map[string]interface{}{
					"id": 1,
				},
			},
			retErr:  gorm.ErrInvalidSQL,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			if tt.retErr == nil {
				mock.ExpectQuery("SELECT \\* FROM `routing_policies` WHERE \\(id in\\(\\?\\)\\)").
					WithArgs(10).
					WillReturnRows(tt.retRows)
			} else {
				mock.ExpectQuery("SELECT \\* FROM `routing_policies` WHERE \\(id in\\(\\?\\)\\)").
					WillReturnError(tt.retErr)
			}

			gdb, _ := gorm.Open("mysql", db)
			h := &RoutingPolicy{
				db:            gdb,
				ID:            tt.fields.ID,
				RecordID:      tt.fields.RecordID,
				HealthCheckID: tt.fields.HealthCheckID,
				Type:          tt.fields.Type,
			}

			got, err := h.FindBy(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("RoutingPolicy.FindBy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RoutingPolicy.FindBy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRoutingPolicy_ChangeState(t *testing.T) {
	tests := []struct {
		name          string
		checkResut    bool
		healthCheckID int
	}{
		{
			name:          "RoutingPolicyTypeFailOverPrimary, checkResult is true",
			healthCheckID: RoutingPolicyTypeFailOverPrimary,
			checkResut:    true,
		},
		{
			name:          "RoutingPolicyTypeFailOverPrimary, checkResult is false",
			healthCheckID: RoutingPolicyTypeFailOverPrimary,
			checkResut:    false,
		},
		{
			name:          "RoutingPolicyTypeDetach, checkResult is true",
			healthCheckID: RoutingPolicyTypeDetach,
			checkResut:    true,
		},
		{
			name:          "RoutingPolicyTypeDetach, checkResult is false",
			healthCheckID: RoutingPolicyTypeDetach,
			checkResut:    false,
		},
		{
			name:          "RoutingPolicyTypeFailOverSecondly, checkResult is true",
			healthCheckID: RoutingPolicyTypeFailOverSecondly,
			checkResut:    true,
		},
		{
			name:          "RoutingPolicyTypeFailOverSecondly, checkResult is false",
			healthCheckID: RoutingPolicyTypeFailOverSecondly,
			checkResut:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &RoutingPolicy{}
			if err := h.ChangeState(tt.checkResut); err != nil {
				t.Errorf("RoutingPolicy.ChangeState failed: %s", err)
			}
		})
	}
}
