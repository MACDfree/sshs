package config

import (
	"reflect"
	"testing"
)

var data = `
alias1:
  ip: "192.168.123.123"
  port: 22
  username: "root"
  password: "123456"
alias2:
  ip: "192.168.123.123"
  port: 22
  username: "root"
  password: "123456"
`

func TestReadConfigData(t *testing.T) {
	wantMap := make(map[string]Session)
	wantMap["alias1"] = Session{
		IP:       "192.168.123.123",
		Port:     22,
		UserName: "root",
		Password: "123456",
	}
	wantMap["alias2"] = Session{
		IP:       "192.168.123.123",
		Port:     22,
		UserName: "root",
		Password: "123456",
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want map[string]Session
	}{
		{
			name: "test1",
			args: args{data: []byte(data)},
			want: wantMap,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readConfigData(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
