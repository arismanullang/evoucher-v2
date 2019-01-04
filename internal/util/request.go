package util

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const (
	defultPage   = 1
	defaultCount = 10
	structTagDB  = "db"
	tagIsNull    = "null"
	tagIsNotNUll = "notnull"
)

//Filter : API filter Query Param
type Filter struct {
	Page   int
	Count  int
	Fields string //Fields : for multyple field , using coma delimiter ex : id , name , etc ..
	Sort   string
}

//NewFilter : initialize filter from query params
func NewFilter(r *http.Request) *Filter {
	return defaultFilter(r)
}

//GetQueryByDefaultStruct get query field from custom Filter.Fields ,or default using Struct Fileds
func (f *Filter) GetQueryByDefaultStruct(i interface{}) string {
	return getQueryFromStruct(f, structTagDB, i)
}

// GetQueryFields get query field from custom Filter.Fields ,or default using model
func (f *Filter) GetQueryFields(stringFiels []string) string {
	if len(strings.TrimSpace(f.Fields)) > 0 {
		return ` SElECT ` + f.Fields
	}
	return ` SElECT ` + strings.Join(stringFiels, ",")
}

//GetQuerySort : generate sql order syntax base on Filter.Sort field
func (f *Filter) GetQuerySort() string {
	if len(f.Sort) > 1 {
		i := 0
		sort := getMapSort(f.Sort)
		q := ` ORDER BY `
		for k, v := range sort {
			if i > 0 {
				q += ` , `
			}
			q += k + ` ` + v
			i++
		}
		return q
	}
	return ``
}

//GetQueryLimit : generate sql syntax of limit & offside
func (f *Filter) GetQueryLimit() string {
	l := strconv.Itoa(f.Count + 1)
	o := strconv.Itoa((f.Page - 1) * f.Count)

	return ` LIMIT ` + l + ` OFFSET ` + o
}

func defaultFilter(r *http.Request) *Filter {
	p, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		p = defultPage
	}
	c, err := strconv.Atoi(r.FormValue("count"))
	if err != nil {
		c = defaultCount
	}
	f := r.FormValue("fields")
	s := r.FormValue("sort")

	return &Filter{
		Page:   p,
		Count:  c,
		Fields: f,
		Sort:   s,
	}
}

func getMapSort(s string) map[string]string {
	field := strings.Split(s, ",")

	sf := make(map[string]string)
	for _, v := range field {
		// fmt.Println("v :", v)
		sortType := func(str string) string {
			if str[len(str)-1:] == "-" {
				return "desc"
			}
			return "asc"
		}
		sf[strings.TrimSuffix(v, "-")] = sortType(v)
	}
	// fmt.Println("sort by :", sf)
	return sf
}

func getQueryFromStruct(f *Filter, tag string, i interface{}) string {
	t := reflect.TypeOf(i)
	q := `SELECT `
	queryParam := strings.Split(f.Fields, ",")
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tableField := field.Tag.Get(tag)
		if len(queryParam) > 1 {
			for _, v := range queryParam {
				if tableField == v {
					q += tableField + `,`
					break
				}
			}
		} else {
			q += tableField + `,`
		}

	}
	return q[:len(q)-1]
}

// func getQueryFromStruct(f *Filter, tag string, i interface{}) string {
// 	t := reflect.TypeOf(i)
// 	q := `SELECT `
// 	queryParam := strings.Split(f.Fields, ",")
// 	for i := 0; i < t.NumField(); i++ {
// 		field := t.Field(i)
// 		tag := strings.Split(field.Tag.Get(tag), ",")
// 		tableField := tag[0]

// 		if len(queryParam) > 1 {
// 			//using query param field
// 			for _, v := range queryParam {
// 				if tableField == v {
// 					if len(tag) > 1 && tag[1] == tagIsNull {
// 						switch field.Type.Kind() {
// 						case reflect.Ptr:
// 							q += coalesce(tableField, `'0001-01-01'::timestamp`) + `,`
// 						case reflect.String:
// 							q += coalesce(tableField, `' '`) + `,`
// 						}
// 					} else {
// 						q += tableField + `,`
// 					}
// 					break
// 				}
// 			}

// 		} else {
// 			//using default struct field db
// 			if len(tag) > 1 && tag[1] == tagIsNull {
// 				switch field.Type.Kind() {
// 				case reflect.Ptr:
// 					q += coalesce(tableField, `'0001-01-01'::timestamp`) + `,`
// 				case reflect.String:
// 					q += coalesce(tableField, `' '`) + `,`
// 				}
// 			} else {
// 				q += tableField + `,`
// 			}
// 		}

// 	}
// 	return q[:len(q)-1]
// }

func coalesce(arg1, arg2 string) string {
	return `coalesce(` + arg1 + `,` + arg2 + `) AS ` + arg1
}
