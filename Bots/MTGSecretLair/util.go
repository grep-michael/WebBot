package mtgsecretlair

import (
	"strconv"
	"time"

	"github.com/grep-michael/WebBot/globals"
)

func productToNotification(prod Product) globals.Notification {
	notification := globals.Notification{
		Id:       prod.ProductID,
		Name:     "",
		Source:   "",
		Message:  getProductName("EN", prod),
		ImageUrl: prod.Image,
		MetaData: globals.NotificationMetaData{
			Timestamp: time.Now(),
			Tags: map[string]string{
				"Price":         getPrice("USD", prod),
				"PurchaseLimit": strconv.Itoa(int(prod.LimitPurchase)),
			},
		},
	}

	return notification
}

func getProductName(langCode string, prod Product) string {
	for _, desc := range prod.Descriptions {
		if desc.Lang == langCode {
			return desc.Title
		}
	}
	return ""
}

func getPrice(currency string, prod Product) string {
	for _, price := range prod.Prices {
		if price.Currency == currency {
			return price.Price
		}
	}
	return "-1"
}
