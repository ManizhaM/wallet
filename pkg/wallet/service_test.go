package wallet

import (
	"reflect"
	"testing"
	"github.com/ManizhaM/wallet/pkg/types"
)

func TestService_FindAccountByID_success(t *testing.T) {
	svc := &Service{}
	_, err := svc.RegisterAccount("+992000000001")
	acc, _ := svc.FindAccountByID(1)
	if err != nil {
		t.Error(err)
	}

	expected := &types.Account{ID: 1, Phone: "+992000000001", Balance: 0}
	if !reflect.DeepEqual(expected, acc) {
		t.Errorf("invalid result, expected: %v, actual: %v", expected, acc)
	}
}

func TestService_FindAccountByID_notFound(t *testing.T) {
	svc := &Service{}
	_, err := svc.FindAccountByID(1)
	expected := ErrAccountNotFound
	if !reflect.DeepEqual(expected, err) {
		t.Errorf("invalid result, expected: %v, actual: %v", expected, err)
	}
}