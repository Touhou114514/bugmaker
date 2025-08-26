package config

import (
	"bytes"
	"reflect"
	"testing"
)

func TestSplitQuotes(t *testing.T) {
	tests := []struct {
		input  string
		quote  rune
		output []string
	}{
		{`a b c`, '"', []string{"a", "b", "c"}},
		{`"a b" c`, '"', []string{"a b", "c"}},
		{`"a \"b\"" c`, '"', []string{`a "b"`, "c"}},
		{`a   "b c" d`, '"', []string{"a", "b c", "d"}},
		{`"a\\b" c`, '"', []string{`a\b`, "c"}},
	}

	for _, tt := range tests {
		got := split_quotes(tt.input, tt.quote)
		if !reflect.DeepEqual(got, tt.output) {
			t.Errorf("split_quotes(%q, %q) = %v, want %v", tt.input, string(tt.quote), got, tt.output)
		}
	}
}

func TestSplit2PartsBySpace(t *testing.T) {
	tests := []struct {
		input  string
		output []string
	}{
		{"a b", []string{"a", "b"}},
		{"  a   b  ", []string{"a", "b"}},
		{"a", []string{"a"}},
		{"", []string{""}},
	}

	for _, tt := range tests {
		got := Split2PartsBySpace(tt.input)
		if !reflect.DeepEqual(got, tt.output) {
			t.Errorf("Split2PartsBySpace(%q) = %v, want %v", tt.input, got, tt.output)
		}
	}
}

func TestArgsType(t *testing.T) {
	var b bool
	var i int
	var s string

	tests := []struct {
		val     string
		field   interface{}
		want    interface{}
		wantErr bool
	}{
		{"true", &b, true, false},
		{"false", &b, false, false},
		{"123", &i, 123, false},
		{"-1", &i, nil, true},
		{`"hello"`, &s, "hello", false},
		{"hello", &s, "hello", false},
	}

	for _, tt := range tests {
		field := reflect.ValueOf(tt.field).Elem()
		err := args_type(tt.val, "test", field)
		if (err != nil) != tt.wantErr {
			t.Errorf("args_type(%q) error = %v, wantErr %v", tt.val, err, tt.wantErr)
			continue
		}
		if !tt.wantErr {
			got := reflect.Indirect(field).Interface()
			if got != tt.want {
				t.Errorf("args_type(%q) = %v, want %v", tt.val, got, tt.want)
			}
		}
	}
}

func TestConfigListAndListByName(t *testing.T) {
	type Conf struct {
		Name  string `tag:"name"`
		Value int    `tag:"value"`
	}
	conf := &Conf{Name: "foo", Value: 42}
	var buf bytes.Buffer
	config_list(&buf, conf, "tag")
	out := buf.String()
	if out == "" {
		t.Error("config_list output is empty")
	}
	val := list_by_name(conf, "name", "tag")
	if val == "" {
		t.Error("list_by_name did not find field")
	}
}
