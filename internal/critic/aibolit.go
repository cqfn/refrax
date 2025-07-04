package critic

import (
	"os/exec"

	"github.com/cqfn/refrax/internal/log"
)

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
