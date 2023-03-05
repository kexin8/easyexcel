package easyexcel

import (
	"io"
	"sync"

	"github.com/kexin8/easyexcel/gopool"
	"github.com/xuri/excelize/v2"
)

type MultipleOption struct {
	*Option

	maxPoolSize int32 // 最大协程池大小,建议设置为当前机器的cpu核数*2
	minPoolSize int32 // 最小协程池大小,建议设置为当前机器的cpu核数
}

// NewMultipleOption 创建Option
func NewMultipleOption() *MultipleOption {
	return &MultipleOption{
		Option:      NewOption(),
		maxPoolSize: 10,
		minPoolSize: 1,
	}
}

// SetMaxPoolSize 设置最大协程池大小
func (o *MultipleOption) SetMaxPoolSize(maxPoolSize int32) *MultipleOption {
	o.maxPoolSize = maxPoolSize
	return o
}

// SetMinPoolSize 设置最小协程池大小
func (o *MultipleOption) SetMinPoolSize(minPoolSize int32) *MultipleOption {
	o.minPoolSize = minPoolSize
	return o
}

// Import 导入
func MultipleImport[T any](file string, opt *MultipleOption, callback func(index int64, data T)) (err error) {
	f, err := excelize.OpenFile(file)
	if err != nil {
		return
	}
	defer f.Close()

	return multipleImport(f, opt, callback)
}

// ImportReader 导入
func MultipleImportReader[T any](file io.Reader, opt *MultipleOption, callback func(index int64, data T)) (err error) {

	f, err := excelize.OpenReader(file)
	if err != nil {
		return
	}
	defer f.Close()

	return multipleImport(f, opt, callback)
}

// multipleImport 多协程导入
// 对于大文件,可以使用多协程来提高效率
func multipleImport[T any](f *excelize.File, opt *MultipleOption, callback func(index int64, data T)) error {

	rows, err := f.Rows(f.GetSheetName(opt.sheetIndex))
	if err != nil {
		return err
	}

	var (
		rowIndex int64 = 0
		fields   Fields

		pool = gopool.NewPool("easyexcel-pool", opt.maxPoolSize, &gopool.Config{ScaleThreshold: opt.minPoolSize})
		wg   sync.WaitGroup
	)

	for rows.Next() {
		if rowIndex < int64(opt.titleIndex) {
			continue
		}

		if rowIndex == int64(opt.titleIndex) {
			titles, err := rows.Columns()
			if err != nil {
				return err
			}

			fields, err = fieldsFromTitls[T](titles)
			if err != nil {
				return err
			}
		}

		var finalRowIndex int64 = 0 // 由于协程执行的时候, rowIndex会变化, 所以需要一个变量来保存当前的rowIndex
		finalRowIndex = rowIndex
		pool.Go(func() {
			wg.Add(1)
			defer wg.Done()

			row, err := rows.Columns()
			if err != nil {
				return
			}

			data, err := rowToStruct[T](row, fields)
			if err != nil {
				return
			}

			callback(finalRowIndex, data)
		})

		rowIndex++
	}

	wg.Wait() // 等待所有协程执行完毕

	return nil
}
