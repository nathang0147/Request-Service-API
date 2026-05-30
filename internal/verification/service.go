package verification

import "context"

type Service struct {
	defaultProvider string
	repository      Repository
	sessionProvider SessionProvider
	resolver        ProviderResolver
}

func NewService(defaultProvider string, repository Repository, sessionProvider SessionProvider) *Service {
	return &Service{
		defaultProvider: defaultProvider,
		repository:      repository,
		sessionProvider: sessionProvider,
	}
}

func NewServiceWithResolver(defaultProvider string, repository Repository, resolver ProviderResolver) *Service {
	return &Service{
		defaultProvider: defaultProvider,
		repository:      repository,
		resolver:        resolver,
	}
}

func (service *Service) CreateRequest(ctx context.Context, input CreateRequestInput) (CreateRequestResponse, error) {
	request, err := service.repository.CreateRequest(ctx, VerificationRequest{
		BusinessRef:  input.BusinessRef,
		CandidateRef: input.CandidateRef,
		Provider:     service.defaultProvider,
		Status:       StatusCreated,
		Verified:     false,
	})
	if err != nil {
		return CreateRequestResponse{}, err
	}

	sessionProvider := service.sessionProvider
	if sessionProvider == nil && service.resolver != nil {
		var err error
		sessionProvider, err = service.resolver.Resolve(ctx, request.Provider)
		if err != nil {
			return CreateRequestResponse{}, err
		}
	}
	if sessionProvider == nil {
		return CreateRequestResponse{}, &Error{
			Code: ErrCodeProviderSessionCreateFailed,
		}
	}

	providerSession, err := sessionProvider.CreateSession(ctx, ProviderSessionInput{
		RequestID:       request.ID,
		BusinessRef:     request.BusinessRef,
		CandidateRef:    request.CandidateRef,
		Provider:        request.Provider,
		CredentialTypes: input.CredentialTypes,
	})
	if err != nil {
		request.Status = StatusFailed
		request.ReasonCode = ErrCodeProviderSessionCreateFailed
		request.Verified = false
		_, _ = service.repository.UpdateRequestStatus(ctx, request)

		return CreateRequestResponse{}, &Error{
			Code: ErrCodeProviderSessionCreateFailed,
			Err:  err,
		}
	}

	session, err := service.repository.CreateSession(ctx, VerificationSession{
		VerificationRequestID: request.ID,
		Provider:              request.Provider,
		ProviderSessionID:     providerSession.ProviderSessionID,
		QRCodeURL:             providerSession.QRCodeURL,
		DeepLink:              providerSession.DeepLink,
		OfferURL:              providerSession.OfferURL,
		ExpiresAt:             providerSession.ExpiresAt,
		RawCreateResponse:     providerSession.RawCreateResponse,
	})
	if err != nil {
		return CreateRequestResponse{}, err
	}

	request.Status = StatusPending
	request.Verified = false
	request.ReasonCode = ""

	updatedRequest, err := service.repository.UpdateRequestStatus(ctx, request)
	if err != nil {
		return CreateRequestResponse{}, err
	}

	return toCreateRequestResponse(updatedRequest, session), nil
}

func (service *Service) GetRequest(ctx context.Context, requestID string) (GetRequestResponse, error) {
	request, err := service.repository.GetRequestByID(ctx, requestID)
	if err != nil {
		return GetRequestResponse{}, err
	}

	session, err := service.repository.GetLatestSessionByRequestID(ctx, requestID)
	if err != nil {
		return toGetRequestResponse(request, nil), nil
	}

	return toGetRequestResponse(request, &session), nil
}
