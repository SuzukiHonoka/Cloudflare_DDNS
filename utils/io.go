package utils

import "io"

func ForceClose(closer io.Closer) {
	_ = closer.Close()
}
