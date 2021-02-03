package trigger

// Type represents trigger type
type Type uint32

// Fields to be sent to the hook
type Fields map[string]interface{}

// Trigger object. Used to register hooks and custom types
type Trigger struct {
	Types []Type
	Hooks TypeHooks
}

var trigger *Trigger

func init() {
	trigger = New()
}

// New creates a new trigger
func New() *Trigger {
	return &Trigger{
		Hooks: TypeHooks{},
	}
}

// WithFields exposed as static function
func WithFields(fields Fields) *Entry {
	return trigger.WithFields(fields)
}

// AddHook exposed as static function
func AddHook(hook Hook, types []Type) {
	trigger.AddHook(hook, types)
}

// AddHook adds hook to the trigger
func (trigger *Trigger) AddHook(hook Hook, types []Type) {
	trigger.Hooks.Add(hook, types)
}

// WithFields call with fields for each field
func (trigger *Trigger) WithFields(fields Fields) *Entry {
	entry := &Entry{Trigger: trigger, Data: Fields{}}
	return entry.WithFields(fields)
}

// Fire triggers hooks based on type
func (trigger *Trigger) Fire(t Type, args ...interface{}) {
	entry := &Entry{Trigger: trigger, Data: Fields{}}
	entry.Fire(t, args...)
}
