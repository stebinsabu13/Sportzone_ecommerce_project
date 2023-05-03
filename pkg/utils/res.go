package utils

type Colours string

type Size string

// type Product string

type ResponseProducts struct {
	Image     string `json:"image"`
	ModelName string `json:"modelname"`
	Price     uint   `json:"price"`
	BrandName string `json:"brandname"`
}

type ResponseProductDetails struct {
	Image              string    `json:"image"`
	ModelName          string    `json:"modelname"`
	Price              uint      `json:"price"`
	BrandName          string    `json:"brandname"`
	DiscountPercentage string    `json:"discountpercentage"`
	AvailableColours   []Colours `json:"availablecolours"`
	AvailableSizes     []Size    `json:"availablesizes"`
	PayableAmount      uint      `json:"payableamount"`
}
