package wrpc

import (
	"testing"

	. "github.com/kong/go-wrpc/wrpc/internal/wrpc"
)

func Test_validateEncoding(t *testing.T) {
	type args struct {
		e Encoding
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "protobuf encoding returns no error",
			args: args{
				e: Encoding_ENCODING_PROTO3,
			},
			wantErr: false,
		},
		{
			name: "unspecified encoding returns an error",
			args: args{
				e: Encoding_ENCODING_UNSPECIFIED,
			},
			wantErr: true,
		},
		{
			name: "random encoding returns an error",
			args: args{
				e: 23,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateEncoding(tt.args.e); (err != nil) != tt.wantErr {
				t.Errorf("validateEncoding() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
