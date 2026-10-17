// provenance checks that nothing in this repository is copied from dreddsa5dies/automateGo,
// the project some early versions of the chapter 7 to 11 code were adapted from. It
// downloads the upstream tree, compares every file by git blob hash, and compares every
// Go file line by line after dropping blank lines, comments and boilerplate.
//
//	go run ./tools/provenance            # exit 1 on any violation
//	go run ./tools/provenance -v         # also print the nearest upstream file per Go file
package main

import (
	"archive/tar"
	"bufio"
	"compress/gzip"
	"crypto/sha1"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	upstreamTarball = "https://codeload.github.com/dreddsa5dies/automateGo/tar.gz/HEAD"
	maxOverlap      = 0.30 // share of our meaningful lines allowed to appear in one upstream file
	maxRun          = 4    // longest allowed run of identical consecutive meaningful lines
)

type file struct {
	path  string
	blob  string   // git blob sha1
	lines []string // meaningful lines (Go files only)
}

// boilerplate is syntax every Go program shares; it carries no authorship signal. Lines
// shorter than minMeaningful characters (`}()`, `wg.Wait()`, `close(ch)`) are dropped for
// the same reason.
var boilerplate = regexp.MustCompile(`^(package \w+|import \(?|\)|\}|\{|\} else \{|\}\(\)|return( err| nil| err, nil|( nil,)+ err)?|if err != nil \{|defer \w+\.Close\(\)|defer wg\.Done\(\)|var wg sync\.WaitGroup|go func\(\) \{|func main\(\) \{|os\.Exit\(1\)|log\.Fatal(ln|f)?\(err\)|"[\w./-]+"|\w+ "[\w./-]+"|fmt\.Fprintln\(os\.Stderr, "error:", err\)|if err := run\(os\.Args\[1:\]\); err != nil \{|if err := run\(\); err != nil \{)$`)

const minMeaningful = 12

func meaningful(src []byte) []string {
	var out []string
	sc := bufio.NewScanner(strings.NewReader(string(src)))
	for sc.Scan() {
		l := strings.TrimSpace(sc.Text())
		if l == "" || len(l) < minMeaningful || strings.HasPrefix(l, "//") || boilerplate.MatchString(l) {
			continue
		}
		out = append(out, l)
	}
	return out
}

func blobSHA(b []byte) string {
	h := sha1.New()
	fmt.Fprintf(h, "blob %d\x00", len(b))
	h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func upstream() ([]file, error) {
	client := &http.Client{Timeout: 2 * time.Minute}
	resp, err := client.Get(upstreamTarball)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download upstream: %s", resp.Status)
	}
	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		return nil, err
	}
	tr := tar.NewReader(gz)
	var files []file
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if hdr.Typeflag != tar.TypeReg || strings.Contains(hdr.Name, "/vendor/") {
			continue
		}
		b, err := io.ReadAll(tr)
		if err != nil {
			return nil, err
		}
		f := file{path: hdr.Name[strings.Index(hdr.Name, "/")+1:], blob: blobSHA(b)}
		if strings.HasSuffix(hdr.Name, ".go") {
			f.lines = meaningful(b)
		}
		files = append(files, f)
	}
	return files, nil
}

func ours(root string) ([]file, error) {
	var files []file
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		name := d.Name()
		if d.IsDir() {
			if name == ".git" || name == "tools" || (path != root && strings.HasPrefix(name, ".")) {
				return filepath.SkipDir
			}
			return nil
		}
		if name == "go.sum" || name == "LICENSE" {
			return nil
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(root, path)
		f := file{path: rel, blob: blobSHA(b)}
		if strings.HasSuffix(name, ".go") {
			f.lines = meaningful(b)
		}
		files = append(files, f)
		return nil
	})
	return files, err
}

// overlap returns the share of a's lines that also occur in b, and the longest run of
// consecutive lines of a that appears verbatim and in order in b.
func overlap(a, b []string) (share float64, run int) {
	if len(a) == 0 {
		return 0, 0
	}
	set := map[string]bool{}
	for _, l := range b {
		set[l] = true
	}
	hits := 0
	for _, l := range a {
		if set[l] {
			hits++
		}
	}
	// longest common substring over lines
	prev := make([]int, len(b)+1)
	for i := 1; i <= len(a); i++ {
		cur := make([]int, len(b)+1)
		for j := 1; j <= len(b); j++ {
			if a[i-1] == b[j-1] {
				cur[j] = prev[j-1] + 1
				run = max(run, cur[j])
			}
		}
		prev = cur
	}
	return float64(hits) / float64(len(a)), run
}

func main() {
	verbose := flag.Bool("v", false, "print the nearest upstream file for every Go file")
	root := flag.String("root", ".", "repository root")
	flag.Parse()

	up, err := upstream()
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	mine, err := ours(*root)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	byBlob := map[string]string{}
	for _, f := range up {
		byBlob[f.blob] = f.path
	}

	violations := 0
	for _, f := range mine {
		if p, ok := byBlob[f.blob]; ok {
			fmt.Printf("IDENTICAL  %s == upstream %s\n", f.path, p)
			violations++
		}
	}
	type near struct {
		path  string
		share float64
		run   int
	}
	var report []string
	for _, f := range mine {
		if f.lines == nil {
			continue
		}
		best := near{}
		for _, u := range up {
			if u.lines == nil {
				continue
			}
			share, run := overlap(f.lines, u.lines)
			if share > best.share || (share == best.share && run > best.run) {
				best = near{u.path, share, run}
			}
		}
		bad := best.share > maxOverlap || best.run > maxRun
		if bad {
			violations++
		}
		if bad || *verbose {
			report = append(report, fmt.Sprintf("%-9s %-55s share=%.2f run=%d  nearest %s", map[bool]string{true: "DERIVED", false: "ok"}[bad], f.path, best.share, best.run, best.path))
		}
	}
	sort.Strings(report)
	for _, r := range report {
		fmt.Println(r)
	}
	fmt.Printf("%d file(s) checked against %d upstream file(s); %d violation(s)\n", len(mine), len(up), violations)
	if violations > 0 {
		os.Exit(1)
	}
}
