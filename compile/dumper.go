package compile

type dumper interface {
	WritePhase(string, string)
	WriteSources(phase string, fn string, lines []string, startline int)
	Close()
}

type nopDumper struct{}

func (nopDumper) WritePhase(string, string)                                           {}
func (nopDumper) WriteSources(phase string, fn string, lines []string, startline int) {}
func (nopDumper) Close()                                                              {}
