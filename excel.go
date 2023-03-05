package easyexcel

import (
	"io"

	"github.com/xuri/excelize/v2"
)

type Option struct {
	Sheet string
TitleIndex int
}

// Import 导入
func Import[T any](file io.Reader, sheet string, titleIndex int) (datas []T, err error) {

	f, err := excelize.OpenReader(file)
	if err != nil {
		return
	}
	defer f.Close()

	rows, err := f.GetRows(sheet)
	if err != nil {
		return
	}

	if len(rows) == 0 {
		return
	}

	return
}
