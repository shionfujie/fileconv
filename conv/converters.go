package conv

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Converter struct {
	SourceType string
	TargetType string
	file       string
	cmd        *exec.Cmd
	pathErr    error
}

type cmd struct {
	name       string
	targetType string
}

var (
	djvu2pdf = cmd{"djvu2pdf", "PDF"}
)

var cmds = map[string]cmd{
	"DJVU": djvu2pdf,
}

func NewConverter(file string) (c *Converter) {
	c = &Converter{
		file: file,
	}
	c.SourceType = sourceType(file)
	if _, err := os.Stat(file); err != nil {
		c.pathErr = err.(*os.PathError).Err
		return
	}

	cmd, ok := cmds[c.SourceType]
	if !ok {
		c.pathErr = fmt.Errorf("unsupported file type %s", c.SourceType)
		return
	}

	c.TargetType = cmd.targetType
	c.cmd = exec.Command(cmd.name, file)
	return
}

func sourceType(file string) string {
	return strings.ToUpper(filepath.Ext(file)[1:])
}

func (c *Converter) Convert() error {
	if c.pathErr != nil {
		return c.pathErr
	}
	if _, err := c.cmd.Output(); err != nil {
		return fmt.Errorf("failed to convert %s file to %s file: %v", c.SourceType, c.TargetType, err)
	}
	return nil
}

func convert(file string) error {
	if _, err := os.Stat(file); err != nil {
		return err.(*os.PathError).Err
	}
	if _, err := exec.Command("djvu2pdf", file).Output(); err != nil {
		return fmt.Errorf("failed to convert DJVU file to PDF file: %v", err)
	}
	return nil
}
