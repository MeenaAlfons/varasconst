package analyzer

type VarMarkedConstFact struct {
}

func (*VarMarkedConstFact) AFact() {}

func (f *VarMarkedConstFact) String() string {
	return "var marked as const"
}
