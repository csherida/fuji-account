package dynamoDB

import (
	"fuji-account/internal/models"
	"math/rand"
	"reflect"
	"strconv"
	"testing"
)

func Test_getItem(t *testing.T) {
	type args struct {
		fujiID string
	}
	tests := []struct {
		name    string
		args    args
		want    *models.FujiAccount
		wantErr bool
	}{
		{
			name: "Simple Account ID Test",
			args: args{
				fujiID: "1",
			},
			want: &models.FujiAccount{
				FujiID:      "1",
				AmazonToken: "1",
				AppleToken:  "1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetItem(tt.args.fujiID)
			if (err != nil) != tt.wantErr {
				t.Errorf("getItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getItem() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPutItem(t *testing.T) {
	type args struct {
		acct *models.FujiAccount
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Add new account",
			args: args{
				&models.FujiAccount{
					FujiID:      "Test" + strconv.Itoa(rand.Intn(1000)),
					AmazonToken: "Test",
					AppleToken:  "Test",
				},
			},
		},
		// TODO: Negative testing
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := PutItem(tt.args.acct); (err != nil) != tt.wantErr {
				t.Errorf("PutItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
