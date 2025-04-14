package banners

import (
	tele "gopkg.in/telebot.v3"
)

type Banner tele.File

func (b *Banner) Caption(caption string) interface{} {
	if b == nil {
		return caption
	}
	return &tele.Photo{File: tele.File{
		FileID:     b.FileID,
		UniqueID:   b.UniqueID,
		FileSize:   b.FileSize,
		FilePath:   b.FilePath,
		FileLocal:  b.FileLocal,
		FileURL:    b.FileURL,
		FileReader: b.FileReader,
	}, Caption: caption}
}

type Banners struct {
	Auth       Banner
	Client     Banner
	Advertiser Banner
}

func New(b *tele.Bot, authBannerID, clientBannerID, advertiserBannerID string) (*Banners, error) {
	authBanner, err := b.FileByID(authBannerID)
	if err != nil {
		return nil, err
	}

	clientBanner, err := b.FileByID(clientBannerID)
	if err != nil {
		return nil, err
	}

	advertiserBanner, err := b.FileByID(advertiserBannerID)
	if err != nil {
		return nil, err
	}

	return &Banners{
		Auth:       Banner(authBanner),
		Client:     Banner(clientBanner),
		Advertiser: Banner(advertiserBanner),
	}, nil
}
