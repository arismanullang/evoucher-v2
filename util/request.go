package util

import (
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx/types"
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
	Page        int
	Count       int
	TableAlias  string //table alias
	Fields      string //Fields : for multyple field , using coma delimiter ex : id , name , etc ..
	Sort        string
	Q           string
	Model       interface{}
	FilterModel interface{}
}

//NewQueryParam : initialize QueryParam from query params
func NewQueryParam(r *http.Request) *QueryParam {
	return defaultQueryParam(r)
}

//NewQueryParamDefault : initialize QueryParam from query params
func NewQueryParamDefault() *QueryParam {
	return defaultQP()
}

//SetTableAlias : set table name base on
func (qp *QueryParam) SetTableAlias(t string) {
	qp.TableAlias = t + `.`
}

//SetFilterModel : Set model which will be used in query
func (qp *QueryParam) SetModel(i interface{}) {
	qp.Model = i
}

//SetFilterModel : Set model of filter which will be used in search query
func (qp *QueryParam) SetFilterModel(i interface{}) {
	qp.FilterModel = i
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
	if qp.Count == -1 {
		return ``
	}

	l := strconv.Itoa(qp.Count + 1)
	o := strconv.Itoa((qp.Page - 1) * (qp.Count))

	return ` LIMIT ` + l + ` OFFSET ` + o
}

func (qp *QueryParam) GetQueryWhereClause(q string, val string) string {
	return q + getQClauseFromStruct(qp, val, qp.Model) + getWhereClauseFromStruct(qp, qp.FilterModel)
}

func (qp *QueryParam) GetQueryWithPagination(q string, sort string, limit string) string {
	return `WITH tbl AS (` + q + sort + `)
			SELECT *
			FROM  (
				TABLE  tbl
				` + limit + `
				) sub
			RIGHT  JOIN (SELECT count(*) FROM tbl) c("count") ON true`
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
	q := r.FormValue("q")

	return &QueryParam{
		Page:   p,
		Count:  c,
		Fields: f,
		Sort:   s,
		Q:      q,
	}
}

func defaultQP() *QueryParam {

	p := defultPage
	c := defaultCount
	f := ""
	s := ""

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
	qp.SetModel(i)
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
			if len(tableField) > 0 && tableField != "count" {
				q += qp.TableAlias + tableField + ` ,`
			}
		}
	}
	return q[:len(q)-1], nil
}

func getQClauseFromStruct(qp *QueryParam, val string, i interface{}) string {
	if len(val) < 1 {
		return ``
	}

	t := reflect.TypeOf(i)
	tv := reflect.ValueOf(i)

	q := ` AND (`
	param := strings.Split(qp.Fields, ",")
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := tv.Field(i)
		f := value.Interface()
		tableField := field.Tag.Get("db")
		switch f.(type) {
		default:
		case string, types.JSONText:
			if len(param) > 1 {
				for _, v := range param {
					if tableField == v {
						q += qp.TableAlias + tableField + `::text ILIKE '%` + val + `%' OR `
						break
					}
				}

			} else {
				if len(tableField) > 0 && tableField != "count" && tableField != "status" {
					q += qp.TableAlias + tableField + `::text ILIKE '%` + val + `%' OR `
				}
			}
		}
	}
	q = q[:len(q)-3]
	q += `) `
	return q
}

func getWhereClauseFromStruct(qp *QueryParam, i interface{}) string {
	if i == nil {
		return ``
	}

	t := reflect.TypeOf(i)
	tv := reflect.ValueOf(i)

	q := ` AND `
	param := strings.Split(qp.Fields, ",")

	if t.NumField() < 1 {
		return ``
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := tv.Field(i)
		tableField := field.Tag.Get("schema")
		tableType := field.Tag.Get("filter")
		if value.String() != "" {
			switch tableType {
			case "string":
				if len(param) > 1 {
					for _, v := range param {
						if tableField == v {
							q += qp.TableAlias + tableField + ` ILIKE '%` + value.String() + `%' `
							break
						}
					}

				} else {
					if len(tableField) > 0 {
						q += qp.TableAlias + tableField + ` ILIKE '%` + value.String() + `%' `
					}
				}
			case "date":
				if len(param) > 1 {
					for _, v := range param {
						if tableField == v {
							dates := strings.Split(value.String(), ",")
							q += ` BETWEEN '` + dates[0] + ` 00:00:00+07'::timestamp AND '` + dates[1] + ` 23:59:59+07'::timestamp `
							break
						}
					}

				} else {
					if len(tableField) > 0 {
						dates := strings.Split(value.String(), ",")
						q += ` BETWEEN '` + dates[0] + ` 00:00:00+07'::timestamp AND '` + dates[1] + ` 23:59:59+07'::timestamp `
						break
					}
				}
			case "array":
				val := arrayToQueryString(strings.Split(value.String(), ","))
				if len(param) > 1 {
					for _, v := range param {
						if tableField == v {
							q += qp.TableAlias + tableField + ` IN (` + val + `) `
							break
						}
					}

				} else {
					if len(tableField) > 0 {
						q += qp.TableAlias + tableField + ` IN (` + val + `) `
					}
				}
			case "json":
				val := arrayToQueryString(strings.Split(value.String(), ","))
				if len(param) > 1 {
					for _, v := range param {
						if tableField == v {
							q += ` json_array_elements (` + qp.TableAlias + tableField + `) @> ARRAY(` + val + `) `
							break
						}
					}

				} else {
					if len(tableField) > 0 {
						q += ` json_array_elements (` + qp.TableAlias + tableField + `) IN (` + val + `) `
					}
				}
			case "json_array":
				fName := `json_array_` + tableField
				q += `EXISTS (
					SELECT 1 FROM json_array_elements(` + qp.TableAlias + tableField + `) ` + fName + `
					WHERE `

				elem := strings.Split(value.String(), ",")
				for _, e := range elem {
					data := strings.Split(e, ":")
					if len(param) > 1 {
						for _, v := range param {
							if tableField == v {
								q += fName + `->>'` + data[0] + `' ILIKE '%` + data[1] + `%' `
								break
							}
						}

					} else {
						if len(tableField) > 0 {
							q += fName + `->>'` + data[0] + `' ILIKE '%` + data[1] + `%' `
						}
					}
					q += `AND `
				}
				q = q[:len(q)-4]
				q += `) `
			case "record":
				elem := strings.Split(value.String(), ",")
				for _, e := range elem {
					data := strings.Split(e, ":")
					if len(param) > 1 {
						for _, v := range param {
							if tableField == v {
								q += qp.TableAlias + tableField + `->>'` + data[0] + `' ILIKE '%` + data[1] + `%' `
								break
							}
						}

					} else {
						if len(tableField) > 0 {
							q += qp.TableAlias + tableField + `->>'` + data[0] + `' ILIKE '%` + data[1] + `%' `
						}
					}
					q += `AND `
				}
				q = q[:len(q)-4]
			default:
				if len(param) > 1 {
					for _, v := range param {
						if tableField == v {
							q += qp.TableAlias + tableField + ` = '` + value.String() + `' `
							break
						}
					}

				} else {
					if len(tableField) > 0 {
						q += qp.TableAlias + tableField + ` = '` + value.String() + `' `
					}
				}
			}
			q += `AND `
		}
	}

	q = q[:len(q)-4]
	return q
}

func arrayToQueryString(arr []string) string {
	if len(arr) < 1 {
		return ""
	}

	str := ""
	for _, v := range arr {
		str += "'" + v + "', "
	}
	return str[:len(str)-2]
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
