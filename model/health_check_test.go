package model

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
)

func TestHealthCheck_FindBy(t *testing.T) {
	type fields struct {
		db              *gorm.DB
		ID              int
		Name            string
		Type            int
		CheckInterval   int
		Threshould      int
		Params          healthCheckParams
		RoutingPolicies RoutingPolicies
	}
	type args struct {
		params map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    HealthChecks
		retRows *sqlmock.Rows
		wantErr bool
		retErr  error
	}{
		{
			name:   "ok",
			fields: fields{},
			args: args{
				params: map[string]interface{}{
					"check_interval": 10,
				},
			},
			retRows: sqlmock.NewRows([]string{
				"id",
				"name",
				"type",
				"check_interval",
				"threshould",
				"params",
			}).
				AddRow(
					1,
					"test",
					2,
					10,
					3,
					`{ "addr": "test.com" }`,
				),
			want: HealthChecks{
				HealthCheck{
					ID:            1,
					Name:          "test",
					Type:          2,
					CheckInterval: 10,
					Threshould:    3,
					Params: healthCheckParams{
						Addr: "test.com",
					},
				},
			},
		},
		{
			name:   "notfound",
			fields: fields{},
			args: args{
				params: map[string]interface{}{
					"check_interval": 1,
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
					"check_interval": 1,
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
				mock.ExpectQuery("SELECT \\* FROM `health_checks` WHERE \\(check_interval in\\(\\?\\)\\)").
					WithArgs(10).
					WillReturnRows(tt.retRows)
			} else {
				mock.ExpectQuery("SELECT \\* FROM `health_checks` WHERE \\(check_interval in\\(\\?\\)\\)").
					WillReturnError(tt.retErr)
			}

			gdb, _ := gorm.Open("mysql", db)

			h := &HealthCheck{
				db:            gdb,
				ID:            tt.fields.ID,
				Name:          tt.fields.Name,
				Type:          tt.fields.Type,
				Threshould:    tt.fields.Threshould,
				CheckInterval: tt.fields.CheckInterval,
				Params:        tt.fields.Params,
			}
			got, err := h.FindBy(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("HealthCheck.FindBy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HealthCheck.FindBy() = %v, want %v", got, tt.want)
			}
		})
	}
}
