package handlers

import (
	"gobi/internal/ai"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NL2SQLHandler(c *gin.Context) {
	var req struct {
		Question string `json:"question"`
		Schema   string `json:"schema"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sql, err := ai.NL2SQL(req.Question, req.Schema)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"sql": sql})
}

func InsightHandler(c *gin.Context) {
	var req struct {
		Table   string `json:"table"`
		Metrics string `json:"metrics"`
		Summary string `json:"summary"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	insight, err := ai.SmartInsight(req.Table, req.Metrics, req.Summary)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"insight": insight})
}

func ReportGenHandler(c *gin.Context) {
	var req struct {
		Requirement string `json:"requirement"`
		Schema      string `json:"schema"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	report, err := ai.SmartReportGen(req.Requirement, req.Schema)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"report": report})
}

func NL2ChartHandler(c *gin.Context) {
	var req struct {
		Requirement string `json:"requirement"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	option, err := ai.NL2Chart(req.Requirement)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"option": option})
}
