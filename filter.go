package main

type filter func(Node) bool

func (f filter) noGit() filter {
	return func(n Node) bool {
		switch {
		case n.Name() == ".git":
			return true
		case f == nil:
			return false
		}

		return f(n)
	}
}

func (f filter) noHG() filter {
	return func(n Node) bool {
		switch {
		case n.Name() == ".hg":
			return true
		case f == nil:
			return false
		}

		return f(n)
	}
}
