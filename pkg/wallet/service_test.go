package wallet

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/ManizhaM/wallet/pkg/types"
)
type testService struct{
	*Service
}

func newTestService() *testService{
	return &testService{Service: &Service{}}
}

type testAccount struct{
	phone 		types.Phone
	balance		types.Money
	payments	[]struct{
		amount		types.Money
		category	types.PaymentCategory
	}
}

var defaultTestAccount = testAccount{
	phone: 		"+992000000001",
	balance: 	1000000,
	payments: 	[]struct{
		amount types.Money;
		category types.PaymentCategory
	}{
		{amount: 100000, category: "auto"},
	},
}	

// Функция добавления аккаунта и баланса в нем
func (s * testService) addAccount(data testAccount) (*types.Account, []*types.Payment, error){
	// регистрируем пользователя
	account, err := s.RegisterAccount(data.phone)
	if(err != nil){
		return nil, nil, fmt.Errorf("can't register account, error = %v", err)
	}

	//пополняем его счет
	err = s.Deposit(account.ID, data.balance)
	if(err != nil){
		return nil, nil, fmt.Errorf("can't deposit account, error = %v", err)
	}

	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("can't make payment, error = %v", err)
		}
	}

	return account, payments, nil
}
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
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	//ищем платеж
	payment := payments[0]
	got, err := s.FindPaymentByID(payment.ID)
	if err != nil{
		t.Errorf("FindPaymentByID(): error = %v", err)
		return
	}
	//сравниваем платежи
	if !reflect.DeepEqual(payment, got) {
		t.Errorf("FindPaymentByID(): wrong payment returned = %v", err)
		return
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

func TestService_Repeat_success(t *testing.T) {
		// создаем сервис
		s := newTestService()
		_, payments, err := s.addAccount(defaultTestAccount)
		if err != nil {
			t.Error(err)
			return
		}
		//осуществляем платеж на его счет
		payment := payments[0]
		repeated, err := s.Repeat(payment.ID)
		if err != nil{
			t.Errorf("Repeat(): error = %v", err)
			return
		}

		got, err := s.FindPaymentByID(repeated.ID)
		if err != nil{
			t.Errorf("FindPaymentByID(): error = %v", err)
		}
		//сравниваем платежи
		if !reflect.DeepEqual(repeated, got) {
			t.Errorf("FindPaymentByID(): wrong payment returned, error: %v", err)
		}
}