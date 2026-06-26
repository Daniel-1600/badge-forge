package rules

type EventType int

const (
	BadgeCreated EventType = iota
	AssertionCreated
)

type Event struct {
	Type     EventType
	BadgeID  string
	PersonID int
}

type Rule interface {
	// Evaluate evaluates the rule against the given event and returns true if the rule is satisfied, false otherwise.
	Evaluate(event Event) bool
}
