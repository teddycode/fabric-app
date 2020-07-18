package v1

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/fabric-app/models"
	"github.com/fabric-app/models/schema"
	"github.com/fabric-app/pkg/app"
	"github.com/fabric-app/pkg/e"
	"github.com/fabric-app/pkg/setting"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"
)

const (
	DATA_TYPE_SENSOR = 0
	DATA_TYPE_PIC    = 1
	DATA_TYPE_FARM   = 2
)

// @Summary  传感器数据溯源
// @Tags 溯源查询
// @Accept json
// @Produce  json
// @Param   body  body   schema.SensorSwag   true "body"
// @Success 200 {string} gin.Context.JSON
// @Failure 401 {string} gin.Context.JSON
// @Router /api/v1/trace/sensor  [POST]
func Sensors(c *gin.Context) {
	st := time.Now()
	appG := app.Gin{C: c}
	var reqInfo schema.SensorSwag
	err := c.BindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
	}
	res, err := BCS.QueryCC("traceable", "query",
		[]string{"s", reqInfo.Point, reqInfo.StarTime, reqInfo.EndTime}, setting.Peers[0])
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_CC_QUERY_FAILED, "Chaincode query failed.")
		return
	}
	appG.C.Writer.Header().Set("t", time.Since(st).String())
	appG.Response(http.StatusOK, e.SUCCESS, res)
}

// @Summary  图片信息溯源
// @Tags 溯源查询
// @Accept json
// @Produce  json
// @Param   body  body   schema.PicSwag   true "body"
// @Success 200 {string} gin.Context.JSON
// @Failure 401 {string} gin.Context.JSON
// @Router /api/v1/trace/picture  [POST]
func Pictures(c *gin.Context) {
	st := time.Now()
	appG := app.Gin{C: c}
	var reqInfo schema.PicSwag
	err := c.BindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
	}
	res, err := BCS.QueryCC("traceable", "query",
		[]string{"p", reqInfo.Point, reqInfo.StarTime, reqInfo.EndTime}, setting.Peers[0])
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_CC_QUERY_FAILED, "Chaincode query failed.")
		return
	}
	appG.C.Writer.Header().Set("t", time.Since(st).String())
	appG.Response(http.StatusOK, e.SUCCESS, res)
}

// @Summary  农事数据溯源
// @Tags 溯源查询
// @Accept json
// @Produce  json
// @Param  body  body   schema.FarmSwag   true "body"
// @Success 200 {string} gin.Context.JSON
// @Failure 401 {string} gin.Context.JSON
// @Router /api/v1/trace/farmData  [POST]
func Farms(c *gin.Context) {
	st := time.Now()
	appG := app.Gin{C: c}
	var reqInfo schema.FarmSwag
	err := c.BindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, nil)
	}
	res, err := BCS.QueryCC("traceable", "query",
		[]string{"f", reqInfo.Point, reqInfo.StarTime, reqInfo.EndTime}, setting.Peers[0])
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_CC_QUERY_FAILED, "Chaincode query failed.")
		return
	}
	appG.C.Writer.Header().Set("t", time.Since(st).String())
	appG.Response(http.StatusOK, e.SUCCESS, res)
}

// @Summary 图片下载
// @Tags 溯源查询
// @Accept json
// @Produce  json
// @Param   body  body   schema.PictureSwag   true "body"
// @Success 200 {string} gin.Context.JSON
// @Failure 401 {string} gin.Context.JSON
// @Router /api/v1/trace/downloadPic  [POST]
//func DownloadPic(c *gin.Context) {
//	appG := app.Gin{C: c}
//	var reqInfo schema.PictureSwag
//	err := c.BindJSON(&reqInfo)
//	if err != nil {
//		appG.Response(http.StatusOK, e.INVALID_PARAMS, "Invalid paras in json")
//	}
//	file, err := transh.GetPicFile(reqInfo.Point, reqInfo.Date, reqInfo.Name)
//	if err != nil {
//		appG.Response(http.StatusOK, e.ERROR_FILE_GET_FAILED, "Get picture file failed.")
//	}
//	defer file.Close()
//
//	buf := bytes.Buffer{}
//	size, err := buf.ReadFrom(file)
//	if err != nil {
//		appG.Response(http.StatusOK, e.ERROR_FILE_GET_FAILED, "File buffer create failed.")
//		return
//	}
//	logging.Debug("File load success,size:", size)
//
//	appG.C.Writer.Header().Add("Content-Type", "application/octet-stream")
//	appG.C.Writer.Header().Add("Content-Disposition", "attachment;filename="+file.Name())
//	appG.Response(http.StatusOK, e.SUCCESS, buf.Bytes())
//}

