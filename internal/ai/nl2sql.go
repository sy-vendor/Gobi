package ai

import (
	"fmt"
)

var NL2SQL = func(question, schema string) (string, error) {
	prompt := fmt.Sprintf(
		`你是一个资深SQL专家，善于将自然语言问题转为高效、标准的SQL查询。
请严格按照以下数据库表结构和字段类型生成SQL，只返回SQL语句本身，不要任何解释或注释。
表结构和示例数据：
%s

用户问题：
%s

请直接输出SQL语句：`, schema, question)
	return CallDeepSeek(prompt, true)
}
