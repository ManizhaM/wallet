package wallet

import (
	"errors"
	"fmt"

	"github.com/ManizhaM/wallet/pkg/types"
)

var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greater than 0")
var ErrAccountNotFound = errors.New("account not found")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrNotEnoughBalance = errors.New("not enough balance")

type Service struct {
	nextAccountID int64
	accounts []*types.Account
	payments []*types.Payment
}

func (s *Service) RegisterAccount(phone types.Phone)(*types.Account, error){
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}

	s.nextAccountID++
	account := &types.Account{
		ID: s.nextAccountID,
		Phone: phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)
	return account, nil
}

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if (amount<=0) {
		return ErrAmountMustBePositive
	}
	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID{
			account = acc
			break
		}
	}
	if account == nil{
		return ErrAccountNotFound
	}
	account.Balance += amount
	return nil
}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error){
	if amount<=0{
		return nil, ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID==accountID {
			account=acc
			break
		}
	}
	if account == nil{
		return nil, ErrAccountNotFound
	}
	if account.Balance<amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= amount
	paymentID := "4122"
	payment := &types.Payment{
		ID: paymentID,
		AccountID: accountID,
		Amount: amount,
		Category: category,
		Status: types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil

}

//FindAccountByID возвращает указатель на найденный аккаунт 
func (s *Service) FindAccountByID(accountID int64)(*types.Account,error){
	var account *types.Account
	for _, acc := range s.accounts {
		if(acc.ID == accountID){
			account = acc
			return account, nil
		}
	}
	return nil, ErrAccountNotFound
}

func (s *Service) FindPaymentByID(paymentID string)(*types.Payment, error){
	var payment *types.Payment
	for _, pay := range s.payments {
		if pay.ID == paymentID {
			payment = pay
			return payment, nil	
		}	
	}
	return nil, ErrPaymentNotFound
}

//Функция отмены платежа
func (s *Service) Reject(paymentID string) error{
	payment, err := s.FindPaymentByID(paymentID)
	if payment != nil{
		account, err1 := s.FindAccountByID(payment.AccountID)
		if account == nil{
			return err1
		}
		payment.Status = types.PaymentStatusFail
		account.Balance += payment.Amount
		return nil
	 }
	return err
} 

type testService struct{
	*Service
}

func newTestService() *testService{
	return &testService{Service: &Service{}}
}
// Функция добавления аккаунта и баланса в нем
func (s * testService) AddAccountWithBalance(phone types.Phone, balance types.Money) (*types.Account, error){
	// регистрируем пользователя
	account, err := s.RegisterAccount(phone)
	if(err != nil){
		return nil, fmt.Errorf("can't register account, error = %v", err)
	}

	//пополняем его счет
	err = s.Deposit(account.ID, balance)
	if(err != nil){
		return nil, fmt.Errorf("can't deposit account, error = %v", err)
	}

	return account, nil
}