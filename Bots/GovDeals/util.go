package govdeals

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	dynamicconfiguration "github.com/grep-michael/WebBot/DynamicConfiguration"
	"github.com/grep-michael/WebBot/globals"
)

func assetToNotification(asset Asset) globals.Notification {
	return globals.Notification{
		Id:       fmt.Sprintf("%d-%d", asset.AccountID, asset.AssetID),
		Name:     asset.CategoryDescription,
		Message:  asset.ShortDescription,
		Source:   fmt.Sprintf("https://www.govdeals.com/en/asset/%d/%d", asset.AssetID, asset.AccountID),
		ImageUrl: getAssetImageUrl(asset),
		MetaData: globals.NotificationMetaData{
			BotName:   "GovDealsBot",
			Timestamp: time.Now(),
			Tags: map[string]string{
				"Price":          fmt.Sprintf("%0.02f %s", asset.CurrentBid, asset.CurrencyCode),
				"AuctionEndDate": asset.AuctionEndDate.Format("01/02/2006"),
			},
		},
	}
}
func getAssetImageUrl(asset Asset) string {
	return fmt.Sprintf("https://webassets.lqdt1.com/assets/photos/%d/%s", asset.AccountID, asset.Photo)
}
func buildSearchRequest(opts dynamicconfiguration.GovDealsBotOptions) searchRequest {
	return searchRequest{
		CategoryIds:   "",
		BusinessID:    "GD",
		SearchText:    opts.SearchTerm,
		IsQAL:         false,
		AuctionTypeID: nil,
		Page:          1,
		DisplayRows:   opts.DisplayRows,
		SortField:     opts.SortField,
		SortOrder:     opts.SortOrder,
		SessionID:     requestSessionID,
		RequestType:   "search",
		ResponseStyle: "productsOnly",
		Facets: []string{
			"categoryName", "auctionTypeID", "condition", "saleEventName",
			"sellerDisplayName", "product_pricecents", "isReserveMet",
			"hasBuyNowPrice", "isReserveNotMet", "sellerType", "warehouseId",
			"region", "currencyTypeCode", "tierId",
		},
		FacetsFilter:         []any{},
		TimeType:             "",
		SellerTypeID:         nil,
		AccountIDs:           []any{},
		IsSimpleTimeSearch:   true,
		SimpleTimeSearchType: "atauction",
		SimpleTimeWithIn:     0,
		ToDate:               nil,
		FromDate:             nil,
		TimeUnitValue:        "",
		IsVehicleSearch:      false,
	}
}

func decompressBody(resp *http.Response) (io.ReadCloser, error) {
	if resp.Header.Get("Content-Encoding") == "gzip" {
		gr, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		return gr, nil
	}
	return resp.Body, nil
}

func dumpJSON(v any, filename string) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Printf("dumpJSON: marshal: %v", err)
		return
	}
	if err := os.WriteFile(filename, data, 0644); err != nil {
		log.Printf("dumpJSON: write file: %v", err)
	}
}
