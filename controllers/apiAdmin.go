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
	c.SetData(&ai.AINames)
}

func (c *ApiAdminController) CreateStat() {
	count, err := c.GetInt("count")
	if c.SetFail(err) {
		return
	}

	save, err := c.GetBool("save_distribution_in_excel", false)
	if c.SetFail(err) {
		return
	}

	aiTypesJSON := c.GetString("ai_types")
	if c.SetFail(err) {
		return
	}
	var types []int
	err = json.Unmarshal([]byte(aiTypesJSON), &types)
	if c.SetFail(err) {
		return
	}
	go models.NewStat(types, count, save)
	c.SetSuccessResponse()
}

func (c *ApiAdminController) ReadStats() {
	stats := models.ReadStats()
	c.SetData(&stats)
}

func (c *ApiAdminController) DeleteStat() {
	id, err := c.GetInt("id")
	if c.SetFail(err) {
		return
	}

	err = models.DeleteStat(id)
	if c.SetFail(err) {
		return
	}
	c.SetSuccessResponse()
}
