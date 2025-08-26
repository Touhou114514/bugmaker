package config

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

func split_quotes(str string, quote rune) []string {
	type states int
	const (
		in_space states = iota
		in_field
		in_quote
		in_quote_escape
	)
	state := in_space
	res := []string{}
	var buffer bytes.Buffer

	for _, c := range str {
		switch state {
		case in_space:
			if c == quote {
				state = in_quote
			} else if !unicode.IsSpace(c) {
				buffer.WriteRune(c)
				state = in_field
			}
		case in_field:
			if c == quote {
				state = in_quote
			} else if unicode.IsSpace(c) {
				res = append(res, buffer.String())
				buffer.Reset()
				state = in_space
			} else {
				buffer.WriteRune(c)
			}

		case in_quote:
			switch c {
			case quote:
				state = in_field
			case '\\':
				state = in_quote_escape
			default:
				buffer.WriteRune(c)
			}

		case in_quote_escape:
			buffer.WriteRune(c)
			state = in_quote
		}
	}

	if state == in_field || buffer.Len() != 0 {
		res = append(res, buffer.String())
	}

	return res
}

func args_type(res string, conf_name string, field reflect.Value) error {
	ak := func(tp reflect.Type) (reflect.Value, error) {
		switch tp.Kind() {
		case reflect.Bool:
			if res != "true" && res != "false" {
				return reflect.Value{}, fmt.Errorf("arg[%q] should be true or false", conf_name)
			}
			v := res == "true"
			return reflect.ValueOf(&v), nil
		case reflect.Int:
			num, err := strconv.Atoi(res)
			if err != nil {
				return reflect.Value{}, fmt.Errorf("arg[%q] should be number", conf_name)
			}
			if num < 0 {
				return reflect.Value{}, fmt.Errorf("arg[%q] should be positive number", conf_name)
			}
			return reflect.ValueOf(&num), nil
		case reflect.String:
			un_quote, err := strconv.Unquote(res)
			if err == nil {
				res = un_quote
			}
			return reflect.ValueOf(&res), nil
		case reflect.Interface:
			fallthrough
		default:
			return reflect.Value{}, fmt.Errorf("arg[%q] unsupported type %s", conf_name, tp.Kind().String())
		}
	}
	if field.Kind() == reflect.Ptr {
		v, err := ak(field.Type().Elem())
		if err != nil {
			return err
		}
		field.Set(v)
	} else {
		v, err := ak(field.Type())
		if err != nil {
			return err
		}
		field.Set(v.Elem())
	}
	return nil
}

type config_iter struct {
	conf_type  reflect.Type
	conf_value reflect.Value
	index      int
	tag        string
}

func inter_conf(conf interface{}, tag string) *config_iter {
	v := reflect.ValueOf(conf).Elem()
	tp := v.Type()
	return &config_iter{
		conf_type:  tp,
		conf_value: v,
		index:      -1,
		tag:        tag,
	}
}

func (it *config_iter) next() bool {
	it.index++
	return it.index < it.conf_type.NumField()
}

func (it *config_iter) field() (name string, field reflect.Value) {
	f := it.conf_type.Field(it.index).Tag.Get(it.tag)
	if comma := strings.Index(f, ","); comma >= 0 {
		f = f[:comma]
	}
	field = it.conf_value.Field(it.index)
	return f, field
}

func config_list(w io.Writer, conf interface{}, tag string) {
	it := inter_conf(conf, tag)
	for it.next() {
		f, field := it.field()
		if f == "" {
			continue
		}
		write_field(w, field, f)
	}
}

func list_by_name(conf interface{}, name, tag string) string {
	if name == "" {
		return ""
	}
	it := inter_conf(conf, tag)
	for it.next() {
		f, field := it.field()
		if f == name {
			var buf bytes.Buffer
			write_field(&buf, field, f)
			return buf.String()
		}
	}

	return ""
}

func find_field_by_name(conf interface{}, name, tag string) reflect.Value {
	it := inter_conf(conf, tag)
	for it.next() {
		f, field := it.field()
		if f == name {
			return field
		}
	}
	return reflect.Value{}
}

func write_field(w io.Writer, field reflect.Value, f string) {
	switch field.Kind() {
	case reflect.Interface:
		switch field := field.Interface().(type) {
		case string:
			fmt.Fprintf(w, "%s\t%q\n", f, field)
		default:
			fmt.Fprintf(w, "%s\t%v\n", f, field)
		}
	case reflect.Ptr:
		if !field.IsNil() {
			fmt.Fprintf(w, "%s\t%v\n", f, field.Len())
		} else {
			fmt.Fprintf(w, "%s\t<no define>\n", f)
		}
	case reflect.String:
		fmt.Fprintf(w, "%s\t%q\n", f, field.String())
	default:
		fmt.Fprintf(w, "%s\t%v\n", f, field)
	}
}

// util
func Split2PartsBySpace(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return []string{""}
	}
	parts := strings.Fields(s)
	if len(parts) <= 2 {
		return parts
	}
	// Merge all parts after the first into the second part
	return []string{parts[0], strings.Join(parts[1:], " ")}
}