// @Summary  链上信息检验
// @Tags 溯源查询
// @Accept json
// @Produce  json
// @Param   body  body   schema.VerifySwag   true "输入交易哈希，返回交易内容（包含文件内容哈希值）"
// @Success 200 {string} gin.Context.JSON
// @Failure 401 {string} gin.Context.JSON
// @Router /api/v1/trace/verify  [POST]
func Verifier(c *gin.Context) {
	st := time.Now()
	appG := app.Gin{C: c}
	var reqInfo schema.VerifySwag
	err := c.BindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, "param bind failed.")
		return
	}
	res, err := BCS.QueryTxByID(reqInfo.Hash, setting.Peers[0])
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_LEDGER_FAILED, "ledger query failed.")
		return
	}
	appG.C.Writer.Header().Set("t", time.Since(st).String())
	appG.Response(http.StatusOK, e.SUCCESS, res)
}

// @Summary  图片哈希校验
// @Tags 溯源查询
// @Accept json
// @Produce  json
// @Param   body  body   schema.CheckPic   true "test"
// @Success 200 {string} gin.Context.JSON
// @Failure 401 {string} gin.Context.JSON
// @Router /api/v1/trace/check  [POST]
func CheckPic(c *gin.Context) {
	st := time.Now()
	appG := app.Gin{C: c}
	var reqInfo schema.CheckPic
	err := c.BindJSON(&reqInfo)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, "param bind failed.")
		return
	}
	var httpClient = &http.Client{}
	url := path.Join( reqInfo.Point, reqInfo.Date, reqInfo.Hash+".jpg")
	httpRequest, _ := http.NewRequest("GET", "http://202.193.60.10/"+url, nil)
	// 发送请求
	resp, err := httpClient.Do(httpRequest)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR, "request picture failed:"+err.Error())
		return
	}
	defer resp.Body.Close()
	response, _ := ioutil.ReadAll(resp.Body)
	picHash := sha256.Sum256(response)
	ph:=fmt.Sprintf("%x", picHash)
	if strings.Compare(ph, reqInfo.Hash) != 0 {
		appG.Response(http.StatusOK, e.ERROR, map[string]interface{}{
			"fileHash":  reqInfo.Hash,
			"checkHash": ph,
			"status":    false,
		})
		return
	}
	appG.C.Writer.Header().Set("t", time.Since(st).String())
	appG.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"fileHash":  reqInfo.Hash,
		"checkHash": ph,
		"status":    true,
	})

}

// @Summary 图片信息上链接口
// @Tags 溯源查询
// @Accept json
// @Produce  json
// @Param   body  body   schema.BCPic   true "返回交易哈希"
// @Security ApiKeyAuth
// @Success 200 {string} gin.Context.JSON
// @Failure 400 {string} gin.Context.JSON
// @Router  /api/v1/trace/upload/pic  [POST]
func UploaderPic(c *gin.Context) {
	appG := app.Gin{C: c}
	var picInfo schema.BCPic
	err := c.BindJSON(&picInfo)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, "Bind json error.")
		return
	}
	raw, _ := json.Marshal(&picInfo)
	//	raw, _ := ioutil.ReadAll(c.Request.Body)
	txID, err := BCS.InvokeCC("traceable", "add",
		[][]byte{[]byte("p"), []byte(picInfo.Point), raw}, setting.Peers)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_CC_INVOKE_FAILED, "Chaincode traceable invoke failed.")
		return
	}
	//local, _ := time.LoadLocation("Local")
	//now, _ := time.ParseInLocation("2006-01-02 15:04:05", "2017-06-20 18:16:15", local)
	models.NewTx(&models.Transaction{
		Timestamp: time.Now(),
		Type:      "p",
		Hash:      string(txID),
		Point:     picInfo.Point,
	})
	appG.Response(http.StatusOK, e.SUCCESS, txID)
}

// @Summary 传感器数据上链接口
// @Tags 溯源查询
// @Accept json
// @Produce  json
// @Param   body  body  schema.BCSensor  true "返回交易哈希"
// @Security ApiKeyAuth
// @Success 200 {string} gin.Context.JSON
// @Failure 400 {string} gin.Context.JSON
// @Router  /api/v1/trace/upload/sensor  [POST]
func UploaderSen(c *gin.Context) {
	appG := app.Gin{C: c}
	var sensor schema.BCSensor
	err := c.BindJSON(&sensor)
	if err != nil {
		appG.Response(http.StatusOK, e.INVALID_PARAMS, "Bind json error.")
		return
	}
	raw, _ := json.Marshal(&sensor)
	//raw, _ := ioutil.ReadAll(c.Request.Body)

	txID, err := BCS.InvokeCC("traceable", "add",
		[][]byte{[]byte("s"), []byte(sensor.Point), raw}, setting.Peers)
	if err != nil {
		appG.Response(http.StatusOK, e.ERROR_CC_INVOKE_FAILED, "Chaincode traceable invoke failed.")
		return
	}
	//local, _ := time.LoadLocation("Local")
	//now, _ := time.ParseInLocation("2006-01-02 15:04:05", "2017-06-20 18:16:15", local)
	models.NewTx(&models.Transaction{
		Timestamp: time.Now(),
		Type:      "s",
		Hash:      string(txID),
		Point:     sensor.Point,
	})
	appG.Response(http.StatusOK, e.SUCCESS, txID)
}
