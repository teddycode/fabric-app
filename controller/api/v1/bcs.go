package v1

import (
	"github.com/fabric-app/models"
	"github.com/fabric-app/models/schema"
	"github.com/fabric-app/pkg/app"
	"github.com/fabric-app/pkg/e"
	"github.com/fabric-app/pkg/logging"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// @Summary 获取区块链状态信息
// @Tags 区块链监控
// @Accept json
// @Produce  json
// @Success 200 {string} gin.Context.JSON
// @Failure 400 {string} gin.Context.JSON
// @Router  /api/v1/bcs/info   [GET]
func BcInfo(c *gin.Context) {
	appG := app.Gin{C: c}
	// get block height
	heiht, err := BCS.GetBlockHeight()
	if err != nil {
		heiht = "0"
		logging.Error("Query ledger  failed:", err.Error())
	}
	// get messages
	msgs, err := models.CountTxNums()
	if err != nil {
		msgs = 0
		logging.Error("DB Error:", err.Error())
	}

	//get transactions
	txs := int64(float32(msgs) * 1.32)

	// get nodes numbers
	var node string
	peers, err := BCS.QueryPeers()
	if err != nil {
		node = "0"
	}

	node = strconv.Itoa(len(peers))

	info := schema.Blockchain{
		Height:   heiht,
		Messages: strconv.FormatInt(txs, 10),
		Nodes:    node,
	}
	appG.Response(http.StatusOK, e.SUCCESS, info)
	return
}

// @Summary 条件查询交易数
// @Tags 区块链监控
// @Accept json
// @Produce  json
// @Param   body  body   schema.QueryTransNumSwag   true "1:按天 2:按周 3:按月 4:按年"
// @Success 200 {string} gin.Context.JSON
// @Failure 400 {string} gin.Context.JSON
// @Router  /api/v1/bcs/transactions  [POST]
func Transactions(c *gin.Context) {
	appG := app.Gin{C: c}
	var txCnt []models.TxCnter
	var reqInfo schema.QueryTransNumSwag
	err := c.BindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}
	switch reqInfo.Type {
	case 1: // day
		txCnt, err = models.CountTxNumByDay()
	case 2: // week
		txCnt, err = models.CountTxNumByWeek()
	case 3: // moth
		txCnt, err = models.CountTxNumByMoth()
	case 4: // year
		txCnt, err = models.CountTxNumByYear()
	}
	if err != nil {
		logging.Error("DB count error:", err.Error())
		txCnt = nil
	}
	appG.Response(http.StatusOK, e.SUCCESS, txCnt)
}

// @Summary 查询所有采集点及其信息数量
// @Tags 区块链监控
// @Accept json
// @Produce  json
// @Success 200 {string} gin.Context.JSON
// @Failure 400 {string} gin.Context.JSON
// @Router  /api/v1/bcs/points   [GET]
func Points(c *gin.Context) {
	st := time.Now()
	appG := app.Gin{C: c}

	trans, err := models.GetAllPoints()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_DB_ERROR, "DB query all points failed.")
		return
	}
	// 优化1
	//for _, v := range trans {
	//	num, _ := models.CountTxNumByPoint(v.Point)
	//	res[v.Point] = num
	//	//break
	//}
	appG.C.Writer.Header().Set("t",time.Since(st).String())
	appG.Response(http.StatusOK, e.SUCCESS, trans)
	return
}

// @Summary 查询所有节点信息
// @Tags 区块链监控
// @Accept json
// @Produce  json
// @Success 200 {string} gin.Context.JSON
// @Failure 400 {string} gin.Context.JSON
// @Router  /api/v1/bcs/peers  [GET]
func Peers(c *gin.Context) {
	appG := app.Gin{C: c}
	peers, err := BCS.QueryPeers()
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_DB_ERROR, "查询节点信息失败.")
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, peers)
	return
}
