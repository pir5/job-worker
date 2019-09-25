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
		Params          *HealthCheckParams
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
					Params: &HealthCheckParams{
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

func TestHealthCheck_UpdateByID(t *testing.T) {
	type fields struct {
		db              *gorm.DB
		ID              int
		Name            string
		Type            int
		CheckInterval   int
		Threshould      int
		Params          *HealthCheckParams
		RoutingPolicies RoutingPolicies
	}
	type args struct {
		id             string
		newHealthCheck *HealthCheck
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		retErr          error
		healthCheckRows *sqlmock.Rows
		want            bool
		wantErr         bool
	}{
		{
			name:   "ok",
			fields: fields{},
			args: args{
				id: "1",
				newHealthCheck: &HealthCheck{
					Type: HealthCheckTypeTCP,
				},
			},
			healthCheckRows: sqlmock.NewRows([]string{
				"id",
				"name",
				"type",
				"check_interval",
				"threshould",
				"params",
			}).
				AddRow(
					1,
					"test check",
					2,
					10,
					3,
					`{ "addr": "test.com" }`,
				),
			want: true,
		},
		{
			name:   "notfound",
			fields: fields{},
			args: args{
				id: "2",
			},
			retErr: gorm.ErrRecordNotFound,
			want:   false,
		},
		{
			name:    "other error",
			fields:  fields{},
			retErr:  gorm.ErrInvalidSQL,
			want:    false,
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
				mock.ExpectQuery("SELECT \\* FROM `health_checks` WHERE \\(id = \\?\\) LIMIT 1").
					WithArgs("1").
					WillReturnRows(tt.healthCheckRows)
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `health_checks` SET `type` = \\? WHERE `health_checks`.`id` = \\?").
					WithArgs(tt.args.newHealthCheck.Type, 1).WillReturnResult(
					sqlmock.NewResult(
						1,
						1,
					),
				)

				mock.ExpectCommit()
			} else {
				mock.ExpectQuery("SELECT \\* FROM `health_checks` WHERE \\(id = \\?\\) LIMIT 1").
					WillReturnError(tt.retErr)
			}

			gdb, _ := gorm.Open("mysql", db)
			d := &HealthCheck{
				db:            gdb,
				ID:            tt.fields.ID,
				Name:          tt.fields.Name,
				Type:          tt.fields.Type,
				Threshould:    tt.fields.Threshould,
				CheckInterval: tt.fields.CheckInterval,
				Params:        tt.fields.Params,
			}

			got, err := d.UpdateByID(tt.args.id, tt.args.newHealthCheck)
			if (err != nil) != tt.wantErr {
				t.Errorf("HealthCheck.UpdateByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HealthCheck.UpdateByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHealthCheck_DeleteByID(t *testing.T) {
	type fields struct {
		db              *gorm.DB
		ID              int
		Name            string
		Type            int
		CheckInterval   int
		Threshould      int
		Params          *HealthCheckParams
		RoutingPolicies RoutingPolicies
	}
	type args struct {
		id string
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		retErr          error
		healthCheckRows *sqlmock.Rows
		want            bool
		wantErr         bool
	}{
		{
			name:   "ok",
			fields: fields{},
			args: args{
				id: "1",
			},
			healthCheckRows: sqlmock.NewRows([]string{
				"id",
				"name",
				"type",
				"check_interval",
				"threshould",
				"params",
			}).
				AddRow(
					1,
					"test check",
					2,
					10,
					3,
					`{ "addr": "test.com" }`,
				),
			want: true,
		},
		{
			name:   "notfound",
			fields: fields{},
			args: args{
				id: "2",
			},
			retErr: gorm.ErrRecordNotFound,
			want:   false,
		},
		{
			name:    "other error",
			fields:  fields{},
			retErr:  gorm.ErrInvalidSQL,
			want:    false,
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
				mock.ExpectQuery("SELECT \\* FROM `health_checks` WHERE \\(id = \\?\\) LIMIT 1").
					WithArgs("1").
					WillReturnRows(tt.healthCheckRows)
				mock.ExpectBegin()
				mock.ExpectExec("DELETE FROM `health_checks` WHERE `health_checks`.`id` = \\?").
					WithArgs(1).WillReturnResult(
					sqlmock.NewResult(
						1,
						1,
					),
				)

				mock.ExpectCommit()
			} else {
				mock.ExpectQuery("SELECT \\* FROM `health_checks` WHERE \\(id = \\?\\) LIMIT 1").
					WillReturnError(tt.retErr)
			}

			gdb, _ := gorm.Open("mysql", db)
			d := &HealthCheck{
				db:            gdb,
				ID:            tt.fields.ID,
				Name:          tt.fields.Name,
				Type:          tt.fields.Type,
				Threshould:    tt.fields.Threshould,
				CheckInterval: tt.fields.CheckInterval,
				Params:        tt.fields.Params,
			}

			got, err := d.DeleteByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("HealthCheck.DeleteByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HealthCheck.DeleteByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHealthCheck_Create(t *testing.T) {
	type fields struct {
		db              *gorm.DB
		ID              int
		Name            string
		Type            int
		CheckInterval   int
		Threshould      int
		Params          *HealthCheckParams
		RoutingPolicies RoutingPolicies
	}
	type args struct {
		newHealthCheck *HealthCheck
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		retErr  error
		wantErr bool
	}{
		{
			name:   "ok",
			fields: fields{},
			args: args{
				newHealthCheck: &HealthCheck{
					ID:   1,
					Name: "new health check",
				},
			},
		},
		{
			name:   "other error",
			fields: fields{},
			args: args{
				newHealthCheck: &HealthCheck{},
			},
			retErr:  gorm.ErrInvalidSQL,
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
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `health_checks` \\(`id`,`name`,`type`,`check_interval`,`threshould`,`params`\\) VALUES \\(\\?,\\?,\\?,\\?,\\?,\\?\\)").
					WithArgs(1, "new health check", 0, 0, 0, "null").WillReturnResult(
					sqlmock.NewResult(
						1,
						1,
					),
				)
				mock.ExpectCommit()
			} else {
				mock.ExpectExec("INSERT INTO `health_checks` \\(`id`,`name`,`type`,`check_interval`,`threshould`,`params`\\) VALUES \\(\\?,\\?,\\?,\\?,\\?,\\?\\)").
					WillReturnError(tt.retErr)
			}

			gdb, _ := gorm.Open("mysql", db)
			d := &HealthCheck{
				db:            gdb,
				ID:            tt.fields.ID,
				Name:          tt.fields.Name,
				Type:          tt.fields.Type,
				Threshould:    tt.fields.Threshould,
				CheckInterval: tt.fields.CheckInterval,
				Params:        tt.fields.Params,
			}

			err = d.Create(tt.args.newHealthCheck)
			if (err != nil) != tt.wantErr {
				t.Errorf("HealthCheck.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
