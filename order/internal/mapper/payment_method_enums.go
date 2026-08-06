package mapper

import (
	orderv1 "github.com/paincake00/microservices-go/shared/pkg/openapi/order/v1"
	paymentv1 "github.com/paincake00/microservices-go/shared/pkg/proto/payment/v1"
)

func MapPaymentMethodFromGrpc(pm paymentv1.PaymentMethod) orderv1.PaymentMethodEnum {
	switch pm {
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CARD:
		return orderv1.PaymentMethodEnumPAYMENTMETHODCARD
	case paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return orderv1.PaymentMethodEnumPAYMENTMETHODCREDITCARD
	case paymentv1.PaymentMethod_PAYMENT_METHOD_SBP:
		return orderv1.PaymentMethodEnumPAYMENTMETHODSBP
	case paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return orderv1.PaymentMethodEnumPAYMENTMETHODINVESTORMONEY
	default:
		return orderv1.PaymentMethodEnumPAYMENTMETHODUNSPECIFIED
	}
}

func MapPaymentMethodToGrpc(pm orderv1.PaymentMethodEnum) paymentv1.PaymentMethod {
	switch pm {
	case orderv1.PaymentMethodEnumPAYMENTMETHODCARD:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CARD
	case orderv1.PaymentMethodEnumPAYMENTMETHODCREDITCARD:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case orderv1.PaymentMethodEnumPAYMENTMETHODSBP:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_SBP
	case orderv1.PaymentMethodEnumPAYMENTMETHODINVESTORMONEY:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}
