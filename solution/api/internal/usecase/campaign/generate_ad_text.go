package campaign

import (
	"api/internal/entity/request"
	"api/internal/entity/response"
	"api/pkg/ai"
	"context"
)

// GenerateText generates ad text based on the provided ad title and advertiser information.
// It uses the AI client to create a compelling advertisement message.
func (uc *Usecase) GenerateText(ctx context.Context, req request.GenerateAdText, advertiserID string) (response.GenerateAdText, error) {
	advertiser, err := uc.advertiserRepository.GetByID(ctx, advertiserID)
	if err != nil {
		return response.GenerateAdText{}, err
	}

	resp, err := uc.aiClient.Call(ctx, ai.CallRequest{
		Prompt:         ai.GenerateAdTextPrompt,
		AdTitle:        req.AdTitle,
		AdvertiserName: advertiser.Name,
		Message:        ai.GenerateAdTextMsg,
		Context:        req.Context,
	})
	if err != nil {
		return response.GenerateAdText{}, err
	}

	return response.GenerateAdText{
		GeneratedText: resp,
	}, nil
}
