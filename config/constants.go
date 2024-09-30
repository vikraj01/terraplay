package config

type Status int

const (
	Pending Status = iota
	Running
	Terminated
)

var StatusNames = map[Status]string{
	Pending:    "pending",
	Running:    "running",
	Terminated: "terminated",
}

var StatusValues = map[string]Status{
	"pending":    Pending,
	"running":    Running,
	"terminated": Terminated,
}

var Timeout = 30
