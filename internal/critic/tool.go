package critic

// @todo #2:45min Implement support for Tool with multiple classes.
//  As for now, we check only the first class, and return imperfections result. Instead, we need to support
//  multiple files instead. Let's implement such Tool struct, that will be able to manage whole project, instead
//  of single Java file. Also see this related issue: https://github.com/cqfn/refrax/issues/28.
type Tool interface {
	Imperfections() string
}
