package utils

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type KendoSort struct {
	Field string `json:"field"`
	Dir   string `json:"dir"`
}

type KendoFilter struct {
	Logic   string           `json:"logic"`
	Filters []KendoSubFilter `json:"filters"`
}

type KendoSubFilter struct {
	Field    string      `json:"field"`
	Operator string      `json:"operator"`
	Value    interface{} `json:"value"`
}

var allowedFieldPattern = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_.]*$`)

func validateField(field string) string {
	if !allowedFieldPattern.MatchString(field) {
		return ""
	}
	return field
}

func validateSortDir(dir string) string {
	switch strings.ToLower(dir) {
	case "asc":
		return "ASC"
	case "desc":
		return "DESC"
	default:
		return "ASC"
	}
}

func SortGenerator(sSort string) ([]KendoSort, string, error) {
	var kendoSort []KendoSort
	var query string
	if err := json.Unmarshal([]byte(sSort), &kendoSort); err != nil {
		return kendoSort, query, err
	}
	if len(kendoSort) == 0 {
		return kendoSort, query, nil
	}
	kendoItem := kendoSort[0]
	field := validateField(kendoItem.Field)
	if field == "" {
		return kendoSort, query, fmt.Errorf("invalid field name: %s", kendoItem.Field)
	}
	dir := validateSortDir(kendoItem.Dir)
	query = fmt.Sprintf(` order by "%s" %s`, field, dir)
	return kendoSort, query, nil
}

func FilterGenerator(sFilter string) (KendoFilter, string, []interface{}, error) {
	var filter KendoFilter
	var params []interface{}

	if err := json.Unmarshal([]byte(sFilter), &filter); err != nil {
		return filter, "", nil, err
	}

	if len(filter.Filters) == 0 {
		return filter, "", nil, nil
	}

	paramIdx := 1
	var parts []string

	for _, filterItem := range filter.Filters {
		field := validateField(filterItem.Field)
		if field == "" {
			continue
		}

		var fragment string
		switch filterItem.Value.(type) {
		case int, float64:
			fragment, params, paramIdx = buildNumericFilter(field, filterItem.Operator, filterItem.Value, params, paramIdx)
		case string:
			fragment, params, paramIdx = buildStringFilter(field, filterItem.Operator, filterItem.Value.(string), params, paramIdx)
		case bool:
			fragment, params, paramIdx = buildBoolFilter(field, filterItem.Operator, filterItem.Value.(bool), params, paramIdx)
		default:
			val := fmt.Sprintf("%v", filterItem.Value)
			fragment, params, paramIdx = buildStringFilter(field, filterItem.Operator, val, params, paramIdx)
		}

		if fragment != "" {
			parts = append(parts, fragment)
		}
	}

	if len(parts) == 0 {
		return filter, "", nil, nil
	}

	logic := "AND"
	if strings.ToLower(filter.Logic) == "or" {
		logic = "OR"
	}

	queryGenerator := "where " + strings.Join(parts, fmt.Sprintf(" %s ", logic))
	return filter, queryGenerator, params, nil
}

func buildStringFilter(field, operator, value string, params []interface{}, idx int) (string, []interface{}, int) {
	switch operator {
	case "uuid", "same":
		fragment := fmt.Sprintf(`"%s" = $%d`, field, idx)
		params = append(params, value)
		return fragment, params, idx + 1
	case "contains":
		fragment := fmt.Sprintf(`LOWER("%s") LIKE LOWER($%d)`, field, idx)
		params = append(params, "%"+value+"%")
		return fragment, params, idx + 1
	case "doesnotcontain":
		fragment := fmt.Sprintf(`LOWER("%s") NOT LIKE LOWER($%d)`, field, idx)
		params = append(params, "%"+value+"%")
		return fragment, params, idx + 1
	case "startswith":
		fragment := fmt.Sprintf(`LOWER("%s") LIKE LOWER($%d)`, field, idx)
		params = append(params, value+"%")
		return fragment, params, idx + 1
	case "endswith":
		fragment := fmt.Sprintf(`LOWER("%s") LIKE LOWER($%d)`, field, idx)
		params = append(params, "%"+value)
		return fragment, params, idx + 1
	case "eq":
		timeValue, err := time.Parse("2006-01-02T15:04:05Z", value)
		if err != nil {
			fragment := fmt.Sprintf(`LOWER("%s") = LOWER($%d)`, field, idx)
			params = append(params, value)
			return fragment, params, idx + 1
		}
		timeStart := timeValue.Local().Format("2006-01-02 15:04:05")
		timeEnd := timeValue.Add(23*time.Hour + 59*time.Minute + 59*time.Second).Local().Format("2006-01-02 15:04:05")
		fragment := fmt.Sprintf(`("%s" BETWEEN $%d AND $%d)`, field, idx, idx+1)
		params = append(params, timeStart, timeEnd)
		return fragment, params, idx + 2
	case "neq":
		timeValue, err := time.Parse("2006-01-02T15:04:05Z", value)
		if err != nil {
			return "", params, idx
		}
		timeStart := timeValue.Local().Format("2006-01-02 15:04:05")
		timeEnd := timeValue.Add(23*time.Hour + 59*time.Minute + 59*time.Second).Local().Format("2006-01-02 15:04:05")
		fragment := fmt.Sprintf(`("%s" NOT BETWEEN $%d AND $%d)`, field, idx, idx+1)
		params = append(params, timeStart, timeEnd)
		return fragment, params, idx + 2
	case "gt":
		timeValue, err := time.Parse("2006-01-02T15:04:05Z", value)
		if err != nil {
			return "", params, idx
		}
		timeEnd := timeValue.Add(23*time.Hour + 59*time.Minute + 59*time.Second).Local().Format("2006-01-02 15:04:05")
		fragment := fmt.Sprintf(`"%s" > $%d`, field, idx)
		params = append(params, timeEnd)
		return fragment, params, idx + 1
	case "gte":
		timeValue, err := time.Parse("2006-01-02T15:04:05Z", value)
		if err != nil {
			return "", params, idx
		}
		timeStart := timeValue.Local().Format("2006-01-02 15:04:05")
		fragment := fmt.Sprintf(`"%s" >= $%d`, field, idx)
		params = append(params, timeStart)
		return fragment, params, idx + 1
	case "lt":
		timeValue, err := time.Parse("2006-01-02T15:04:05Z", value)
		if err != nil {
			return "", params, idx
		}
		timeStart := timeValue.Local().Format("2006-01-02 15:04:05")
		fragment := fmt.Sprintf(`"%s" < $%d`, field, idx)
		params = append(params, timeStart)
		return fragment, params, idx + 1
	case "lte":
		timeValue, err := time.Parse("2006-01-02T15:04:05Z", value)
		if err != nil {
			return "", params, idx
		}
		timeEnd := timeValue.Add(23*time.Hour + 59*time.Minute + 59*time.Second).Local().Format("2006-01-02 15:04:05")
		fragment := fmt.Sprintf(`"%s" <= $%d`, field, idx)
		params = append(params, timeEnd)
		return fragment, params, idx + 1
	case "time":
		fragment := fmt.Sprintf(`TO_CHAR(("%s" AT TIME ZONE 'UTC')::time, 'HH24:MI') = $%d`, field, idx)
		params = append(params, value)
		return fragment, params, idx + 1
	default:
		return "", params, idx
	}
}

func buildNumericFilter(field, operator string, value interface{}, params []interface{}, idx int) (string, []interface{}, int) {
	switch operator {
	case "eq":
		fragment := fmt.Sprintf(`"%s" = $%d`, field, idx)
		params = append(params, value)
		return fragment, params, idx + 1
	case "neq":
		fragment := fmt.Sprintf(`"%s" <> $%d`, field, idx)
		params = append(params, value)
		return fragment, params, idx + 1
	case "gt":
		fragment := fmt.Sprintf(`"%s" > $%d`, field, idx)
		params = append(params, value)
		return fragment, params, idx + 1
	case "gte":
		fragment := fmt.Sprintf(`"%s" >= $%d`, field, idx)
		params = append(params, value)
		return fragment, params, idx + 1
	case "lt":
		fragment := fmt.Sprintf(`"%s" < $%d`, field, idx)
		params = append(params, value)
		return fragment, params, idx + 1
	case "lte":
		fragment := fmt.Sprintf(`"%s" <= $%d`, field, idx)
		params = append(params, value)
		return fragment, params, idx + 1
	default:
		return "", params, idx
	}
}

func buildBoolFilter(field, operator string, value bool, params []interface{}, idx int) (string, []interface{}, int) {
	switch operator {
	case "eq":
		fragment := fmt.Sprintf(`"%s" = $%d`, field, idx)
		params = append(params, value)
		return fragment, params, idx + 1
	default:
		return "", params, idx
	}
}

func FilterGeneratorNonKendo(sFilter string) (KendoFilter, error) {
	var filter KendoFilter
	if err := json.Unmarshal([]byte(sFilter), &filter); err != nil {
		return filter, err
	}
	return filter, nil
}

func KendoSortFilter(vSort, vFilters string) (sSorts, sFilters string, filterParams []interface{}, err error) {
	if vSort != "" {
		_, sSorts, err = SortGenerator(vSort)
		if err != nil {
			return sSorts, sFilters, nil, err
		}
	}

	if vFilters != "" {
		_, sFilters, filterParams, err = FilterGenerator(vFilters)
		if err != nil {
			return sSorts, sFilters, nil, err
		}
	}
	return sSorts, sFilters, filterParams, nil
}

func FilterSpecific(value, param, operator, filter string) (string, []interface{}) {
	field := validateField(param)
	if field == "" {
		return filter, nil
	}

	if len(filter) > 0 {
		return fmt.Sprintf("%s %s \"%s\" = $1", filter, operator, field), []interface{}{value}
	}
	return fmt.Sprintf(`WHERE "%s" = $1`, field), []interface{}{value}
}
