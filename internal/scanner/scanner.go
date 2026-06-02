package scanner

type Options struct {
	ProjectPath string
	AssetsDir   string
	IgnoreDirs  []string
	verbose     bool
}

type Scanner struct {
	opts Options
}

func NewScanner(opts Options) *Scanner {
	return &Scanner{
		opts: opts,
	}
}
