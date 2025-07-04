package critic

import (
	"os/exec"

	"github.com/cqfn/refrax/internal/log"
)

// @todo #2:45min Implement support for Aibolit with multiple classes.
//  As for now, we check only the first class, and return imperfections result. Instead, we need to support
//  multiple files instead. Let's implement such Aibolit struct, that will be able to manage whole project, instead
//  of single Java file. Also see this related issue: https://github.com/cqfn/refrax/issues/28.
type Aibolit struct {
	filename string
}

func NewAibolit(filename string) Tool {
	return &Aibolit{filename}
}

func (a* Aibolit) Imperfections() string {
	cmd := exec.Command("aibolit", "check", "--filenames", "Foo.java")
	opportunities, _ := cmd.CombinedOutput()
    log.Debug("Identified refactoring opportunities with aibolit: \n%s", opportunities)
    return string(opportunities)
}
