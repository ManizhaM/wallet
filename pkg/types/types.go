package types

//Money - денежная сумма в мин.еденицах (центы, копейки, дирамы и тд)
type Money int64

//PaymentCategory - категория в которой был совершен платеж(авто, аптеки, рестораны и тд)
type PaymentCategory string

//PaymentStatus - статус платежа
type PaymentStatus string

//Предопределенные статусы платежей
const (
	PaymentStatusOk PaymentStatus = "OK"
	PaymentStatusFail PaymentStatus = "FAIL"
	PaymentStatusInProgress PaymentStatus = "INPROGRESS"
)

//Payment - информация о платеже
type Payment struct {
	ID			string
	AccountID	int64
	Amount 		Money
	Category 	PaymentCategory
	Status		PaymentStatus
}

//Phone - номер телефона
type Phone string

//Account - информация о счете покупателя
type Account struct{
	ID 		int64
	Phone	Phone
	Balance	Money
}

type Favorite struct{
	ID 			string
	AccountID	int64
	Name		string
	Amount 		Money
	Category	PaymentCategory
}