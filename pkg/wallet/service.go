package wallet

import (
	"errors"
	"fmt"

	"github.com/ManizhaM/wallet/pkg/types"
	"github.com/google/uuid"
)

var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greater than 0")
var ErrAccountNotFound = errors.New("account not found")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrNotEnoughBalance = errors.New("not enough balance")
var ErrFavoriteNotFound = errors.New("favorite not found")

type Service struct {
	nextAccountID int64
	accounts []*types.Account
	payments []*types.Payment
	favorites []*types.Favorite
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
	paymentID := uuid.New().String()
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

func (s *Service) Repeat(paymentID string)(*types.Payment, error)  {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil{
		return nil, fmt.Errorf("FindPaymentByID(): can't find payment, error = %v", err)
	}
	if payment != nil {
		account, err := s.FindAccountByID(payment.AccountID)
		if account == nil{
			return nil, fmt.Errorf("FindAccountByID(): can't find account, error = %v", err)
		}	
		account.Balance -= payment.Amount
	}
	newPayment := &types.Payment{
		ID: 			uuid.New().String(), 
		AccountID: 		payment.AccountID,
		Amount: 		payment.Amount,
		Category: 		payment.Category,
		Status: 		payment.Status,			
	}
	s.payments = append(s.payments, newPayment)
	return newPayment, nil
}
// FavoritePayment - создаёт избранное из конкретного платежа
func (s *Service) FavoritePayment(paymentID string, name string)(*types.Favorite, error)  {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil{
		return nil, ErrPaymentNotFound
	}
	favorite := &types.Favorite{
		ID: uuid.New().String(),
		AccountID: payment.AccountID,
		Name: name,
		Amount: payment.Amount,
		Category: payment.Category,
	}
	s.favorites = append(s.favorites, favorite)
	return favorite, nil
}

//PayFromFavorite – совершает платёж из конкретного избранного
func (s *Service) PayFromFavorite(favoriteID string)(*types.Payment, error) {
	for _, favorite := range s.favorites{
		if(favorite.ID == favoriteID){
			account, err := s.FindAccountByID(favorite.AccountID)
			if err != nil{
				return nil, ErrAccountNotFound
			}
			payment := &types.Payment{
				ID: uuid.New().String(),
				AccountID: favorite.AccountID,
				Amount: favorite.Amount,
				Category: favorite.Category,
				Status: types.PaymentStatusInProgress,
			}
			account.Balance -= favorite.Amount
			s.payments = append(s.payments, payment)
			return payment, nil
		}
	}
	return nil, ErrFavoriteNotFound
}


