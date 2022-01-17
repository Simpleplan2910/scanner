package git

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"unicode/utf8"
)

type Repos interface {
	GetTextFiles() (l []string, err error)
	ReadFile(filename string) (reader io.Reader, err error)
}

type repos struct {
	dir string
}

func (r *repos) GetTextFiles() (l []string, err error) {
	l = []string{}
	err = filepath.Walk(r.dir, func(path string, infos fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if infos.IsDir() && infos.Name() == ".git" {
			// skip .git direction
			return filepath.SkipDir
		}
		if infos.IsDir() {
			// skip direction so it's not append to list files
			return nil
		}

		// check if file is text file
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		if IsTextFile(f) {
			l = append(l, path)
		}
		return nil
	})
	return l, err
}

func (r *repos) ReadFile(filename string) (reader io.Reader, err error) {
	return os.Open(filename)
}

// steal from godoc/util/util.go with a little of modifying

// IsText reports whether a significant prefix of s looks like correct UTF-8;
// that is, if it is likely that s is human-readable text.
func IsText(s []byte) bool {
	const max = 1024 // at least utf8.UTFMax
	if len(s) > max {
		s = s[0:max]
	}
	for i, c := range string(s) {
		if i+utf8.UTFMax > len(s) {
			// last char may be incomplete - ignore
			break
		}
		if c == 0xFFFD || c < ' ' && c != '\n' && c != '\t' && c != '\f' {
			// decoding error or control character - not a text file
			return false
		}
	}
	return true
}

// IsTextFile reports whether the file has a known extension indicating
// a text file, or if a significant chunk of the specified file looks like
// correct UTF-8; that is, if it is likely that the file contains human-
// readable text.
func IsTextFile(f io.Reader) bool {
	// the extension is not known; read an initial chunk
	// of the file and check if it looks like text
	var buf [1024]byte
	n, err := f.Read(buf[0:])
	if err != nil {
		return false
	}

	return IsText(buf[0:n])
}
