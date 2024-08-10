package store

import (
	"sync"
	"time"

	"url-shortener/models"
)

type URL interface {
	Create(req *models.CreateRequest) string
	GetAnalytics(alias string) *models.Analytics
	GetAllAnalytics() []models.Analytics
	Update(req *models.UpdateRequest, alias string)
	UpdateAccessAnalytics(alias string)
	Delete(alias string)
}

type url struct {
	analytics sync.Map
}

func New() URL {
	u := &url{}

	u.analytics = sync.Map{}

	return u
}

func (u *url) Create(req *models.CreateRequest) string {
	analytics := models.Analytics{
		Alias:       req.CustomAlias,
		LongURL:     req.LongURL,
		TTLSeconds:  req.TTLSeconds,
		CreatedAt:   time.Now().UTC().String(),
		AccessTimes: []string{},
	}

	u.analytics.Store(req.CustomAlias, analytics)

	return analytics.Alias
}

func (u *url) GetAnalytics(alias string) *models.Analytics {
	value, ok := u.analytics.Load(alias)
	if !ok {
		return nil
	}

	res := value.(models.Analytics)

	return &res
}

func (u *url) GetAllAnalytics() []models.Analytics {
	analytics := make([]models.Analytics, 0)

	u.analytics.Range(func(key, value any) bool {
		v := value.(models.Analytics)

		analytics = append(analytics, v)

		return true
	})

	return analytics
}

func (u *url) Update(req *models.UpdateRequest, alias string) {
	value, ok := u.analytics.Load(alias)
	if !ok {
		return
	}

	analytics := value.(models.Analytics)

	if req.CustomAlias != "" {
		analytics.Alias = req.CustomAlias
	}

	if req.TTLSeconds != 0 {
		analytics.TTLSeconds = req.TTLSeconds
		analytics.CreatedAt = time.Now().UTC().String()
	}

	analytics.AccessCount = 0
	analytics.AccessTimes = nil

	u.analytics.Delete(alias)

	u.analytics.Store(analytics.Alias, analytics)
}

func (u *url) UpdateAccessAnalytics(alias string) {
	value, ok := u.analytics.Load(alias)
	if !ok {
		return
	}

	analytics := value.(models.Analytics)

	analytics.AccessCount += 1

	u.updateAccessTimes(&analytics, time.Now().UTC().String())

	u.analytics.Store(alias, analytics)
}

func (u *url) Delete(alias string) {
	u.analytics.Delete(alias)
}

func (u *url) updateAccessTimes(analytics *models.Analytics, accessTime string) {
	if len(analytics.AccessTimes) == 10 {
		analytics.AccessTimes = analytics.AccessTimes[0:9]
	}

	analytics.AccessTimes = append([]string{accessTime}, analytics.AccessTimes...)
}
