package easyexcel

import (
	"reflect"
	"testing"
)

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
			wantData: User{"张三", 18, 0},
			wantErr:  false,
		},
		{name: "row sex is nil",
			args: args{
				titles: []string{"姓名", "年龄", "性别"},
				row:    []string{"张三", "18"},
			},
			wantData: User{"张三", 18, 0},
			wantErr:  false,
		},
		{name: "title 缺少性别",
			args: args{
				titles: []string{"姓名", "年龄"},
				row:    []string{"张三", "18", "男"},
			},
			wantData: User{"张三", 18, 0},
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

type User struct {
	Name string `excel:"name:姓名"`
	Age  int    `excel:"name:年龄"`
	Sex  int    `excel:"name:性别;convertExp:0=男,1=女,2=未知"`
}
