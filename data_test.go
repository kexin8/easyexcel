package easyexcel

import (
	"reflect"
	"testing"
)

type User struct {
	Name  string  `excel:"name:姓名"`
	Age   int     `excel:"name:年龄"`
	Sex   int     `excel:"name:性别;convertExp:0=男,1=女,2=未知"`
	Money float64 `excel:"name:金额"`
}

func Test_rowToStruct(t *testing.T) {

	type args struct {
		titles []string
		row    []string
	}
	tests := []struct {
		name     string
		args     args
		wantData User
		wantErr  bool
	}{
		// TODO: Add test cases.
		{name: "first test",
			args: args{
				titles: []string{"姓名", "年龄", "性别"},
				row:    []string{"张三", "18", "男"},
			},
			wantData: User{Name: "张三", Age: 18, Sex: 0},
			wantErr:  false,
		},
		{name: "row sex is nil",
			args: args{
				titles: []string{"姓名", "年龄", "性别"},
				row:    []string{"张三", "18"},
			},
			wantData: User{Name: "张三", Age: 18, Sex: 0},
			wantErr:  false,
		},
		{name: "title 缺少性别",
			args: args{
				titles: []string{"姓名", "年龄"},
				row:    []string{"张三", "18", "男"},
			},
			wantData: User{Name: "张三", Age: 18, Sex: 0},
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields, err := fieldsFromTitls[User](tt.args.titles)
			if (err != nil) != tt.wantErr {
				t.Errorf("rowToStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			u, err := rowToStruct[User](tt.args.row, fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("rowToStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(u, tt.wantData) {
				t.Errorf("rowToStruct() = %v, want %v", u, tt.wantData)
			}
		})
	}
}

func Test_structToRow(t *testing.T) {
	type args struct {
		data   User
		titles []string
	}
	tests := []struct {
		name    string
		args    args
		wantRow []string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "first test",
			args: args{
				data:   User{Name: "张三", Age: 18, Sex: 0, Money: 1.1},
				titles: []string{"姓名", "年龄", "性别", "金额"},
			},
			wantRow: []string{"张三", "18", "男", "1.1"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields, err := fieldsFromTitls[User](tt.args.titles)
			if (err != nil) != tt.wantErr {
				t.Errorf("rowToStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			gotRow, err := structToRow(tt.args.data, fields)
			if (err != nil) != tt.wantErr {
				t.Errorf("structToRow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if !reflect.DeepEqual(gotRow, tt.wantRow) {
				t.Errorf("structToRow() = %v, want %v", gotRow, tt.wantRow)
			}
		})
	}
}
