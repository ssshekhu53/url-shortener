package service

import (
	"errors"
	"math/rand"

	"url-shortener/models"
	"url-shortener/store"
	"url-shortener/utils"
)

const (
	defaultTTLS = 120
	maxRetries  = 10
	aliasLength = 6
)

type URL interface {
	Create(req *models.CreateRequest) (string, error)
	Get(alias string) (string, error)
	GetAnalytics(alias string) (*models.Analytics, error)
	GetAllAnalytics() []models.Analytics
	Update(req *models.UpdateRequest, alias string) error
	Delete(alias string) error
}

type url struct {
	store store.URL
}

func New(store store.URL) URL {
	return &url{store: store}
}

func (u *url) Create(req *models.CreateRequest) (string, error) {
	var err error

	if req.CustomAlias == "" {
		req.CustomAlias, err = u.generateAlias(maxRetries)
		if err != nil {
			return "", err
		}
	}

	if req.TTLSeconds == 0 {
		req.TTLSeconds = defaultTTLS
	}

	analytics := u.store.GetAnalytics(req.CustomAlias)

	hasExpired, err := utils.HasExpired(analytics)
	if err != nil {
		return "", err
	}

	if analytics != nil && !hasExpired {
		return "", errors.New("alias already exists")
	}

	alias := u.store.Create(req)

	return alias, nil
}

func (u *url) Get(alias string) (string, error) {
	analytics := u.store.GetAnalytics(alias)
	if analytics == nil {
		return "", errors.New("alias does not exist or has expired")
	}

	hasExpired, err := utils.HasExpired(analytics)
	if err != nil {
		return "", err
	}

	if hasExpired {
		return "", errors.New("alias does not exist or has expired")
	}

	u.store.UpdateAccessAnalytics(alias)

	return analytics.LongURL, nil
}

func (u *url) GetAnalytics(alias string) (*models.Analytics, error) {
	analytics := u.store.GetAnalytics(alias)
	if analytics == nil {
		return nil, errors.New("alias does not exist or has expired")
	}

	hasExpired, err := utils.HasExpired(analytics)
	if err != nil {
		return nil, err
	}

	if hasExpired {
		return nil, errors.New("alias does not exist or has expired")
	}

	return analytics, nil
}

func (u *url) GetAllAnalytics() []models.Analytics {
	return u.store.GetAllAnalytics()
}

func (u *url) Update(req *models.UpdateRequest, alias string) error {
	analytics := u.store.GetAnalytics(alias)
	if analytics == nil {
		return errors.New("alias does not exist or has expired")
	}

	hasExpired, err := utils.HasExpired(analytics)
	if err != nil {
		return err
	}

	if hasExpired {
		return errors.New("alias does not exist or has expired")
	}

	if req.CustomAlias == "" {
		req.CustomAlias = analytics.Alias
	}

	if req.TTLSeconds == 0 {
		req.TTLSeconds = analytics.TTLSeconds
	}

	u.store.Update(req, alias)

	return nil
}

func (u *url) Delete(alias string) error {
	analytics := u.store.GetAnalytics(alias)
	if analytics == nil {
		return errors.New("alias does not exist or has expired")
	}

	hasExpired, err := utils.HasExpired(analytics)
	if err != nil {
		return err
	}

	if hasExpired {
		return errors.New("alias does not exist or has expired")
	}

	u.store.Delete(alias)

	return nil
}

func (u *url) generateAlias(maxRetries int) (string, error) {
	var (
		analytics *models.Analytics
		alias     string
	)

	analytics = &models.Analytics{}
	for analytics != nil {
		alias = u.generateRandomString()

		analytics = u.store.GetAnalytics(alias)

		if maxRetries == 0 {
			return "", errors.New("cannot generate alias")
		}

		maxRetries--
	}

	return alias, nil
}

func (u *url) generateRandomString() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, aliasLength)

	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}

	return string(s)
}
