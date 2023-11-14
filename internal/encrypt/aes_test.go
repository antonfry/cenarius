package encrypt

import (
	"testing"
)

func TestAESEncrypted(t *testing.T) {
	type args struct {
		decrypted string
		key       string
		iv        string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Valid",
			args: args{
				decrypted: "Valid test",
				key:       "f1c68defcac1715234f1b9a9906c0a7c",
				iv:        "f1c68defcac17152",
			},
			want:    "YiN/N9QtBJtmjYu6rxL3cA==",
			wantErr: false,
		},
		{
			name: "InValidKey",
			args: args{
				decrypted: "InValidKey test",
				key:       "f1c68defcac1715234f1b9a9906c0a7",
				iv:        "f1c68defcac17152",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AESEncrypted(tt.args.decrypted, tt.args.key, tt.args.iv)
			if (err != nil) != tt.wantErr {
				t.Errorf("AESEncrypted() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AESEncrypted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAESDecrypted(t *testing.T) {
	type args struct {
		encrypted string
		key       string
		iv        string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Valid",
			args: args{
				encrypted: "YiN/N9QtBJtmjYu6rxL3cA==",
				key:       "f1c68defcac1715234f1b9a9906c0a7c",
				iv:        "f1c68defcac17152",
			},
			want:    "Valid test",
			wantErr: false,
		},
		{
			name: "InValidKey",
			args: args{
				encrypted: "YiN/N9QtBJtmjYu6rxL3cA==",
				key:       "f1c68defcac1715234f1b9a9906c0a7",
				iv:        "f1c68defcac17152",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AESDecrypted(tt.args.encrypted, tt.args.key, tt.args.iv)
			if (err != nil) != tt.wantErr {
				t.Errorf("AESDecrypted() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AESDecrypted() = %v, want %v", got, tt.want)
			}
		})
	}
}
