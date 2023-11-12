package crypt

import "testing"

func TestEncryptMess(t *testing.T) {
	type args struct {
		mess     string
		encrType Enrypter
		opts     []Option
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "sha256 test with key",
			args: args{
				mess:     "test",
				encrType: Sha256Encoder{},
				opts:     []Option{WithKey("test key")},
			},
			want: "6PJVMUMjLXuYGlmET-Zr2lJGJYHYmM2PR-DSniaEalE=",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncryptMess(tt.args.mess, tt.args.encrType, tt.args.opts...); got != tt.want {
				t.Errorf("EncryptMess() = %v, want %v", got, tt.want)
			}
		})
	}
}
