package configflag

import "strings"

// arrayFlags contains an array of command flags
type arrayFlags struct {
	items *[]string
	reset bool
}

func newArrayFlags(items *[]string) *arrayFlags {
	return &arrayFlags{
		items: items,
		reset: true,
	}
}

// Values return the values of a flag array
func (i *arrayFlags) Values() []string {
	if i.items == nil {
		return []string{}
	}
	return *i.items
}

// String return the string value of a flag array
func (i *arrayFlags) String() string {
	return strings.Join(i.Values(), ",")
}

// Set is used to add a value to the flag array
func (i *arrayFlags) Set(value string) error {
	if i.reset {
		i.reset = false
		*i.items = []string{value}
	} else {
		*i.items = append(*i.items, value)
	}
	return nil
}
