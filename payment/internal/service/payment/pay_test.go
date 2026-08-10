package payment

func (s *ServiceSuite) TestPayOrderSuccess() {
	id, err := s.payService.PayOrder()
	s.NoError(err)
	s.NotNil(id)
}

func (s *ServiceSuite) TestPayOrderFailure() {
	id, err := s.corruptedPayService.PayOrder()
	s.Error(err)
	s.Equal(err.Error(), "generation failed")
	s.Empty(id)
}
