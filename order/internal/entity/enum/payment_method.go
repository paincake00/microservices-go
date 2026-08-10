package enum

type PaymentMethod string

const (
	PaymentMethodUnspecified   PaymentMethod = "PAYMENT_METHOD_UNSPECIFIED"
	PaymentMethodCard          PaymentMethod = "PAYMENT_METHOD_CARD"
	PaymentMethodSbp           PaymentMethod = "PAYMENT_METHOD_SBP"
	PaymentMethodCreditCard    PaymentMethod = "PAYMENT_METHOD_CREDIT_CARD"
	PaymentMethodInvestorMoney PaymentMethod = "PAYMENT_METHOD_INVESTOR_MONEY"
)
