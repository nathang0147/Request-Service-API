package verification

func toCreateRequestResponse(request VerificationRequest, session VerificationSession) CreateRequestResponse {
	return CreateRequestResponse{
		RequestID:  request.ID,
		Status:     request.Status,
		Verified:   request.Verified,
		ReasonCode: request.ReasonCode,
		Session: SessionDetails{
			QRCodeURL: session.QRCodeURL,
			DeepLink:  session.DeepLink,
			OfferURL:  session.OfferURL,
			ExpiresAt: session.ExpiresAt,
		},
	}
}

func toGetRequestResponse(request VerificationRequest, session *VerificationSession) GetRequestResponse {
	response := GetRequestResponse{
		RequestID:  request.ID,
		Status:     request.Status,
		Verified:   request.Verified,
		ReasonCode: request.ReasonCode,
	}

	if session != nil {
		response.Session = &SessionDetails{
			QRCodeURL: session.QRCodeURL,
			DeepLink:  session.DeepLink,
			OfferURL:  session.OfferURL,
			ExpiresAt: session.ExpiresAt,
		}
	}

	return response
}
