package delay_server

import (
	"reflect"
	"testing"
)

func TestNewDelayServer(t *testing.T) {
	type args struct {
		opt []SOption
	}
	tests := []struct {
		name    string
		args    args
		want    *delayServer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewDelayServer(tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDelayServer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDelayServer() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_delayServer_PutMsg(t *testing.T) {
	type fields struct {
		storagePath string
		ring        *ring
	}
	type args struct {
		msg DelayMessage
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &delayServer{
				storagePath: tt.fields.storagePath,
				ring:        tt.fields.ring,
			}
			if err := s.PutMsg(tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("PutMsg() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_delayServer_validate(t *testing.T) {
	type fields struct {
		storagePath string
		ring        *ring
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "empty",
			fields: fields{
				storagePath: "",
				ring:        nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &delayServer{
				storagePath: tt.fields.storagePath,
				ring:        tt.fields.ring,
			}
			if err := s.validate(); (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
