package mapper

import (
	"github.com/paincake00/microservices-go/order/internal/entity/enum"
	orderv1 "github.com/paincake00/microservices-go/shared/pkg/openapi/order/v1"
	paymentv1 "github.com/paincake00/microservices-go/shared/pkg/proto/payment/v1"
)

// PaymentMethodToOAS : OAS - Open API Specification
func PaymentMethodToOAS(pm enum.PaymentMethod) orderv1.PaymentMethodEnum {
	switch pm {
	case enum.PaymentMethodCard:
		return orderv1.PaymentMethodEnumPAYMENTMETHODCARD
	case enum.PaymentMethodCreditCard:
		return orderv1.PaymentMethodEnumPAYMENTMETHODCREDITCARD
	case enum.PaymentMethodSbp:
		return orderv1.PaymentMethodEnumPAYMENTMETHODSBP
	case enum.PaymentMethodInvestorMoney:
		return orderv1.PaymentMethodEnumPAYMENTMETHODINVESTORMONEY
	default:
		return orderv1.PaymentMethodEnumPAYMENTMETHODUNSPECIFIED
	}
}

func PaymentMethodFromOAS(pm orderv1.PaymentMethodEnum) enum.PaymentMethod {
	switch pm {
	case orderv1.PaymentMethodEnumPAYMENTMETHODCARD:
		return enum.PaymentMethodCard
	case orderv1.PaymentMethodEnumPAYMENTMETHODCREDITCARD:
		return enum.PaymentMethodCreditCard
	case orderv1.PaymentMethodEnumPAYMENTMETHODSBP:
		return enum.PaymentMethodSbp
	case orderv1.PaymentMethodEnumPAYMENTMETHODINVESTORMONEY:
		return enum.PaymentMethodInvestorMoney
	default:
		return enum.PaymentMethodUnspecified
	}
}

func PaymentMethodToPaymentProto(pm enum.PaymentMethod) paymentv1.PaymentMethod {
	switch pm {
	case enum.PaymentMethodCard:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CARD
	case enum.PaymentMethodCreditCard:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case enum.PaymentMethodSbp:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_SBP
	case enum.PaymentMethodInvestorMoney:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}
