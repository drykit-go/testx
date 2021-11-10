package metatype

// Var gathers a variable informations so it can be generated
// via a go template.
type Var struct {
	Name, Value string
	DocLines    []string
}
