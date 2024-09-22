package util

type Division = uint8

const (
	SupportDivision Division = iota + 1
	ItDivision
	BillingDivision
)

var Divisions = [...]string{"Support", "IT", "Billing"}

type PayloadField uint8

const (
	DivisionField PayloadField = iota
	SubjectField
	TextField
)

func (p PayloadField) String() string {
	return PayloadFields[p]
}

var PayloadFields = [...]string{"division", "subject", "text"}
