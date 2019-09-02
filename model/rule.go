package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"
)

type (
	//Rules :
	Rules struct {
		Rules RulesExpression `json:"rule"`
	}
	//RulesExpression :
	RulesExpression struct {
		Or  []map[string]RulesArgument `json:"$OR"`
		And []map[string]RulesArgument `json:"$AND"`
	}
	//RulesArgument :
	// Gte : Greater Than Equals (>=)
	// Lte : Lesser Than Equals (<=)
	// In : Include in, kind of array which compare data `in` expression generaly (*sql where `IN` ('string1', 'string2', 'string3'))
	// Gt : Greater Than (>)
	// Lt : Lesser Than (<)
	// Eq : Equal
	RulesArgument struct {
		Gte interface{}   `json:"$gte,omitempty"`
		Lte interface{}   `json:"$lte,omitempty"`
		In  []interface{} `json:"$in,omitempty"`
		Gt  interface{}   `json:"$gt,omitempty"`
		Lt  interface{}   `json:"$lt,omitempty"`
		Eq  interface{}   `json:"$eq,omitempty"`
	}
	//RulesInput : transactional data requirement for validate rule
	RulesInput struct {
		ProgramPeriod *time.Time
		VoucherPeriod *time.Time
	}
)

const (
	ruleActiveProgramPeriod = "active_program_period"
	ruleValidVoucherPeriod  = "valid_voucher_period"
	ruleAllowCrossProgram   = "cross_program"

	//allow accumulative change to number 0 = unlimited
	ruleAccumulative = "allow_accumulative"

	//rule max usage by day
	ruleMaxUsageByDay = "max_usage_by_day"

	ruleSpending      = "spending"
	ruleValidityHours = "validity_hours"
	ruleValidityDays  = "validity_days"

	//GET Voucher
	//TODO
	ruleGetVoucherType = "get_voucher_type" //1days, program
	ruleMaxGetVoucher  = "max_get_voucher"  //2
)

var (
	//ErrorRuleUnexpectedTimeFormat :
	ErrorRuleUnexpectedTimeFormat = errors.New("Invalid string time format could not be converted to time")
	//ErrorRuleUnexpectedNumericType :
	ErrorRuleUnexpectedNumericType = errors.New("Non-numeric type could not be converted to float")
)

// Unmarshal from JSONExpr.String to Rule struct
func (rule *Rules) Unmarshal(exp JSONExpr) error {
	err := json.Unmarshal([]byte(exp.String()), rule)
	return err
}

func (rule *Rules) getRuleKeys() ([]string, error) {
	//get keys of rules
	keys := make([]string, len(rule.Rules.And))

	i := 0
	for _, v := range rule.Rules.And {
		for k := range v {
			keys[i] = k
			i++
			break
		}
	}
	return keys, nil
}

//Validate Rules
func (rule *Rules) Validate() (bool, error) {
	r := false
	for _, v := range rule.Rules.And {
		r, err := rule.validateRulesAnd(v)
		if !r {
			return r, err
		}
	}
	return r, nil
}

//DEPRECATED
//DEPRECATED
//DEPRECATED
func (rule *Rules) validateRulesAnd(ra map[string]RulesArgument) (bool, error) {
	r := false
	for k, v := range ra {
		switch k {
		case ruleActiveProgramPeriod:
			v.validateTime(time.Now())
			break
		case ruleValidVoucherPeriod:
			v.validateTime(time.Now())
			break
		case ruleAllowCrossProgram:
			v.validateString("")
			break
		case ruleAccumulative:
			//validate bool
			break
		case ruleMaxUsageByDay:
			v.validateNumber(0)
			break
		case ruleSpending:
			v.validateNumber(0)
			break
		case ruleValidityHours:
			//??
			break
		case ruleValidityDays:
			v.validateNumber(0)
			break
		}
	}
	return r, nil
}

