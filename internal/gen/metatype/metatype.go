package metatype

type Interface struct {
	Docs     []string
	Name     string
	Embedded []string
	Funcs    []Func
}

// EmbedInterface appends interfaceName to MetaInterface.Embedded
// if not already exists, else it is ignored.
func (mi *Interface) EmbedInterface(interfaceName string) {
	for _, itf := range mi.Embedded {
		if itf == interfaceName {
			return
		}
	}
	mi.Embedded = append(mi.Embedded, interfaceName)
}

// AddFunc creates a MetaFunc from *doc.Func and appends it
// to MetaInterface.Funcs.
func (mi *Interface) AddFunc(f Func) {
	mi.Funcs = append(mi.Funcs, f)
}

type Func struct {
	Sign string
	Docs []string
}

type Var struct {
	Name, Type, Value string
}
