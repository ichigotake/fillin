package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/mitchellh/go-homedir"
)

// Run fillin
func Run(configDir string, args []string, in *bufio.Reader, out *bufio.Writer) (string, error) {
	dir, err := homedir.Expand(normalizeFileSeparator(configDir))
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	if err := os.Chmod(dir, 0700); err != nil {
		return "", err
	}

	path := filepath.Join(dir, "fillin.json")
	rfile, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0600)
	if err != nil {
		return "", err
	}
	w := new(bytes.Buffer)
	filled, err := Fillin(args, rfile, w, in, out)
	if err != nil {
		return "", err
	}
	cmd := escapeJoin(filled)
	rfile.Close() // not be defered due to rename

	tmp := filepath.Join(dir, fmt.Sprintf("fillin.%d-%d.json", os.Getpid(), rand.Int()))
	defer os.Remove(tmp)
	wfile, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return "", err
	}
	if n, err := wfile.Write(w.Bytes()); n != w.Len() || err != nil {
		wfile.Close()
		return "", err
	}
	wfile.Close() // not be defered due to rename

	if err := os.Rename(tmp, path); err != nil {
		return "", err
	}

	if cmd != "" {
		histfile := filepath.Join(dir, ".fillin.histfile")
		hfile, err := os.OpenFile(histfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
		defer hfile.Close()
		if err != nil {
			return "", err
		}
		hfile.WriteString(fmt.Sprintf(": %d:0;%s\n", time.Now().Unix(), cmd))
		hfile.Chmod(0600)
	}

	return cmd, nil
}

func normalizeFileSeparator(path string) string {
	p := filepath.Join(strings.Split(path, "/")...)
	if path[0] == '/' {
		p = "/" + p
	}
	return p
}

func escapeJoin(args []string) string {
	if len(args) == 1 {
		return args[0]
	}
	for i, arg := range args {
		args[i] = escape(arg)
	}
	return strings.Join(args, " ")
}

func escape(arg string) string {
	switch arg {
	case "|", "||", "&&", ">", ">>", "<":
		return arg
	}
quote:
	for _, quote := range []bool{false, true} {
		s := arg
		isHead, afterFd, afterRedirect := true, false, false
		var buf bytes.Buffer
		if quote {
			buf.WriteByte('\'')
		}
		for len(s) > 0 {
			c, l := utf8.DecodeRuneInString(s)
			s = s[l:]
			if (isHead || afterFd) && strings.ContainsRune("<>", c) && !quote {
				buf.WriteRune(c)
				isHead, afterFd, afterRedirect = false, false, true
				continue
			} else if afterRedirect && strings.ContainsRune("<>", c) && !quote {
				buf.WriteRune(c)
			} else if isHead && strings.ContainsRune("12", c) && !quote {
				buf.WriteRune(c)
				isHead, afterFd = false, true
				continue
			} else if afterRedirect && strings.ContainsRune("&", c) && !quote {
				buf.WriteRune(c)
			} else if strings.ContainsRune("\\'\"`${[|&;<>()*?!", c) && !quote {
				buf.WriteByte('\\')
				buf.WriteRune(c)
			} else if c == rune(' ') && !quote {
				continue quote
			} else if c == rune('\t') {
				if quote {
					buf.WriteByte('\\')
					buf.WriteByte('t')
				} else {
					continue quote
				}
			} else {
				buf.WriteRune(c)
			}
			isHead, afterFd, afterRedirect = false, false, false
		}
		if quote {
			buf.WriteByte('\'')
		}
		return buf.String()
	}
	return ""
}
