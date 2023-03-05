package easyexcel

import (
	"io"

	"github.com/xuri/excelize/v2"
)

type Option struct {
	sheetIndex int //sheet索引
	titleIndex int //标题索引

	ignoreEmptyRow bool
}

// NewOption 创建Option
func NewOption() *Option {
	return &Option{
		sheetIndex:     0,
		titleIndex:     0,
		ignoreEmptyRow: true,
	}
}

// SheetIndex 设置sheet索引,默认0
func (o *Option) SheetIndex(sheetIndex int) *Option {
	o.sheetIndex = sheetIndex
	return o
}

// TitleIndex 设置标题索引,默认0
func (o *Option) TitleIndex(titleIndex int) *Option {
	o.titleIndex = titleIndex
	return o
}

// IgnoreEmptyRow 设置是否忽略空行,默认true
func (o *Option) IgnoreEmptyRow(ignoreEmptyRow bool) *Option {
	o.ignoreEmptyRow = ignoreEmptyRow
	return o
}

// Import 导入
func Import[T any](file string, opt *Option) (datas []T, err error) {
	f, err := excelize.OpenFile(file)
	if err != nil {
		return
	}
	defer f.Close()

	return importing[T](f, opt)
}

// ImportReader 导入
func ImportReader[T any](file io.Reader, opt *Option) (datas []T, err error) {

	f, err := excelize.OpenReader(file)
	if err != nil {
		return
	}
	defer f.Close()

	return importing[T](f, opt)
}

func importing[T any](f *excelize.File, opt *Option) (datas []T, err error) {
	rows, err := f.GetRows(f.GetSheetName(opt.sheetIndex))
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 || len(rows) <= opt.titleIndex {
		return
	}

	titles := rows[opt.titleIndex]

	fields, err := fieldsFromTitls[T](titles)
	if err != nil {
		return nil, err
	}

	for _, row := range rows[opt.titleIndex+1:] {

		if opt.ignoreEmptyRow && isEmptyRow(row) {
			continue
		}

		data, err := rowToStruct[T](row, fields)
		if err != nil {
			return nil, err
		}

		datas = append(datas, data)
	}

	return
}

// isEmptyRow 判断是否是空行
func isEmptyRow(row []string) bool {
	for _, v := range row {
		if v != "" {
			return false
		}
	}
	return true
}
