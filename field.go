package easyexcel

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

const (
	EXCEL = "excel"

	//tag标签
	TagName       = "name"       //字段名
	TagColumnType = "columnType" //导出类型(number数字、string文本、image图片)
	TagType       = "type"       //字段类型(all导入导出、import导入、export导出)
	TagConvertExp = "convertExp" //读取内容转表达式 (如: 0=男,1=女,2=未知)
	TagHeight     = "height"     //行高
	TagWidth      = "width"      //列宽
	TagDefault    = "default"    //默认值
	TagPrompt     = "prompt"     //提示信息
	TagCombo      = "combo"      //下拉框,以,分隔

	//类型
	ALL    = "all"    //导入导出
	IMPORT = "import" //导入
	EXPORT = "export" //导出
	//导出类型
	NUMBER = "number" //数字
	STRING = "string" //文本
	IMAGE  = "image"  //图片

	//分隔符
	separator = ";"
)

type Fields map[int]*Field

type Field struct {
	name string //字段名
	typ  string //字段类型(all导入导出、import导入、export导出)
	sort int    //排序,默认按照字段顺序排序;0->∞
	tags string //tag标签

	column          string            //excel列名
	columnType      string            //导出类型(number数字、string文本、image图片)
	readConverters  map[string]string //读取excel转换器
	writeConverters map[string]string //写入excel转换器

	height       float64  //行高
	width        float64  //列宽
	defaultValue string   //默认值
	prompt       string   //提示信息
	combo        []string //下拉框
}

func NewField(name, tag string, sort int) Field {
	return Field{
		name:            name,
		sort:            sort,
		tags:            tag,
		typ:             ALL,    //默认导入导出
		columnType:      STRING, //默认导出类型文本
		readConverters:  make(map[string]string),
		writeConverters: make(map[string]string),
		height:          14, //默认行高
		width:           16, //默认列宽
	}
}

// parseTag 解析tag标签
func (f *Field) parseTag() error {
	//e.g. excel:"name:姓名;convertExp:0=男,1=女,2=未知"

	if f.tags == "-" {
		//e.g. excel:"-"
		f.column = f.name
		return nil
	}

	tags := strings.Split(f.tags, separator)
	for _, tag := range tags {
		kv := strings.Split(tag, ":")
		if len(kv) == 1 && len(tags) == 1 {
			//e.g. excel"姓名" => name:姓名
			f.column = kv[0]
			continue
		} else if len(kv) == 2 {
			switch kv[0] {
			case TagName:
				f.column = kv[1]
			case TagType:
				f.typ = kv[1]
			case TagColumnType:
				f.columnType = kv[1]
			case TagConvertExp:
				converters := strings.Split(kv[1], ",")
				for _, converter := range converters {
					kv := strings.Split(converter, "=")
					if len(kv) == 2 {
						f.readConverters[kv[1]] = kv[0]
						f.writeConverters[kv[0]] = kv[1]
					}
				}
			case TagHeight:
				h, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return err
				}
				f.height = h
			case TagWidth:
				w, err := strconv.ParseFloat(kv[1], 64)
				if err != nil {
					return err
				}
				f.width = w
			case TagDefault:
				f.defaultValue = kv[1]
			case TagPrompt:
				f.prompt = kv[1]
			case TagCombo:
				f.combo = strings.Split(kv[1], ",")
			}
		} else {
			return errors.New("tag format error: " + tag)
		}

	}

	return nil
}

func fieldsFromTitls[T any](titles []string) (fs Fields, err error) {
	fs = make(Fields, len(titles))

	fields, err := fieldsFromStruct[T]()
	if err != nil {
		return
	}

	for i, title := range titles {
		for _, field := range fields {
			if field.column == strings.TrimSpace(title) {
				fs[i] = &field
				break
			}
		}
	}

	return
}

// fieldsFromStruct 从结构体中获取字段信息
func fieldsFromStruct[T any]() ([]Field, error) {
	var fields []Field
	t := reflect.TypeOf(new(T))

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Anonymous {
			continue
		}
		fieldName := field.Name
		tag := field.Tag.Get(EXCEL)

		if tag == "" {
			continue
		}

		f := NewField(fieldName, tag, i)
		if err := f.parseTag(); err != nil {
			return nil, err
		}
		fields = append(fields, f)
	}
	return fields, nil
}
