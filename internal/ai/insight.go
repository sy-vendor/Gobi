package ai

import (
	"fmt"
)

var SmartInsight = func(table, metrics, summary string) (string, error) {
	prompt := fmt.Sprintf(
		`你是一个资深商业智能分析师，善于从数据中发现业务洞察。
请根据下方表名、分析指标和数据摘要，自动生成简明的业务分析结论，包括趋势、异常、同比/环比等要点。
表名：%s
分析指标：%s
数据摘要：%s

请用简明中文输出分析结论：`, table, metrics, summary)
	return CallDeepSeek(prompt, true)
}
