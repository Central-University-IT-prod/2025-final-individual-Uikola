package advertiser

import (
	"bot/internal/repository"
	"bot/internal/telegram/banners"
	"bot/internal/telegram/menu"
	"bot/internal/telegram/middlewares"
	"bot/pkg/advertising"
	"strconv"
	"strings"

	"github.com/nlypage/intele"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/layout"
)

type Handler struct {
	layout       *layout.Layout
	inputManager *intele.InputManager
	banners      *banners.Banners

	userRepository     repository.UserRepository
	campaignRepository repository.CampaignRepository

	menuHandler *menu.Handler

	advertisingClient advertising.Client
}

func NewHandler(layout *layout.Layout, inputManager *intele.InputManager, banners *banners.Banners, userRepository repository.UserRepository, campaignRepository repository.CampaignRepository, menuHandler *menu.Handler, advertisingClient advertising.Client) *Handler {
	return &Handler{
		layout:       layout,
		inputManager: inputManager,
		banners:      banners,

		userRepository:     userRepository,
		campaignRepository: campaignRepository,

		menuHandler: menuHandler,

		advertisingClient: advertisingClient,
	}
}

func (h *Handler) AdvertiserSetup(group *tele.Group, middlewaresHandler *middlewares.Handler) {
	group.Use(middlewaresHandler.IsAdvertiser)
	group.Handle(h.layout.Callback("main_menu:create_campaign"), h.createCampaign)
	group.Handle(h.layout.Callback("advertiser:create_campaign:refill"), h.createCampaign)
	group.Handle(h.layout.Callback("advertiser:create_campaign:confirm"), h.confirmCampaignCreation)

	group.Handle(h.layout.Callback("main_menu:my_campaigns"), h.campaignsList)
	group.Handle(h.layout.Callback("advertiser:campaigns:back"), h.campaignsList)
	group.Handle(h.layout.Callback("advertiser:campaigns:prev_page"), h.campaignsList)
	group.Handle(h.layout.Callback("advertiser:campaigns:next_page"), h.campaignsList)

	group.Handle(h.layout.Callback("campaigns:campaign"), h.campaign)
	group.Handle(h.layout.Callback("campaigns:campaign:back"), h.campaign)
	group.Handle(h.layout.Callback("campaign:update"), h.updateCampaign)
	group.Handle(h.layout.Callback("campaign:update:back"), h.updateCampaign)
	group.Handle(h.layout.Callback("campaign:update:impressions_limit"), h.updateImpressionsLimit)
	group.Handle(h.layout.Callback("campaign:update:clicks_limit"), h.updateClicksLimit)
	group.Handle(h.layout.Callback("campaign:update:cost_per_impression"), h.updateCostPerImpression)
	group.Handle(h.layout.Callback("campaign:update:cost_per_click"), h.updateCostPerClick)
	group.Handle(h.layout.Callback("campaign:update:ad_title"), h.updateAdTitle)
	group.Handle(h.layout.Callback("campaign:update:ad_text"), h.updateAdText)
	group.Handle(h.layout.Callback("campaign:update:start_date"), h.updateStartDate)
	group.Handle(h.layout.Callback("campaign:update:end_date"), h.updateEndDate)
	group.Handle(h.layout.Callback("campaign:update:gender"), h.updateGender)
	group.Handle(h.layout.Callback("campaign:update:location"), h.updateLocation)
	group.Handle(h.layout.Callback("campaign:update:age_from"), h.updateAgeFrom)
	group.Handle(h.layout.Callback("campaign:update:age_to"), h.updateAgeTo)

	group.Handle(h.layout.Callback("campaign:delete"), h.deleteCampaign)
	group.Handle(h.layout.Callback("campaign:delete:confirm"), h.confirmCampaignDeletion)
	group.Handle(h.layout.Callback("campaign:delete:cancel"), h.campaign)

	group.Handle(h.layout.Callback("campaign:statistic"), h.campaignStatistic)

	group.Handle(h.layout.Callback("main_menu:advertiser_statistic"), h.advertiserStatistic)
	group.Handle(h.layout.Callback("main_menu:generate_ad_text"), h.generateAdText)
}

func parseCallback(callbackData string) (string, int, error) {
	var (
		campaignID string
		p          int
		err        error
	)

	data := strings.Split(callbackData, " ")
	if len(data) == 2 {
		campaignID = data[0]
		p, err = strconv.Atoi(data[1])
		if err != nil {
			return "", 0, err
		}
	} else if len(data) == 1 {
		p, err = strconv.Atoi(data[0])
		if err != nil {
			return "", 1, nil
		}
	} else {
		return "", 1, nil
	}

	return campaignID, p, nil
}
