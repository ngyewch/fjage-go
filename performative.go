package fjage

type Performative string

const (
	PerformativeAgree         Performative = "AGREE"
	PerformativeCancel        Performative = "CANCEL"
	PerformativeCFP           Performative = "CFP"
	PerformativeConfirm       Performative = "CONFIRM"
	PerformativeDisconfirm    Performative = "DISCONFIRM"
	PerformativeFailure       Performative = "FAILURE"
	PerformativeInform        Performative = "INFORM"
	PerformativeNotUnderstood Performative = "NOT_UNDERSTOOD"
	PerformativePropose       Performative = "PROPOSE"
	PerformativeQueryIf       Performative = "QUERY_IF"
	PerformativeRefuse        Performative = "REFUSE"
	PerformativeRequest       Performative = "REQUEST"
)
