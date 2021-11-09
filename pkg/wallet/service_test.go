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


func TestService_Reject_success(t *testing.T) {
	s := &Service{}
	acc, err := s.RegisterAccount("+992000000001")
	if err != nil {
		t.Error(err)
	}
	err1 := s.Deposit(acc.ID, 2000)
	if err1 != nil {
		t.Error(err1)
	}
	pay, err2 := s.Pay(1, 1000, "auto")
	if err2 != nil {
		t.Error(err2)
	}
	err3 := s.Reject(pay.ID)

	//var expected *error
	if !reflect.DeepEqual(nil, err3) {
		t.Errorf("invalid result, expected: %v, actual: %v", nil, err3)
	}
}

func TestService_Reject_notFound(t *testing.T) {
	s := &Service{}
	err := s.Reject("1111")

	expected := ErrPaymentNotFound
	if !reflect.DeepEqual(expected, err) {
		t.Errorf("invalid result, expected: %v, actual: %v", expected, err)
	}
}


func TestService_FindPaymentByID_success(t *testing.T) {
	// создаем сервис
	s := newTestService()
	account, err := s.AddAccountWithBalance("+992000000001", 1000000)
	if err != nil {
		t.Error(err)
		return
	}
	//осуществляем платеж на его счет
	payment, err := s.Pay(account.ID, 100000, "auto")
	if err != nil {
		t.Errorf("Pay(): can't create payment, error = %v", err)
	}
	got, err := s.FindPaymentByID(payment.ID)
	if err != nil{
		t.Errorf("FindPaymentByID(): error = %v", err)
	}
	//сравниваем платежи
	if !reflect.DeepEqual(got, payment) {
		t.Errorf("invalid result, expected: %v, actual: %v", got, payment)
	}
}

func TestService_FindPaymentByID_notFound(t *testing.T) {
	s := &Service{}
	_, err := s.FindPaymentByID("1111")
	expected := ErrPaymentNotFound
	if !reflect.DeepEqual(expected, err) {
		t.Errorf("invalid result, expected: %v, actual: %v", expected, err)
	}
}