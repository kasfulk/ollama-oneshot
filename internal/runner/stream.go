package runner

import (
	"io"
	"os"
)

func streamOutput(src io.Reader) error {
	_, err := io.Copy(os.Stdout, src)
	return err
}