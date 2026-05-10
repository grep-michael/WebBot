package mtgsecretlair

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type SearchResponse struct {
	Filters []Filter `json:"filters"`
}

type Filter struct {
	Count    int       `json:"count"`
	Total    int       `json:"total"`
	Products []Product `json:"products"`
}

type Product struct {
	ProductID       string        `json:"productID"`
	State           string        `json:"state"`
	ReleaseDate     string        `json:"release_date"`
	Type            string        `json:"type"`
	Refs            Refs          `json:"refs"`
	Country         string        `json:"country"`
	PriceInfo       PriceInfo     `json:"price_info"`
	PriceHide       any           `json:"price_hide"` // int or string in the data
	HasLowestPrice  string        `json:"has_lowest_price"`
	MetaRefIDLinked string        `json:"meta_refid_linked"`
	Image           string        `json:"image"`
	SecondaryImage  string        `json:"secondary_image"`
	Sample          []string      `json:"sample"`
	Note            string        `json:"note"`
	Descriptions    []Description `json:"descriptions"`
	Stock           Stock         `json:"stock"`
	IsReserved      bool          `json:"is_reserved"`
	Categories      []Category    `json:"categories"`
	//Specific          Specific          `json:"specific"`
	Publisher         string            `json:"publisher"`
	Preorder          string            `json:"preorder"`
	PreorderDisplay   string            `json:"preorder_display"`
	WaitingList       bool              `json:"waiting_list"`
	Feeds             []Feed            `json:"feeds"`
	RawData           RawData           `json:"raw_data"`
	Upselling         Upselling         `json:"upselling"`
	FlashSales        FlashSales        `json:"flash_sales"`
	LimitPurchase     FlexInt           `json:"limit_purchase"`
	Backorder         Backorder         `json:"backorder"`
	Protected         bool              `json:"protected"`
	ImageAlt          ImageAlt          `json:"image_alt"`
	SecondaryImageAlt []ImageAlt        `json:"secondary_image_alt"`
	PreviewsPrimary   []PreviewsPrimary `json:"previews_primary"`
	Metaproduct       []any             `json:"metaproduct"`
	Prices            []Price           `json:"prices"`
	CustomFields      CustomFields      `json:"custom_fields"`
	Score             *float64          `json:"_score"`
}

type Refs struct {
	RefID      string `json:"refID"`
	SKU        string `json:"sku"`
	EAN        string `json:"ean"`
	UPC        string `json:"upc"`
	ISBN       string `json:"isbn"`
	JAN        string `json:"jan"`
	ShippingID string `json:"shipping_id"`
	GeoID      string `json:"geoID"`
}

type PriceInfo struct {
	CurrenciesList []string                   `json:"currencies_list"`
	Currencies     map[string][]CurrencyPrice `json:"currencies"`
	Prices         []LicencePrice             `json:"prices"`
}

type CurrencyPrice struct {
	Price        string `json:"price"`
	CrossedPrice string `json:"crossed_price,omitempty"`
}

type LicencePrice struct {
	Licence string `json:"licence"`
	Price   string `json:"price"`
}

type Description struct {
	Lang      string `json:"lang"`
	Title     string `json:"title"`
	HomeTitle string `json:"home_title"`
}

type Stock struct {
	Stock  any    `json:"stock"` // int or string in the data
	OStock string `json:"ostock"`
}

type Category struct {
	CategoryID   string `json:"categoryID"`
	ExternalID   string `json:"externalID"`
	CategoryName string `json:"categoryName"`
}

type Specific struct {
	ContentType string    `json:"content_type,omitempty"`
	State       string    `json:"state,omitempty"`
	Gender      string    `json:"gender,omitempty"`
	Color       string    `json:"color,omitempty"`
	Size        string    `json:"size,omitempty"`
	Extension   Extension `json:"extension,omitempty"`
	// Bundle-specific fields
	ElemVisible        string          `json:"elem_visible,omitempty"`
	ElemVisibleProduct string          `json:"elem_visible_product,omitempty"`
	DiscountType       string          `json:"discount_type,omitempty"`
	Discount           string          `json:"discount,omitempty"`
	GlobalPreorder     string          `json:"global_preorder,omitempty"`
	IsManualDiscount   string          `json:"is_manual_discount,omitempty"`
	Elements           []BundleElement `json:"elements,omitempty"`
}

