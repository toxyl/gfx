package meta

type FilterMetaDataArg struct {
	Name    string
	Default any
}

type FilterMetaData struct {
	Name string
	Args []*FilterMetaDataArg // name to default
}

func (m *FilterMetaData) Arg(i int) *FilterMetaDataArg {
	if i > len(m.Args)-1 {
		return &FilterMetaDataArg{
			Name:    "undefined",
			Default: nil,
		}
	}
	return m.Args[i]
}

func (m *FilterMetaData) NameOf(i int) string { return m.Arg(i).Name }
func (m *FilterMetaData) DefaultOf(i int) any { return m.Arg(i).Default }

func (m *FilterMetaData) ArgNames() []string {
	res := []string{}
	for _, a := range m.Args {
		res = append(res, a.Name)
	}
	return res
}

func New(name string, args []*FilterMetaDataArg) *FilterMetaData {
	m := FilterMetaData{
		Name: name,
		Args: args,
	}
	return &m
}