//Default true if 1 validate return false then return `not valid transaction rule`
func (ra *RulesArgument) validateTime(tx time.Time) (bool, error) {
	r := true
	if ra.Eq != nil {
		r, err := ra.validateEqTime(tx)
		if err != nil {
			return r, err
		}
	}
	if ra.Gte != nil {
		r, err := ra.validateGteTime(tx)
		if err != nil {
			return r, err
		}
	}
	if ra.Lte != nil {
		r, err := ra.validateLteTime(tx)
		if err != nil {
			return r, err
		}
	}
	if ra.Gt != nil {
		r, err := ra.validateGtTime(tx)
		if err != nil {
			return r, err
		}
	}
	if ra.Lt != nil {
		r, err := ra.validateLtTime(tx)
		if err != nil {
			return r, err
		}
	}
	if ra.In != nil {
		r, err := ra.validateInTime(tx)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

func (ra *RulesArgument) validateEqTime(tx time.Time) (bool, error) {
	r := false
	//convert to time
	t, err := stringToTime(fmt.Sprint(ra.Gte))
	if err != nil {
		return r, err
	}
	//validate rule
	r = tx.Equal(t)
	return r, nil
}

func (ra *RulesArgument) validateGteTime(tx time.Time) (bool, error) {
	r := false
	//convert to time
	t, err := stringToTime(fmt.Sprint(ra.Gte))
	if err != nil {
		return r, err
	}
	//validate rule
	r = tx.Equal(t) || tx.After(t)
	return r, nil
}

func (ra *RulesArgument) validateLteTime(tx time.Time) (bool, error) {
	r := false
	//convert to time
	t, err := stringToTime(fmt.Sprint(ra.Gte))
	if err != nil {
		return r, err
	}
	//validate rule
	r = tx.Equal(t) || tx.Before(t)
	return r, nil
}

func (ra *RulesArgument) validateGtTime(tx time.Time) (bool, error) {
	r := false
	//convert to time
	t, err := stringToTime(fmt.Sprint(ra.Gte))
	if err != nil {
		return r, err
	}
	//validate rule
	r = tx.After(t)
	return r, nil
}

func (ra *RulesArgument) validateLtTime(tx time.Time) (bool, error) {
	r := false
	//convert to time
	t, err := stringToTime(fmt.Sprint(ra.Gte))
	if err != nil {
		return r, err
	}
	//validate rule
	r = tx.Before(t)
	return r, nil
}

func (ra *RulesArgument) validateInTime(tx time.Time) (bool, error) {
	r := false
	//convert to time
	for _, value := range ra.In {
		t, err := stringToTime(fmt.Sprint(value))
		if err != nil {
			return r, err
		}
		//validate rule
		r = tx.Equal(t)
		if r {
			break
		}
	}
	return r, nil
}

//Default true if 1 validate return false then return `not valid transaction rule`
func (ra *RulesArgument) validateString(val string) (bool, error) {
	r := true
	if ra.Eq != nil {
		r, err := ra.validateEqString(val)
		if err != nil {
			return r, err
		}
	}
	if ra.In != nil {
		r, err := ra.validateInString(val)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

func (ra *RulesArgument) validateEqString(val string) (bool, error) {
	r := false
	//validate rule
	r = val == ra.Eq
	return r, nil
}

func (ra *RulesArgument) validateInString(val string) (bool, error) {
	r := false
	//convert to time
	for _, value := range ra.In {
		//validate rule
		r = value == val
		if r {
			break
		}
	}
	return r, nil
}

//Default true if 1 validate return false then return `not valid transaction rule`
func (ra *RulesArgument) validateNumber(val interface{}) (bool, error) {
	r := true
	if ra.Eq != nil {
		r, err := ra.validateEqNumber(val)
		if err != nil {
			return r, err
		}
	}
	if ra.Gte != nil {
		r, err := ra.validateGteNumber(val)
		if err != nil {
			return r, err
		}
	}
	if ra.Lte != nil {
		r, err := ra.validateLteNumber(val)
		if err != nil {
			return r, err
		}
	}
	if ra.Gt != nil {
		r, err := ra.validateGtNumber(val)
		if err != nil {
			return r, err
		}
	}
	if ra.Lt != nil {
		r, err := ra.validateLtNumber(val)
		if err != nil {
			return r, err
		}
	}
	if ra.In != nil {
		r, err := ra.validateInNumber(val)
		if err != nil {
			return r, err
		}
	}
	return r, nil
}

func (ra *RulesArgument) validateEqNumber(val interface{}) (bool, error) {
	r := false
	//validate rule
	r = ra.Eq == val
	return r, nil
}

func (ra *RulesArgument) validateGteNumber(val interface{}) (bool, error) {
	r := false
	//convert to float
	gte, err := toFloat(ra.Gte)
	if err != nil {
		return r, ErrorRuleUnexpectedNumericType
	}
	c, err := toFloat(val)
	if err != nil {
		return r, ErrorRuleUnexpectedNumericType
	}
	//validate rule
	r = gte <= c
	return r, nil
}

func (ra *RulesArgument) validateLteNumber(val interface{}) (bool, error) {
	r := false
	//convert to float
	lte, err := toFloat(ra.Lte)
	if err != nil {
		return r, ErrorRuleUnexpectedNumericType
	}
	c, err := toFloat(val)
	if err != nil {
		return r, ErrorRuleUnexpectedNumericType
	}
	//validate rule
	r = lte >= c
	return r, nil
}

func (ra *RulesArgument) validateGtNumber(val interface{}) (bool, error) {
	r := false
	//convert to float
	gt, err := toFloat(ra.Gt)
	if err != nil {
		return r, ErrorRuleUnexpectedNumericType
	}
	c, err := toFloat(val)
	if err != nil {
		return r, ErrorRuleUnexpectedNumericType
	}
	//validate rule
	r = gt < c
	return r, nil
}

func (ra *RulesArgument) validateLtNumber(val interface{}) (bool, error) {
	r := false
	//convert to time
	lt, err := toFloat(ra.Gt)
	if err != nil {
		return r, ErrorRuleUnexpectedNumericType
	}
	c, err := toFloat(val)
	if err != nil {
		return r, ErrorRuleUnexpectedNumericType
	}
	//validate rule
	r = lt > c
	return r, nil
}

func (ra *RulesArgument) validateInNumber(val interface{}) (bool, error) {
	r := false
	//convert to time
	c, err := toFloat(val)
	if err != nil {
		return r, ErrorRuleUnexpectedNumericType
	}
	for _, value := range ra.In {
		//validate rule
		in, err := toFloat(value)
		if err != nil {
			return r, ErrorRuleUnexpectedNumericType
		}
		r = in == c
		if r {
			break
		}
	}
	return r, nil
}

func (rule *Rules) checkActiveProgramPeriod(tx *time.Time, arg *RulesArgument) (bool, error) {
	r := false
	//convert to time
	if arg.Gte != nil {
		t, err := stringToTime(fmt.Sprint(arg.Gte))
		if err != nil {
			return r, err
		}
		//validate rule
		r = tx.Equal(t) || tx.After(t)
	}
	return r, nil
}

func stringToTime(value string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return t, ErrorRuleUnexpectedTimeFormat
	}
	return t, nil
}

func toFloat(val interface{}) (float64, error) {
	switch i := val.(type) {
	case float64:
		return i, nil
	case float32:
		return float64(i), nil
	case int64:
		return float64(i), nil
	case int32:
		return float64(i), nil
	case int:
		return float64(i), nil
	case uint64:
		return float64(i), nil
	case uint32:
		return float64(i), nil
	case uint:
		return float64(i), nil
	default:
		return math.NaN(), ErrorRuleUnexpectedNumericType
	}
}

func (rule *Rules) isTimeValid(opr string, val time.Time) bool {
	r := false
	now := time.Now()
	switch opr {
	case "gte":
		r = now.Equal(val) || now.After(val)
		break
	case "lte":
		r = now.Equal(val) || now.Before(val)
		break
	case "gt":
		r = now.After(val)
		break
	case "lt":
		r = now.Before(val)
		break
	case "eq":
		r = now.Equal(val)
		break
	case "in":
		r = false
		break
	}
	return r
}

func (rule *Rules) isNumberValid(opr string, val, exp int) bool {
	r := false
	switch opr {
	case "gte":
		r = exp >= val
		break
	case "lte":
		r = exp <= val
		break
	case "gt":
		r = exp > val
		break
	case "lt":
		r = exp < val
		break
	case "eq":
		r = exp == val
		break
	case "in":
		r = false
		break
	}
	return r
}
