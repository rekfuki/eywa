package trigger

// Hook interface allows to attach custom hooks to be triggered based on type
type Hook interface {
	Fire(*Entry) error
}

// TypeHooks for storing hooks for type
type TypeHooks map[Type][]Hook

// Add adds a hook to specific trigger type
func (hooks TypeHooks) Add(hook Hook, types []Type) {
	for _, t := range types {
		hooks[t] = append(hooks[t], hook)
	}
}

// Fire all the hooks attached to that trigger type
func (hooks TypeHooks) Fire(t Type, entry *Entry) error {
	for _, hook := range hooks[t] {
		if err := hook.Fire(entry); err != nil {
			return err
		}
	}

	return nil
}
