package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatParam struct {
	JWT     string `json:"jwt"`
	PopugID string `json:"taskID"`
}

type StatResult struct {
	Total   uint64 `json:"total"`
	Average uint64 `json:"average"`
	Count   uint64 `json:"count"`
}

func (s *Service) StatHandle(c *gin.Context) {
	p := StatParam{}
	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	result := StatResult{}

	for _, record := range s.billingStorage.log {
		result.Total += record.debit
		result.Total -= record.credit
		result.Count++

	}

	result.Average = result.Total / result.Count

	c.JSON(http.StatusOK, gin.H{
		"list": result,
	})
}

type ListParam struct {
	JWT     string `json:"jwt"`
	PopugID string `json:"popug_id"`
}

type ListResult struct {
	records []Record
	total   uint64
}

func (s *Service) ListHandle(c *gin.Context) {
	p := ListParam{}
	if err := c.BindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
	}

	result := ListResult{
		records: nil,
		total:   0,
	}

	for _, record := range s.billingStorage.log {
		if record.PopugID == p.PopugID {
			result.records = append(result.records)
			result.total += record.debit
			result.total -= record.credit
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"list": result,
	})
}
