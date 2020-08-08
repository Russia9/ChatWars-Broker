package messages

type OfferMessage struct {
	SellerID     string `json:"sellerId"`
	SellerName   string `json:"sellerName"`
	SellerCastle string `json:"sellerCastle"`
	Item         string `json:"item"`
	Quantity     int    `json:"qty"`
	Price        int    `json:"price"`
}
