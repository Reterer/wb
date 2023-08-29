package main

import (
	"bufio"
	"bytes"
	ciclebuf "dev05/cicle_buf"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// grep [OPTION...] PATTERN [FILE...]

type config struct {
	after      int
	before     int
	context    int
	count      bool
	ignoreCase bool
	invert     bool
	fixed      bool
	lineNum    bool

	pattern string
	files   []string
}

func parseArgs() *config {
	var cfg config
	flag.IntVar(&cfg.after, "A", 0, "печатать +N строк после совпадения")
	flag.IntVar(&cfg.before, "B", 0, "печатать +N строк до совпадения")
	flag.IntVar(&cfg.context, "C", 0, "печатать ±N строк вокруг совпадения")
	flag.BoolVar(&cfg.count, "c", false, "количество строк")
	flag.BoolVar(&cfg.ignoreCase, "i", false, "игнорировать регистр")
	flag.BoolVar(&cfg.invert, "v", false, "вместо совпадения, исключать")
	flag.BoolVar(&cfg.fixed, "F", false, "очное совпадение со строкой, не паттерн")
	flag.BoolVar(&cfg.lineNum, "n", false, "печатать номер строки")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	cfg.pattern = args[0]
	cfg.files = args[1:]

	return &cfg
}

type grepOpts struct {
	count      bool
	ignoreCase bool
	invert     bool
	fixed      bool
	lineNum    bool
	after      int
	before     int
}

type Grep struct {
	opts    grepOpts
	pattern string
	files   []string
}

type grepPrinter struct {
	w             io.Writer
	matchPrefix   string
	contextPrefix string
	matchSep      byte
	contextSep    byte

	needNumerateLines bool
	wasPritned        bool // Первая строка никогда не начинается с \n
	// Правильней было бы сделать циклический буффер для любого типа
}

func newGrepPrinter(w io.Writer, filename string, msep byte, csep byte, needNum bool) grepPrinter {
	res := grepPrinter{
		w:                 w,
		matchSep:          msep,
		contextSep:        csep,
		needNumerateLines: needNum,
		matchPrefix:       "\n",
		contextPrefix:     "\n",
	}

	if filename != "" {
		res.matchPrefix = "\n" + filename + string(msep)
		res.contextPrefix = "\n" + filename + string(csep)
	}

	return res
}

func (p *grepPrinter) formatLine(nfilenamesep, line string, sep byte, num int) []byte {
	var buf bytes.Buffer
	if p.needNumerateLines {
		snum := strconv.Itoa(num)
		buf.Grow(len(nfilenamesep) + len(snum) + len(line) + 1)
		buf.WriteString(nfilenamesep)
		buf.WriteString(snum)
		buf.WriteByte(sep)
		buf.WriteString(line)
	} else {
		buf.Grow(len(nfilenamesep) + len(line))
		buf.WriteString(nfilenamesep)
		buf.WriteString(line)
	}
	return buf.Bytes()
}

func (p *grepPrinter) formatMatched(line string, num int) []byte {
	return p.formatLine(p.matchPrefix, line, p.matchSep, num)
}

func (p *grepPrinter) formatContext(line string, num int) []byte {
	return p.formatLine(p.contextPrefix, line, p.contextSep, num)
}

func (p *grepPrinter) print(line []byte) error {
	if !p.wasPritned {
		if len(line) > 0 {
			line = line[1:]
		}
		p.wasPritned = true
	}

	if _, err := p.w.Write(line); err != nil {
		return err
	}
	return nil
}

func (p *grepPrinter) printMatched(line string, num int) error {
	b := p.formatMatched(line, num)
	return p.print(b)
}
func (p *grepPrinter) printContext(line string, num int) error {
	b := p.formatContext(line, num)
	return p.print(b)
}

func (g *Grep) grepio(filename string, reader io.Reader, writer io.Writer) error {
	scn := bufio.NewScanner(reader)

	// Создаем регулярное выражение -F; -i
	pattern := g.pattern
	if g.opts.fixed {
		pattern = "^" + regexp.QuoteMeta(pattern) + "$"
	}
	if g.opts.ignoreCase {
		pattern = "(?i)" + pattern
	}
	mexp, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}

	// ---
	printer := newGrepPrinter(writer, filename, ':', '-', g.opts.lineNum)

	count := 0      // количество найденных строк (для -c)
	lineNumber := 0 // количество строк (для -n)

	endAfter := 0                               // -A
	beforeBuf := ciclebuf.NewBuf(g.opts.before) // -B

	for scn.Scan() {
		line := scn.Text()
		lineNumber++

		// Match cond (-v)
		match := mexp.MatchString(line)
		if g.opts.invert {
			match = !match
		}
		// ---

		if match {
			if !g.opts.count { // Если нужно выводить результат
				{
					// Вывод строк before, если нужно: -B
					line, ok := beforeBuf.Pop()
					for ok {
						if err := printer.print(line); err != nil {
							return err
						}
						line, ok = beforeBuf.Pop()
					}
				}

				if err := printer.printMatched(line, lineNumber); err != nil {
					return err
				}
				endAfter = lineNumber + g.opts.after // -A
			}
			count++
		} else {
			if !g.opts.count {
				if lineNumber <= endAfter {
					if err := printer.printContext(line, lineNumber); err != nil {
						return err
					}
				} else {
					// Если мы не вывели эту строку, то ее нужно добавить в буффер, если он нужен
					if g.opts.before > 0 {
						// -B
						outbuf := printer.formatContext(line, lineNumber)
						beforeBuf.Push(outbuf)
					}
				}
			}
		}
	}

	if g.opts.count { // -c
		if _, err := writer.Write([]byte(printer.matchPrefix[1:] + strconv.Itoa(count))); err != nil {
			return err
		}
	}

	return nil
}

func NewGrep(cfg *config) (*Grep, error) {
	files := make([]string, len(cfg.files))
	copy(files, cfg.files)

	after := cfg.after
	before := cfg.before
	if cfg.context != 0 {
		after = cfg.context
		before = cfg.context
	}

	return &Grep{
		opts: grepOpts{
			count:      cfg.count,
			ignoreCase: cfg.ignoreCase,
			invert:     cfg.invert,
			fixed:      cfg.fixed,
			lineNum:    cfg.lineNum,
			after:      after,
			before:     before,
		},
		pattern: cfg.pattern,
		files:   files,
	}, nil
}

func (g *Grep) Start() {
	if len(g.files) == 0 {
		g.files = append(g.files, "-")
	}
	for _, file := range g.files {
		if file == "-" {
			err := g.grepio("", os.Stdin, os.Stdout)
			if err != nil {
				fmt.Fprintf(os.Stderr, "grep: %s: %s", "(standard input)\n", err)
			}
		} else {
			f, err := os.Open(file)
			if err != nil {
				fmt.Fprintf(os.Stderr, "grep: %s: %s\n", file, err)
				continue
			}
			err = g.grepio(file, f, os.Stdout)
			if err != nil {
				fmt.Fprintf(os.Stderr, "grep: %s: %s\n", file, err)
			}
			f.Close()
		}
	}
}

func main() {
	cfg := parseArgs()
	grep, _ := NewGrep(cfg)
	grep.Start()
}
