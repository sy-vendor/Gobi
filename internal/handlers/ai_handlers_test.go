package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"gobi/config"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	// 切换到项目根目录
	rootDir, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		panic("Failed to get project root: " + err.Error())
	}
	if err := os.Chdir(rootDir); err != nil {
		panic("Failed to change working directory: " + err.Error())
	}
	os.Exit(m.Run())
}

func TestNL2SQLHandler(t *testing.T) {
	err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	router := gin.Default()
	router.POST("/api/ai/nl2sql", NL2SQLHandler)

	reqBody := map[string]string{
		"question": "本月销售额是多少？",
		"schema":   "CREATE TABLE sales (id INT, amount DECIMAL, month VARCHAR(10));",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/ai/nl2sql", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", w.Code)
	}
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["sql"] == "" {
		t.Fatalf("Expected sql in response, got %v", resp)
	}
	t.Logf("Response: %v", resp)
}

func TestInsightHandler(t *testing.T) {
	err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	router := gin.Default()
	router.POST("/api/ai/insight", InsightHandler)

	reqBody := map[string]string{
		"table":   "sales",
		"metrics": "revenue,growth",
		"summary": "2024年1-6月销售数据...",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/ai/insight", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", w.Code)
	}
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["insight"] == "" {
		t.Fatalf("Expected insight in response, got %v", resp)
	}
	t.Logf("Response: %v", resp)
}

func TestReportGenHandler(t *testing.T) {
	err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	router := gin.Default()
	router.POST("/api/ai/reportgen", ReportGenHandler)

	reqBody := map[string]string{
		"requirement": "需要一个销售趋势分析报表",
		"schema":      "CREATE TABLE sales (id INT, amount DECIMAL, month VARCHAR(10));",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/ai/reportgen", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", w.Code)
	}
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["report"] == "" {
		t.Fatalf("Expected report in response, got %v", resp)
	}
	t.Logf("Response: %v", resp)
}

func TestNL2ChartHandler(t *testing.T) {
	err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	router := gin.Default()
	router.POST("/api/ai/nl2chart", NL2ChartHandler)

	reqBody := map[string]string{
		"requirement": "请生成2023年每季度销售额变化的图表",
	}
	body, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/api/ai/nl2chart", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", w.Code)
	}
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["option"] == "" {
		t.Fatalf("Expected option in response, got %v", resp)
	}
	t.Logf("Response: %v", resp)
}
