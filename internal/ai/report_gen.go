package ai

import (
	"fmt"
)

var SmartReportGen = func(requirement, schema string) (string, error) {
	prompt := fmt.Sprintf(
		`你是一个BI报表专家，擅长根据业务需求和数据库结构设计高质量的数据分析报表。
请根据下方数据库表结构和用户需求，生成一份报表设计方案，内容包括：
1. 推荐的图表类型（如柱状图、折线图、饼图等）
2. 核心分析指标
3. 主要分析维度
4. 推荐的SQL查询示例

数据库表结构：
%s

用户需求：
%s

请用简明中文输出，结构清晰。`, schema, requirement)
	return CallDeepSeek(prompt, true)
}
