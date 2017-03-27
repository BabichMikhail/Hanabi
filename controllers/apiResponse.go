package controllers

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type FailRespone struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type SuccessResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

type ApiController struct {
	BaseController
}

const (
	StatusSuccess = "success"
	StatusFail    = "fail"
)

func (c *ApiController) SetFail(err error) bool {
	if err == nil {
		return false
	}
	result := ErrorResponse{StatusFail, err.Error()}
	c.Data["json"] = &result
	c.ServeJSON()
	return true
}

func (c *ApiController) SetSuccessResponse() {
	c.Data["json"] = &SuccessResponse{StatusSuccess, nil}
	c.ServeJSON()
}

func (c *ApiController) SetData(result interface{}) {
	c.Data["json"] = &result
	c.Data["json"] = struct {
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
	}{StatusSuccess, result}
	c.ServeJSON()
}
