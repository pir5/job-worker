package model

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestTCPChecker_Check(t *testing.T) {
	type fields struct {
		params HealthCheckParams
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				params: HealthCheckParams{
					Addr: "127.0.0.1",
					Port: 3000,
				},
			},
		},
		{
			name: "ng",
			fields: fields{
				params: HealthCheckParams{
					Addr: "127.0.0.1",
					Port: 3001,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			done := make(chan struct{})
			l, err := net.Listen("tcp", "127.0.0.1:3000")
			if err != nil {
				t.Fatal(err)
			}
			defer l.Close()

			go func() {
				conn, err := l.Accept()
				if err != nil {
					if strings.Index(err.Error(), "use of closed network connection") < 0 {
						t.Errorf("TCPChecker.Check() accept error = %v", err)
					}
					done <- struct{}{}
					return
				}
				done <- struct{}{}
				defer conn.Close()
			}()

			if err := TCPCheck(&tt.fields.params); (err != nil) != tt.wantErr {
				t.Errorf("TCPChecker.Check() error = %v, wantErr %v", err, tt.wantErr)
			}
			l.Close()
			<-done
		})
	}
}

var testHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello HTTP Test")
})

func Test_httpCheck(t *testing.T) {
	type args struct {
		params   *HealthCheckParams
		protocol string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				params:   &HealthCheckParams{},
				protocol: "http",
			},
		},
		{
			name: "ng",
			args: args{
				params:   &HealthCheckParams{},
				protocol: "https",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		ts := httptest.NewServer(testHandler)
		defer ts.Close()
		u, err := url.Parse(ts.URL)
		if err != nil {
			log.Fatal(err)
		}
		tt.args.params.Addr = strings.Split(u.Host, ":")[0]
		n, _ := strconv.Atoi(strings.Split(u.Host, ":")[1])
		tt.args.params.Port = n

		t.Run(tt.name, func(t *testing.T) {
			if err := HTTPCheck(tt.args.params, tt.args.protocol); (err != nil) != tt.wantErr {
				t.Errorf("httpCheck() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
