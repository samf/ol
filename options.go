package main

type options struct {
	debug   bool
	workers int
}

func (opt *options) valid() error {
	return nil
}
