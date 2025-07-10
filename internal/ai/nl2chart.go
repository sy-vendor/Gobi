package ai

import (
	"fmt"
)

// NL2Chart 通过自然语言生成 ECharts option JSON
var NL2Chart = func(requirement string) (string, error) {
	prompt := fmt.Sprintf(
		`你是一个数据可视化专家，擅长将业务需求转化为 ECharts 标准的 option JSON。
请根据下方用户需求，生成一个 ECharts 柱状图的 option JSON，要求：
1. x 轴为 2023 年四个季度（第一季度、第二季度、第三季度、第四季度）
2. y 轴为销售额
3. 图表标题为“2023年每季度销售额变化”
4. series 为柱状图，数据请合理虚构
5. 只返回标准 ECharts option JSON，不要任何解释或注释

用户需求：
%s

请直接输出 ECharts option JSON：`, requirement)
	return CallDeepSeek(prompt, false)
}
