package kube_secrets_exporter

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	Stderr = File("stderr")
	Stdout = File("stdout")
)

type File string

func (instance *File) Set(plain string) error {
	switch File(strings.ToLower(plain)) {
	case "":
		*instance = ""
		return nil
	case Stderr:
		*instance = Stderr
		return nil
	case Stdout:
		*instance = Stdout
		return nil
	default:
		fi, err := os.Stat(plain)
		if err != nil && !os.IsNotExist(err) {
			return err
		} else if fi != nil && fi.IsDir() {
			return fmt.Errorf("should be file but is directory: %s", plain)
		}
		*instance = File(plain)
		return nil
	}
}

func (instance File) String() string {
	return string(instance)
}

func (instance File) Open() (io.WriteCloser, error) {
	switch instance {
	case "":
		return &noopWriter{}, nil
	case Stdout:
		return &fileWriterWrapper{os.Stdout, false}, nil
	case Stderr:
		return &fileWriterWrapper{os.Stderr, false}, nil
	default:
		dir := filepath.Dir(instance.String())
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("cannot ensure directory for '%v': %w", instance, err)
		}
		if f, err := os.OpenFile(instance.String(), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
			return nil, fmt.Errorf("cannot open '%v': %w", instance, err)
		} else {
			return f, nil
		}
	}
}

type noopWriter struct {
}

func (instance *noopWriter) Write([]byte) (n int, err error) {
	return 0, nil
}

func (instance *noopWriter) Close() error {
	return nil
}

type fileWriterWrapper struct {
	io.WriteCloser
	canClose bool
}

func (instance *fileWriterWrapper) Close() error {
	if instance.canClose {
		return instance.WriteCloser.Close()
	}
	return nil
}
