package govdeals

import (
	"fmt"
	"strings"
	"time"
)

type searchResponse struct {
	Assets             []Asset `json:"assetSearchResults"`
	AssetSearchFacets  []any   `json:"assetSearchFacets"`
	IsAPIFailureActive bool    `json:"isAPIFailureActive"`
	EasMsg             string  `json:"easMsg"`
}

type Asset struct {
	AccountID               int      `json:"accountId"`
	AssetID                 int      `json:"assetId"`
	AuctionID               int      `json:"auctionId"`
	InventoryID             int      `json:"inventoryId"`
	ShortDescription        string   `json:"assetShortDescription"`
	LongDescription         string   `json:"assetLongDescription"`
	Make                    string   `json:"makebrand"`
	Model                   string   `json:"model"`
	ModelYear               string   `json:"modelYear"`
	AssetCategory           string   `json:"assetCategory"`
	CategoryDescription     string   `json:"categoryDescription"`
	LocationID              int      `json:"locationId"`
	LocationAddr1           string   `json:"locationAddress1"`
	LocationAddr2           string   `json:"locationAddress2"`
	LocationCity            string   `json:"locationCity"`
	LocationState           string   `json:"locationState"`
	LocationZip             string   `json:"locationZip"`
	Country                 string   `json:"country"`
	StateDescription        string   `json:"stateDescription"`
	Latitude                float64  `json:"latitude"`
	Longitude               float64  `json:"longitude"`
	AuctionStartDate        JSONTime `json:"assetAuctionStartDate"`
	AuctionEndDate          JSONTime `json:"assetAuctionEndDate"`
	AuctionEndDateUTC       string   `json:"assetAuctionEndDateUtc"`
	AuctionStartDateDisplay string   `json:"assetAuctionStartDateDisplay"`
	AuctionEndDateDisplay   string   `json:"assetAuctionEndDateDisplay"`
	TimeRemaining           string   `json:"timeRemaining"` // format: "D:H:M:S"
	AssetBidPrice           float64  `json:"assetBidPrice"`
	CurrentBid              float64  `json:"currentBid"`
	BidIncrement            float64  `json:"assetBidIncrement"`
	StrikePrice             float64  `json:"assetStrikePrice"`
	HighBidder              string   `json:"highBidder"`
	BidCount                int      `json:"bidCount"`
	BidWatchID              int      `json:"bidWatchId"`
	CurrencyCode            string   `json:"currencyCode"`
	HasReservePrice         bool     `json:"hasReservePrice"`
	IsReserveNotMet         bool     `json:"isReserveNotMet"`
	IsReserveReduced        bool     `json:"isReserveReduced"`
	WillNextBidMeetReserve  bool     `json:"willNextBidMeetReserve"`
	IsFreeAsset             bool     `json:"isFreeAsset"`
	IsNewAsset              bool     `json:"isNewAsset"`
	IsSoldAuction           bool     `json:"isSoldAuction"`
	AllowProbationBidders   bool     `json:"allowProbationBidders"`
	DenyProbBids            bool     `json:"denyProbBids"`
	BusinessID              string   `json:"businessId"`
	GroupID                 string   `json:"groupId"`
	CompanyName             string   `json:"companyName"`
	DisplaySellerName       string   `json:"displaySellerName"`
	DisplayEventID          string   `json:"displayEventId"`
	EventID                 int      `json:"eventId"`
	AuctionTypeID           int      `json:"auctionTypeId"`
	LotNumber               int      `json:"lotNumber"`
	Photo                   string   `json:"photo"`
	ClickURL                string   `json:"clickUrl"`
	Keywords                string   `json:"keywords"`
	CommDesc                string   `json:"commDesc"`
	TermsAndConditions      string   `json:"termsAndConditions"`
	AssetRestrictionCode    string   `json:"assetRestrictionCode"`
	ProximityDistance       float64  `json:"proximityDistance"`
	WarehouseID             int      `json:"warehouseId"`
	CategoryRoutepath       string   `json:"categoryRoutepath"`
}

type JSONTime struct {
	time.Time
}

func (t *JSONTime) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	parsed, err := time.Parse("2006-01-02T15:04:05", str)
	if err != nil {
		return fmt.Errorf("JSONTime: parse %q: %w", str, err)
	}
	t.Time = parsed
	return nil
}

type searchRequest struct {
	CategoryIds          string   `json:"categoryIds"`
	BusinessID           string   `json:"businessId"`
	SearchText           string   `json:"searchText"`
	IsQAL                bool     `json:"isQAL"`
	AuctionTypeID        *int     `json:"auctionTypeId"`
	Page                 int      `json:"page"`
	DisplayRows          int      `json:"displayRows"`
	SortField            string   `json:"sortField"`
	SortOrder            string   `json:"sortOrder"`
	SessionID            string   `json:"sessionId"`
	RequestType          string   `json:"requestType"`
	ResponseStyle        string   `json:"responseStyle"`
	Facets               []string `json:"facets"`
	FacetsFilter         []any    `json:"facetsFilter"`
	TimeType             string   `json:"timeType"`
	SellerTypeID         *int     `json:"sellerTypeId"`
	AccountIDs           []any    `json:"accountIds"`
	IsSimpleTimeSearch   bool     `json:"isSimpleTimeSearch"`
	SimpleTimeSearchType string   `json:"simpleTimeSearchType"`
	SimpleTimeWithIn     int      `json:"simpleTimeWithIn"`
	ToDate               *string  `json:"toDate"`
	FromDate             *string  `json:"fromDate"`
	TimeUnitValue        string   `json:"timeUnitValue"`
	IsVehicleSearch      bool     `json:"isVehicleSearch"`
}
