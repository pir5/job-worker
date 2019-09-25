package health_worker

import (
	"os"
	"reflect"
	"testing"
)

func TestNewConfig(t *testing.T) {
	type args struct {
		confPath string
	}
	tests := []struct {
		name    string
		args    args
		want    Config
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				confPath: "./test.toml",
			},
			want: Config{
				WorkerID:     0,
				PollInterval: 10,
				Concurrency:  10000,
				Listen:       "0.0.0.0:1000",
				DB: database{
					Host:     "127.0.0.1",
					Port:     3306,
					DBName:   "testdb",
					UserName: "testuser",
					Password: "db_password",
				},
				Redis: redis{
					Host:     "127.0.0.1",
					Port:     6379,
					PoolSize: 30,
					DB:       1,
					TTL:      60,
					Password: "redis_password",
				},
				PdnsAPI:pdnsAPI{
					Host: "127.0.0.1",
					Port: 8080,
				},
			},
		},
	}
	for _, tt := range tests {
		os.Setenv("PIR5_DATABASE_PASSWORD", "db_password")
		os.Setenv("PIR5_REDIS_PASSWORD", "redis_password")
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfig(tt.args.confPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
