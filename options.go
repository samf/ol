package main

type options struct {
	debug     bool
	workers   int
	ignoreVCS bool

	filter filter
}

func (opt *options) valid() error {
	opt.filter = noopFilter
	if opt.ignoreVCS {
		opt.filter = opt.filter.noGit().noHG()
	}

	return nil
}
