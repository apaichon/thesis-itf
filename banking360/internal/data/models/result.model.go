package models
import
(
	"net/http"
)

type ResultModel struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func NewResultModel(resp *http.Response, message string) *ResultModel {
	result := ResultModel{}
	result.Code = resp.StatusCode
	result.Status = resp.Status
	result.Message = message
	return &result
}

func NewSuccessModel(message string) *ResultModel {
	result := ResultModel{}
	result.Code = 200
	result.Status = "Ok"
	result.Message = message
	return &result
}

