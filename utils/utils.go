package utils

import (
	"fmt"
	"time"

	"url-shortener/models"
)

func HasExpired(analytics *models.Analytics) (bool, error) {
	if analytics == nil {
		return false, nil
	}

	createdAt, err := time.Parse("2006-01-02 15:04:05 -0700 MST", analytics.CreatedAt)
	if err != nil {
		return false, fmt.Errorf("error while parsing times: %v", err)
	}

	now := time.Now().UTC()

	if uint(now.Sub(createdAt).Seconds()) > analytics.TTLSeconds {
		return true, nil
	}

	return false, nil
}
