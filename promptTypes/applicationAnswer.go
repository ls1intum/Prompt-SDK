package promptTypes

import (
	"encoding/json"
	"fmt"
)

const (
	TypeMultiSelect = "multiselect"
	TypeText        = "text"
)

type AnswersMultiSelect struct {
	Type     string   `json:"type"`
	OrderNum int      `json:"order_num"`
	Key      string   `json:"key"`
	Answer   []string `json:"answer"`
}

type AnswersText struct {
	Type     string `json:"type"`
	OrderNum int    `json:"order_num"`
	Key      string `json:"key"`
	Answer   string `json:"answer"`
}

// rawAnswer is the intermediate struct used for unmarshaling the raw JSON.
type rawAnswer struct {
	Type     string      `json:"type"`
	OrderNum int         `json:"order_num"`
	Answer   interface{} `json:"answer"` // can be a string or list of strings
	Key      string      `json:"key"`
}

// Allows to parse the answers from the metadata in any later phase.
// 1) []AnswersText for type="text"
// 2) []AnswersMultiSelect for type="multiselect"
func ReadApplicationAnswersFromMetaData(data interface{}) ([]AnswersText, []AnswersMultiSelect, error) {
	// Step 1: Marshal interface{} into JSON
	rawBytes, err := json.Marshal(data)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	// Step 2: Unmarshal into a slice of rawAnswer
	var rawAnswers []rawAnswer
	if err := json.Unmarshal(rawBytes, &rawAnswers); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal into rawAnswers: %w", err)
	}

	// Separate the answers into textAnswers and multiAnswers
	var (
		textAnswers  []AnswersText
		multiAnswers []AnswersMultiSelect
	)

	// Step 3: Convert each rawAnswer to the correct final struct
	for _, ra := range rawAnswers {
		switch ra.Type {
		case TypeMultiSelect:
			// Expect ra.Answer to be an array of strings
			arr, ok := ra.Answer.([]interface{})
			if !ok {
				return nil, nil, fmt.Errorf("multiselect answer was not an array: %v", ra.Answer)
			}
			stringArr := make([]string, 0, len(arr))
			for _, item := range arr {
				s, ok := item.(string)
				if !ok {
					return nil, nil, fmt.Errorf("multiselect array contains non-string element: %v", item)
				}
				stringArr = append(stringArr, s)
			}
			multiAnswers = append(multiAnswers, AnswersMultiSelect{
				Type:     ra.Type,
				OrderNum: ra.OrderNum,
				Key:      ra.Key,
				Answer:   stringArr,
			})

		case TypeText:
			// Expect ra.Answer to be a single string
			s, ok := ra.Answer.(string)
			if !ok {
				return nil, nil, fmt.Errorf("text answer was not a string: %v", ra.Answer)
			}
			textAnswers = append(textAnswers, AnswersText{
				Type:     ra.Type,
				OrderNum: ra.OrderNum,
				Key:      ra.Key,
				Answer:   s,
			})

		default:
			return nil, nil, fmt.Errorf("unknown answer type: %s", ra.Type)
		}
	}

	return textAnswers, multiAnswers, nil
}
