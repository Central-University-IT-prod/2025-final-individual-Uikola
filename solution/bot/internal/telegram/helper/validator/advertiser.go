package validator

import "github.com/google/uuid"

func AdvertiserID(advertiserID string) bool {
	_, err := uuid.Parse(advertiserID)
	return err == nil
}
