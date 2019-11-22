package util

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const (
	defultPage      = 1
	defaultCount    = 10
	defaultMaxCount = 100
	structTagDB     = "db"
	tagIsNull       = "null"
	tagIsNotNUll    = "notnull"
)

//QueryParam : API QueryParam Query Param
type QueryParam struct {
	Page       int
	Count      int
	TableAlias string //table alias
	Fields     string //Fields : for multyple field , using coma delimiter ex : id , name , etc ..
	Sort       string
}

//NewQueryParam : initialize QueryParam from query params
func NewQueryParam(r *http.Request) *QueryParam {
	return defaultQueryParam(r)
}

//SetTableAlias : set table name base on
func (qp *QueryParam) SetTableAlias(t string) {
	qp.TableAlias = t + `.`
}

//GetQueryByDefaultStruct get query field from custom QueryParam.Fields ,or default using Struct Fileds
func (qp *QueryParam) GetQueryByDefaultStruct(i interface{}) (string, error) {
	return getQueryFromStruct(qp, structTagDB, i)
}

// GetQueryFields : get query field from custom QueryParam.Fields ,or default using model
func (qp *QueryParam) GetQueryFields(stringFiels []string) string {
	if len(strings.TrimSpace(qp.Fields)) > 0 {
		return ` SElECT ` + qp.Fields
	}
	return ` SElECT ` + strings.Join(stringFiels, ",")
}

//GetQuerySort : generate sql order syntax base on QueryParam.Sort field , default sort "ASC" ,
func (qp *QueryParam) GetQuerySort() string {
	if len(qp.Sort) > 1 {
		i := 0
		sort := getMapSort(qp.Sort)
		q := ` ORDER BY `
		for k, v := range sort {
			if i > 0 {
				q += ` , `
			}
			q += qp.TableAlias + k + ` ` + v
			i++
		}
		return q
	}
	return ``
}

//GetQueryLimit : generate sql syntax of limit & offside
func (qp *QueryParam) GetQueryLimit() string {
	l := strconv.Itoa(qp.Count)
	o := strconv.Itoa((qp.Page - 1) * (qp.Count))

	return ` LIMIT ` + l + ` OFFSET ` + o
}

func defaultQueryParam(r *http.Request) *QueryParam {
	p, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		p = defultPage
	}
	c, err := strconv.Atoi(r.FormValue("count"))
	if err != nil {
		c = defaultCount
	}
	//limit max
	if c >= 100 {
		c = defaultMaxCount
	}

	f := r.FormValue("fields")
	s := r.FormValue("sort")

	return &QueryParam{
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

func getQueryFromStruct(qp *QueryParam, tag string, i interface{}) (string, error) {
	t := reflect.TypeOf(i)
	qp.SetTableAlias(t.Name())

	q := `SELECT `
	param := strings.Split(qp.Fields, ",")
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tableField := field.Tag.Get(tag)
		if len(param) > 1 {
			for _, v := range param {
				if tableField == v {
					q += qp.TableAlias + tableField + ` ,`
					break
				}
			}
		} else {
			if len(tableField) > 0 {
				q += qp.TableAlias + tableField + ` ,`
			}
		}
	}
	return q[:len(q)-1], nil
}

// func getQueryFromStruct(f *QueryParam, tag string, i interface{}) string {
// 	t := reflect.TypeOf(i)
// 	q := `SELECT `
// 	queryParam := strings.Split(qp.Fields, ",")
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

type (
	//Filter : custom filter joson obj parser
	Filter struct {
		obj map[string]interface{}
	}
)
