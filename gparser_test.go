package gparser

import (
	"reflect"
	"testing"
)

func TestGoParser_Match(t *testing.T) {
	tests := []struct {
		name string
		expr string
		data map[string]interface{}
		want bool
	}{
		{
			name: "test_case1",
			expr: "a == 1 && b == 2",
			data: map[string]interface{}{
				"a": 1,
				"b": 2,
			},
			want: true,
		},
		{
			name: "test_case2",
			expr: "a == 1 && b == 2",
			data: map[string]interface{}{
				"a": 1,
				"b": 3,
			},
			want: false,
		},
		{
			name: "test_case3",
			expr: "a == 1 && b == 2 || c == \"test\"",
			data: map[string]interface{}{
				"a": 1,
				"b": 3,
				"c": "test",
			},
			want: true,
		},
		{
			name: "test_case4",
			expr: "a == 1 && b == 2 && c == \"test\"",
			data: map[string]interface{}{
				"a": 1,
				"b": 3,
				"c": "test",
			},
			want: false,
		},
		{
			name: "test_case5",
			expr: "a == 1 && b == 2 && c == \"test\" && d == true",
			data: map[string]interface{}{
				"a": 1,
				"b": 2,
				"c": "test",
				"d": true,
			},
			want: true,
		},
		{
			name: "test_case6",
			expr: "a == 1 && b == 2 && c == \"test\" && d == false",
			data: map[string]interface{}{
				"a": 1,
				"b": 2,
				"c": "test",
				"d": true,
			},
			want: false,
		},
		{
			name: "test_case6",
			expr: "!(a == 1 && b == 2 && c == \"test\" && d == false)",
			data: map[string]interface{}{
				"a": 1,
				"b": 2,
				"c": "test",
				"d": true,
			},
			want: true,
		},
		{
			name: "test_case7",
			expr: "!(a == 1 && b == 2) || (c == \"test\" && d == false)",
			data: map[string]interface{}{
				"a": 1,
				"b": 2,
				"c": "test",
				"d": false,
			},
			want: true,
		},
		{
			name: "test_case8",
			expr: "a >= \"0.0.1\" && c <= \"1.0.0\" && b < \"1.0.0\"",
			data: map[string]interface{}{
				"a": "0.0.2",
				"b": "0.9.9",
				"c": "1.0.0",
			},
			want: true,
		},
		{
			name: " test_case9",
			expr: `in_array(a, []string{"12131","0989988"})`,
			data: map[string]interface{}{
				"a": "12131",
			},
			want: true,
		},
		{
			name: "tset_case10",
			expr: "start_with(a,\"111111111/222222222\")",
			data: map[string]interface{}{
				"a": "111111111/222222222/333333",
			},
			want: true,
		},
		{
			// 测试部门匹配的case
			name: "test_case11",
			expr: "in_organization(a,\"111111111/222222222\")",
			data: map[string]interface{}{
				"a": "111111111/222222222/333333",
			},
			want: true,
		},
		{
			// 测试部门不匹配的case
			name: "test_case13",
			expr: "in_organization(a,\"111111111/44444444\")",
			data: map[string]interface{}{
				"a": "111111111/222222222/333333",
			},
			want: false,
		},
		{
			// 测试data中的变量不覆盖规则中的变量
			name: "test_case14",
			expr: " a == \"test\" && b == \"test02\"",
			data: map[string]interface{}{
				"a": "test",
			},
			want: false,
		},
		{
			// 测试data中的变量不覆盖规则中的变量
			name: "test_case15",
			expr: " a == \"test\" || b == \"test02\"",
			data: map[string]interface{}{
				"a": "test",
			},
			want: true,
		},
		{
			name: "test_case_16",
			expr: " contain_organization(a,[]string{\"111111/222222/333333\",\"1111111/222222/444444\"})",
			data: map[string]interface{}{
				"a": "333333",
			},
			want: false,
		},
		{
			name: "test_case_16",
			expr: "contain_organization(a,[]string{\"111111/222222\"})",
			data: map[string]interface{}{
				"a": "111111/222222/333333",
			},
			want: true,
		},
		{
			name: "test_case_17",
			expr: "osVersion > \"10\"",
			data: map[string]interface{}{
				"osVersion": "7",
			},
			want: false,
		},
		{
			name: "test_case_18",
			expr: "compare_version(version,\"1.0\",\">\")",
			data: map[string]interface{}{
				"version": "1.1",
			},
			want: true,
		},
		{
			name: "test_case_18",
			expr: "compare_version(version,\"1.0.0\",\">=\")",
			data: map[string]interface{}{
				"version": "0.9.1",
			},
			want: false,
		},
		{
			name: "test_case_18",
			expr: "compare_version(version,\"1.0\",\"<=\") && a == 1",
			data: map[string]interface{}{
				"version": "0.9.1",
				"a":       1,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := Match(tt.expr, tt.data); !reflect.DeepEqual(got, tt.want) || err != nil {
				t.Errorf("goParser match failed, want=%v, got=%v, err=%v", tt.want, got, err)
			}
		})
	}
}

func BenchmarkGoParser_Match(b *testing.B) {
	// 规则表达式
	expr := `(a == 1 && b == "b" && in_array(c, []int{100,99,98,97})) || (d == false)`
	// 映射数据
	data := map[string]interface{}{
		"a": 1,
		"b": "b",
		"c": 100,
		"d": true,
	}
	for i := 0; i < b.N; i++ {
		Match(expr, data)
	}
}
