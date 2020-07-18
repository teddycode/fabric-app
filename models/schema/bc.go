package schema

type Blockchain struct {
	Height       string `json:"height"`       // 区块高度
	Messages     string `json:"messages"`     // 上链信息数量
	transactions string `json:"transactions"` // 交易数量
	Nodes        string `json:"nodes"`        // 节点数量
}

type VerifySwag struct {
	Hash string `json:"hash" example:"b9c52e66c1ebfc826e324a394a106f9dc9550fed4390808b2d8932ff91c92b5a"` // 脸上哈希
}

type CheckPic struct {
	Date  string `json:"date" example:"2020-07-13"`
	Hash  string `json:"hash" example:"d8963f1db11e7ca690f159b7e8173205854d7914aba7b47ede48684175b8bc39"`
	Point string `json:"point" example:"0018DE743E31"`
}

//type UploadSwag struct {
//	Type string `json:"type" example:"p"`
//	Raw  string `json:"raw" example:"{\"point\":\"point001\",\"type\":\"temperature\",\"value\":\"26.2\",\"unit\":\"C\"}"`
//}

type BCSensor struct {
	Point string `json:"point" example:"19372180"`    // 采集点
	Type  string `json:"type"  example:"temperature"` // 传感器参数类型
	Value string `json:"value" example:"26.3"`        // 参数数值
	Unit  string `json:"unit"  example:"℃"`           // 单位
}

type BCPic struct {
	Point string `json:"point"  example:"19372180"`                                                        // 采集点
	Name  string `json:"name"  example:"b9c52e66c1ebfc826e324a394a106f9dc9550fed4390808b2d8932ff91c92b5a"` // 图片文件名（也是脸上哈希）
	Size  string `json:"size" example:"102400"`                                                            // 图片大小
	Type  string `json:"type" example:"sensor"`                                                            // 来源 0: sensor, 1:pic,
}

//type BCFarmSwag struct {
//	Name string `json:"name"`
//
//}
