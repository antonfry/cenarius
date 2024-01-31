package mcache

import (
	"cenarius/internal/model"
	"reflect"
	"testing"
)

func TestMCacheRepo(t *testing.T) {
	tests := []struct {
		name    string
		args    *model.SecretCache
		wantErr bool
	}{
		{
			name: "Valid",
			args: &model.SecretCache{
				LoginWithPasswords: []*model.LoginWithPassword{{Login: "TestLogin", Password: "TestPassword"}, {Login: "TestLogin2", Password: "TestPassword2"}},
				CreditCards:        []*model.CreditCard{{CVC: "323", Number: "239209355363"}, {CVC: "313", Number: "239376832093"}},
				SecretTexts:        []*model.SecretText{{Text: "Some Secret text"}, {Text: "another very secret test text"}},
				SecretFiles:        []*model.SecretFile{{Path: "/Some/Path"}, {Path: "another/path"}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MCacheRepo{
				store: New(),
			}
			if err := r.Save(tt.args); (err != nil) != tt.wantErr {
				t.Errorf("MCacheRepo.Save() error = %v, wantErr %v", err, tt.wantErr)
			}
			got, err := r.Get()
			if (err != nil) != tt.wantErr {
				t.Errorf("MCacheRepo.Save() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.args) {
				t.Errorf("PersistMetricRepository.Get() = %v, want %v\n", got, tt.args)
			}
		})
	}
}
