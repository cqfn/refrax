package aibolit

import (
	"os/exec"

	"github.com/cqfn/refrax/internal/log"
)

type DefaultAibolit struct {
	filename string
}

func NewDefaultAibolit(filename string) Aibolit {
	return &DefaultAibolit{filename}
}

func (a* DefaultAibolit) Imperfections() string {
	cmd := exec.Command("aibolit", "check", "--filenames", "Foo.java")
	opportunities, _ := cmd.CombinedOutput()
    log.Debug("Identified refactoring opportunities with aibolit: \n%s", opportunities)
    return string(opportunities)
}
