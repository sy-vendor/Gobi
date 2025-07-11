package main

import (
	"fmt"
	"os"
	"time"

	"github.com/xuri/excelize/v2"
)

func main() {
	// 创建新的Excel文件
	f := excelize.NewFile()
	defer f.Close()

	// 设置默认工作表名称
	f.SetSheetName("Sheet1", "销售柱状图日报")

	// 添加标题
	f.SetCellValue("销售柱状图日报", "A1", "销售柱状图日报")
	f.SetCellValue("销售柱状图日报", "A2", fmt.Sprintf("生成时间: %s", time.Now().Format("2006-01-02 15:04:05")))

	// 设置标题样式
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Size:  16,
			Color: "1F4E79",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	f.SetCellStyle("销售柱状图日报", "A1", "A1", titleStyle)
	f.MergeCell("销售柱状图日报", "A1", "C1")

	// 设置时间样式
	timeStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:  10,
			Color: "666666",
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})
	f.SetCellStyle("销售柱状图日报", "A2", "A2", timeStyle)
	f.MergeCell("销售柱状图日报", "A2", "C2")

	// 添加数据表头
	headers := []string{"月份", "产品类别", "销售额"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c4", 'A'+i)
		f.SetCellValue("销售柱状图日报", cell, header)
	}

	// 设置表头样式
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold:  true,
			Color: "FFFFFF",
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	f.SetCellStyle("销售柱状图日报", "A4", "C4", headerStyle)

	// 添加示例数据
	sampleData := [][]interface{}{
		{"2024-01", "电子产品", 150000},
		{"2024-01", "服装", 80000},
		{"2024-01", "食品", 120000},
		{"2024-02", "电子产品", 180000},
		{"2024-02", "服装", 95000},
		{"2024-02", "食品", 135000},
		{"2024-03", "电子产品", 220000},
		{"2024-03", "服装", 110000},
		{"2024-03", "食品", 150000},
	}

	// 填充数据
	for i, row := range sampleData {
		rowNum := i + 5
		for j, value := range row {
			cell := fmt.Sprintf("%c%d", 'A'+j, rowNum)
			f.SetCellValue("销售柱状图日报", cell, value)
		}
	}

	// 设置数据行样式
	dataStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})
	f.SetCellStyle("销售柱状图日报", "A5", fmt.Sprintf("C%d", 5+len(sampleData)-1), dataStyle)

	// 设置数字格式（销售额列）
	numberStyle, _ := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		NumFmt: 44, // 货币格式
	})
	f.SetCellStyle("销售柱状图日报", "C5", fmt.Sprintf("C%d", 5+len(sampleData)-1), numberStyle)

	// 添加汇总行
	summaryRow := len(sampleData) + 5
	f.SetCellValue("销售柱状图日报", fmt.Sprintf("A%d", summaryRow), "总计")
	f.SetCellValue("销售柱状图日报", fmt.Sprintf("C%d", summaryRow), "=SUM(C5:C13)")

	// 设置汇总行样式
	summaryStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"E7E6E6"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		NumFmt: 44,
	})
	f.SetCellStyle("销售柱状图日报", fmt.Sprintf("A%d", summaryRow), fmt.Sprintf("C%d", summaryRow), summaryStyle)

	// 添加图表说明
	chartDescRow := summaryRow + 2
	f.SetCellValue("销售柱状图日报", fmt.Sprintf("A%d", chartDescRow), "图表说明:")
	f.SetCellValue("销售柱状图日报", fmt.Sprintf("A%d", chartDescRow+1), "• 此模板用于生成按月份和产品类别分组的销售柱状图")
	f.SetCellValue("销售柱状图日报", fmt.Sprintf("A%d", chartDescRow+2), "• 数据来源SQL: SELECT month AS 月份, category AS 产品类别, SUM(amount) AS 销售额 FROM sales GROUP BY month, category ORDER BY month, category")
	f.SetCellValue("销售柱状图日报", fmt.Sprintf("A%d", chartDescRow+3), "• 建议图表类型: 分组柱状图，X轴为月份，Y轴为销售额，不同颜色代表不同产品类别")

	// 设置说明文字样式
	descStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Size:  10,
			Color: "666666",
		},
	})
	f.SetCellStyle("销售柱状图日报", fmt.Sprintf("A%d", chartDescRow), fmt.Sprintf("A%d", chartDescRow+3), descStyle)

	// 添加数据验证说明
	validationRow := chartDescRow + 5
	f.SetCellValue("销售柱状图日报", fmt.Sprintf("A%d", validationRow), "数据验证:")
	f.SetCellValue("销售柱状图日报", fmt.Sprintf("A%d", validationRow+1), "• 月份格式: YYYY-MM")
	f.SetCellValue("销售柱状图日报", fmt.Sprintf("A%d", validationRow+2), "• 销售额: 必须为正数")
	f.SetCellValue("销售柱状图日报", fmt.Sprintf("A%d", validationRow+3), "• 产品类别: 不能为空")

	f.SetCellStyle("销售柱状图日报", fmt.Sprintf("A%d", validationRow), fmt.Sprintf("A%d", validationRow+3), descStyle)

	// 设置列宽
	f.SetColWidth("销售柱状图日报", "A", "A", 15)
	f.SetColWidth("销售柱状图日报", "B", "B", 20)
	f.SetColWidth("销售柱状图日报", "C", "C", 15)

	// 设置行高
	f.SetRowHeight("销售柱状图日报", 1, 30)
	f.SetRowHeight("销售柱状图日报", 2, 20)
	f.SetRowHeight("销售柱状图日报", 4, 25)

	// 保存文件
	filename := fmt.Sprintf("sales_chart_daily_report_template_%s.xlsx", time.Now().Format("20060102_150405"))
	if err := f.SaveAs(filename); err != nil {
		fmt.Printf("保存文件失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("销售柱状图日报模板已生成: %s\n", filename)
	fmt.Println("模板包含以下内容:")
	fmt.Println("1. 销售数据表格（月份、产品类别、销售额）")
	fmt.Println("2. 数据汇总行")
	fmt.Println("3. 图表说明和SQL查询语句")
	fmt.Println("4. 数据验证规则")
	fmt.Println("5. 适合生成分组柱状图的格式化数据")
}
