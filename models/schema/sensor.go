package schema

type SensorSwag struct {
	StarTime string `json:"star_time" example:"1594382265"`
	EndTime  string `json:"end_time" example:"1595382265"`
	Point    string `json:"point" example:"0018DE743E31"`
}

type PicSwag struct {
	StarTime string `json:"star_time" example:"1594382265"`
	EndTime  string `json:"end_time" example:"1595382265"`
	Point    string `json:"point" example:"0018DE743E31"`
}

type FarmSwag struct {
	StarTime string `json:"star_time" example:"1594382265"`
	EndTime  string `json:"end_time" example:"9999999999"`
	Point    string `json:"point" example:"teddy"`
}