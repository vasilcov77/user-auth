//go:build integration

package test

func (s *Suite) Test_Register() {
	email := "user123@example.com"
	phone := "+79031234567"
	password := "securePassword123"
	company := "ООО «Компания123»"

	id, err := s.client.Register(ctx, email, phone, password, company)
	s.NoError(err)
	s.NotEmpty(id)

	aToken, rToken, err := s.client.Login(ctx, email, password)
	s.NoError(err)
	s.NotEmpty(aToken)
	s.NotEmpty(rToken)

	u, err := s.client.GetUser(ctx, id, aToken)

	s.NoError(err)
	s.NotNil(u)

	s.Equal(email, u.Email)
	s.Equal(phone, u.Phone)
	s.True(u.Active)
	s.False(u.VerifiedPhone)
	s.False(u.VerifiedEmail)
}
