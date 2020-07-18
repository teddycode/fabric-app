package schema

type TransactionSwag struct {
	Timestamp int    `json:"timestamp" binding:"required"`
	Type      string `json:"type" binding:"required"`
	Hash      string `json:"hash" binding:"required"`
	Point     string `json:"point" binding:"required"`
}

type QueryTransNumSwag struct {
	Type int `json:"type" binding:"required" example:"1"`
}

type QueryPointTxNumSwag struct {
	Point string `json:"type" binding:"required"`
}
