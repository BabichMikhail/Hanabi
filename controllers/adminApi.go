package controllers

import (
	"encoding/json"

	ai "github.com/BabichMikhail/Hanabi/AI"

	"github.com/BabichMikhail/Hanabi/models"
)

type ApiAdminController struct {
	ApiController
}

func (c *ApiAdminController) GetAINames() {
	result := struct {
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
	}{StatusSuccess, ai.AINames}
	c.SetData(&result)
}

func (c *ApiAdminController) CreateStat() {
	count, err := c.GetInt("count")
	if c.SetError(err) {
		return
	}
	aiTypesJSON := c.GetString("ai_types")
	if c.SetError(err) {
		return
	}
	var types []int
	err = json.Unmarshal([]byte(aiTypesJSON), &types)
	if c.SetError(err) {
		return
	}
	go models.NewStat(types, count)
	c.SetSuccessResponse()
}

func (c *ApiAdminController) ReadStats() {
	result := struct {
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
	}{StatusSuccess, models.ReadStats()}
	c.SetData(&result)
}
