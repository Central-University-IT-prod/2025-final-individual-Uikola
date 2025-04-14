package validator

import "github.com/google/uuid"

func ClientID(clientID string) bool {
	_, err := uuid.Parse(clientID)
	return err == nil
}
