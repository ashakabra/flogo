package queryticket

import (
	"encoding/json"
	"fmt"
	"reflect"

	"git.tibco.com/git/product/ipaas/wi-contrib.git/engine/conversion"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

type Parameters struct {
	QueryParams []*TypedValue `json:"queryParams"`
}

type TypedValue struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
	//	Type  string      `json:"type"`
}

type Param struct {
	Name      string
	Type      string
	Repeating string
	Required  string
}

func ParseParams(paramSchema map[string]interface{}) ([]Param, error) {

	if paramSchema == nil {
		return nil, nil
	}

	var parameter []Param

	//Structure expected to be JSON schema like
	props := paramSchema["properties"].(map[string]interface{})
	for k, v := range props {
		param := &Param{}
		param.Name = k
		propValue := v.(map[string]interface{})
		for k1, v1 := range propValue {
			if k1 == "required" {
				param.Required = v1.(string)
			} else if k1 == "type" {
				if v1 != "array" {
					param.Repeating = "false"
				}
				param.Type = v1.(string)
			} else if k1 == "items" {
				param.Repeating = "true"
				items := v1.(map[string]interface{})
				s, err := conversion.ConvertToString(items["type"])
				if err != nil {
					return nil, err
				}
				param.Type = s
			}
		}
		parameter = append(parameter, *param)
	}

	return parameter, nil
}

func GetParameter(valueIN interface{}) (params *Parameters, err error) {
	params = &Parameters{}
	//QueryParams
	activityLog.Info("Reading Query parameters")
	queryParamsMap, _ := LoadJsonSchemaFromMetadata(valueIN)
	if queryParamsMap != nil {
		queryParams, err := ParseParams(queryParamsMap)
		if err != nil {
			return params, err
		}

		if queryParams != nil {
			inputQueryParams, err := GetComplexValueAsMap(valueIN)
			if err != nil {
				return params, err
			}
			var typeValuesQueryParams []*TypedValue
			for _, qParam := range queryParams {
				isRequired := qParam.Required
				paramName := qParam.Name
				if isRequired == "true" && inputQueryParams[paramName] == nil {
					return nil, fmt.Errorf("Required query parameter [%s] is not configured.", paramName)
				}
				if inputQueryParams[paramName] != nil {
					if qParam.Repeating == "true" {
						val := inputQueryParams[paramName]

						switch reflect.TypeOf(val).Kind() {
						case reflect.Slice:
							s := reflect.ValueOf(val)
							//working for array
							var value string
							for i := 0; i < s.Len(); i++ {
								stringValue := fmt.Sprint(s.Index(i).Interface())
								value = value + stringValue + ", "
							}
							value = value[0 : len(value)-2] //remove last extra comma

							typeValue := &TypedValue{}
							typeValue.Name = paramName
							//typeValue.Type = hParam.Type
							typeValue.Value = value
							typeValuesQueryParams = append(typeValuesQueryParams, typeValue)
						}
					} else {
						typeValue := &TypedValue{}
						typeValue.Name = paramName
						typeValue.Value = inputQueryParams[paramName]
						//typeValue.Type = hParam.Type
						typeValuesQueryParams = append(typeValuesQueryParams, typeValue)
					}
					params.QueryParams = typeValuesQueryParams
				}
			}
		}
	}

	return params, err
}

func LoadJsonSchemaFromMetadata(valueIN interface{}) (map[string]interface{}, error) {
	if valueIN != nil {
		complex := valueIN.(*data.ComplexObject)
		if complex != nil {
			params, err := convertToMap(complex.Metadata)
			if err != nil {
				return nil, err
			}
			return params, nil
		}
	}
	return nil, nil
}

func GetComplexValueAsMap(valueIN interface{}) (map[string]interface{}, error) {
	if valueIN != nil {
		complex := valueIN.(*data.ComplexObject)
		if complex != nil {
			switch t := complex.Value.(type) {
			case string:
				m := map[string]interface{}{}
				err := json.Unmarshal([]byte(t), &m)
				if err != nil {
					return nil, err
				}
				return m, nil
			default:
				return convertToMap(complex.Value)

			}
		}
	}
	return nil, nil
}

func convertToMap(data interface{}) (map[string]interface{}, error) {
	switch t := data.(type) {
	case string:
		if t != "" {
			m := map[string]interface{}{}
			err := json.Unmarshal([]byte(t), &m)
			if err != nil {
				return nil, err
			}
			return m, nil
		}
	case map[string]interface{}:
		return t, nil
	case interface{}:
		b, err := json.Marshal(t)
		if err != nil {
			return nil, err
		}
		m := map[string]interface{}{}
		err = json.Unmarshal(b, &m)
		if err != nil {
			return nil, err
		}
		return m, nil
	}

	return nil, nil
}
