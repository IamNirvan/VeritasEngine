package facts

import "github.com/IamNirvan/veritasengine/internal/models/response"

type GeneralInput struct {
	Input    *interface{}   `json:"input"`
	Response *[]interface{} `json:"response"`
}

func NewFact() *GeneralInput {
	return &GeneralInput{
		Input:    nil,
		Response: &[]interface{}{},
	}
}

func (gi *GeneralInput) StringListHaveMatchingItems(listA []string, listB []string) bool {
	itemMap := make(map[string]struct{}, len(listA))

	for _, itemA := range listA {
		itemMap[itemA] = struct{}{}
	}
	for _, itemB := range listB {
		if _, found := itemMap[itemB]; found {
			return true
		}
	}
	return false
}

func (gi *GeneralInput) AddToResponse(responseType int64, statusType int64, data interface{}) {
	*gi.Response = append(*gi.Response, response.RuleEvaluationResponse{
		ResponseType: responseType,
		StatusType:   statusType,
		Payload:      &data,
	})
}
