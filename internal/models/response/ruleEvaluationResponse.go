package response

type RuleEvaluationResponse struct {
	ResponseType int
	StatusType   int
	Payload      *interface{}
}
