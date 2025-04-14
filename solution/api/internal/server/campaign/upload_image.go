package campaign

import (
	"api/internal/entity/response"
	"api/internal/errorz"
	"api/internal/server/helper"
	"api/pkg/s3"
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// UploadImage handles the HTTP request to upload an image for a specific campaign.
// It accepts a multipart/form-data request with a file, uploads the file to S3,
// associates the uploaded image with the campaign, and returns the image URL.
func (h *Handler) UploadImage(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	advertiserID := chi.URLParam(r, "advertiserId")
	valid := helper.IsValidUUID(advertiserID)
	if !valid {
		return errorz.APIError{
			Status: http.StatusBadRequest,
			Err:    errorz.ErrInvalidAdvertiserID,
			Msg:    "invalid advertiser id",
		}
	}

	campaignID := chi.URLParam(r, "campaignId")
	valid = helper.IsValidUUID(campaignID)
	if !valid {
		return errorz.APIError{
			Status: http.StatusBadRequest,
			Err:    errorz.ErrInvalidCampaignID,
			Msg:    "invalid campaign id",
		}
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return errorz.APIError{
			Status: http.StatusBadRequest,
			Err:    err,
			Msg:    "failed to parse multipart form",
		}
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		return errorz.APIError{
			Status: http.StatusBadRequest,
			Err:    err,
			Msg:    "failed to receive file from form",
		}
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return errorz.APIError{
			Status: http.StatusBadRequest,
			Err:    err,
			Msg:    "failed to read file",
		}
	}

	fileData := s3.FileDataType{
		FileName: header.Filename,
		Data:     fileBytes,
	}

	link, err := h.campaignUsecase.UploadImage(ctx, fileData, advertiserID, campaignID)
	switch {
	case errors.Is(err, errorz.ErrAdvertiserNotFound):
		return errorz.APIError{
			Status: http.StatusNotFound,
			Err:    err,
			Msg:    "advertiser not found",
		}
	case errors.Is(err, errorz.ErrCampaignNotFound):
		return errorz.APIError{
			Status: http.StatusNotFound,
			Err:    err,
			Msg:    "campaign not found",
		}
	case err != nil:
		return errorz.APIError{
			Status: http.StatusInternalServerError,
			Err:    err,
			Msg:    "failed to upload image",
		}
	}

	resp := response.UploadImage{
		ImageURL: link,
	}

	return helper.WriteJSON(w, http.StatusOK, resp)
}