type Extension struct {
	Videogame       *VideogameExtension `json:"videogame,omitempty"`
	Custom          ExtensionCustom     `json:"custom"`
	Physical        PhysicalExtension   `json:"physical"`
	Other           ExtensionOther      `json:"other"`
	ExtensionCustom []any               `json:"extensiom_custom,omitempty"` // typo in API
	Groups          map[string]string   `json:"groups,omitempty"`
}

type VideogameExtension struct {
	Edition string `json:"edition"`
}

type ExtensionCustom struct {
	IDToUseWithProvider string `json:"id_to_use_with_provider"`
}

type PhysicalExtension struct {
	Weight          string `json:"weight"`
	WeightUnits     string `json:"weight_units"`
	WeightBuffer    string `json:"weight_buffer,omitempty"`
	WeightPartial   string `json:"weight_partial,omitempty"`
	PackHeight      string `json:"pack_height,omitempty"`
	PackHeightUnits string `json:"pack_height_units,omitempty"`
	PackWidth       string `json:"pack_width,omitempty"`
	PackWidthUnits  string `json:"pack_width_units,omitempty"`
	PackDepth       string `json:"pack_depth,omitempty"`
	PackDepthUnits  string `json:"pack_depth_units,omitempty"`
}

type ExtensionOther struct {
	RelatedFeeds string `json:"related_feeds"`
}

type BundleElement struct {
	ProductID    int           `json:"productID"`
	Quantity     string        `json:"quantity"`
	Refs         Refs          `json:"refs"`
	PriceInfo    PriceInfo     `json:"price_info"`
	Descriptions []Description `json:"descriptions"`
}

type Feed struct {
	FeedID string `json:"feedID"`
	Pos    int    `json:"pos"`
}

type RawData struct {
	CrowdfundingEnabled bool `json:"crowdfunding_enabled"`
}

type Upselling struct {
	Crosssell []any      `json:"crosssell"`
	Upsell    UpsellData `json:"upsell"`
}

type UpsellData struct {
	RelatedFeeds string `json:"related_feeds"`
}

type FlashSales struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type Backorder struct {
	ReleaseDate string `json:"release_date"`
	DelayDebit  int    `json:"delay_debit"`
}

type ImageAlt struct {
	AltName string `json:"alt_name"`
	Name    string `json:"name"`
}

type PreviewsPrimary struct {
	Previews []PreviewGroup `json:"previews"`
}

type PreviewGroup struct {
	Preview []PreviewItem `json:"preview"`
}

type PreviewItem struct {
	AltName string `json:"alt_name"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	URL     string `json:"url"`
	Visible string `json:"visible"`
}

type Price struct {
	Currency                 string `json:"currency"`
	Price                    string `json:"price"`
	CrossedPrice             string `json:"crossed_price,omitempty"`
	CfKwProductType          string `json:"cf_kw_product_type"`
	CfKwEventStartDate       string `json:"cf_kw_event_start_date"`
	CfKwEventEndDate         string `json:"cf_kw_event_end_date"`
	CfKwDisableStrikethrough string `json:"cf_kw_disable_strikethrough"`
	CfKwFakeProduct          string `json:"cf_kw_fake_product,omitempty"`
	CfKwEnablePrequeue       string `json:"cf_kw_enable_prequeue,omitempty"`
	CfKwDonateURL            string `json:"cf_kw_donate_url,omitempty"`
}

type CustomFields struct {
	CfKwProductType          []string `json:"cf_kw_product_type"`
	CfKwEventStartDate       []string `json:"cf_kw_event_start_date"`
	CfKwEventEndDate         []string `json:"cf_kw_event_end_date"`
	CfKwDisableStrikethrough []string `json:"cf_kw_disable_strikethrough"`
	CfKwFakeProduct          []string `json:"cf_kw_fake_product,omitempty"`
	CfKwEnablePrequeue       []string `json:"cf_kw_enable_prequeue,omitempty"`
	CfKwDonateURL            []string `json:"cf_kw_donate_url,omitempty"`
}

type FlexInt int

func (f *FlexInt) UnmarshalJSON(data []byte) error {
	// Try int first
	var i int
	if err := json.Unmarshal(data, &i); err == nil {
		*f = FlexInt(i)
		return nil
	}

	// Fall back to string
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("cannot unmarshal %s into FlexInt", data)
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("cannot convert %q to int: %w", s, err)
	}

	*f = FlexInt(i)
	return nil
}
