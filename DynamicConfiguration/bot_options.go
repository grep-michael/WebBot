package dynamicconfiguration

type GovDealsBotOptions struct {
	PollingInterval string `json:"PollingInterval"`
	SearchTerm      string `json:"SearchTerm"`
	SortField       string `json:"SortField"`   //bestfit currentbid auctionclose latestonline location
	SortOrder       string `json:"SortOrder"`   //asc - ascending, desc - descending
	DisplayRows     int    `json:"DisplayRows"` //512 max - number of results returned from api
}

type MTGSecretLairBotOptions struct {
	PollingInterval    string `json:"PollingInterval"`
	CollectionInterval string `json:"CollectionInterval"` //time between request when collecting the product listings
	NotifyInital       bool   `json:"NotifyInital"`       //Will send a notification of all detected inital product listings
	SearchCount        int    `json:"SearchCount"`        //max is 50, high number means less requests, just leave max
}
