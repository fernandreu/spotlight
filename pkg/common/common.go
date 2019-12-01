package common

import (
	"bufio"
	"os"
)

func TryCloseFile(f *os.File) {
	_ = f.Close()
}

func TryFlush(w *bufio.Writer) {
	_ = w.Flush()
}
