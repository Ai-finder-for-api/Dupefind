// dupefind: find duplicate files under a directory by content hash.
package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/Ai-finder-for-api/test-proj-1/internal/dupe"
)

func main() {
	minSize := flag.Int64("min", 1, "ignore files smaller than this many bytes")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: dupefind [-min bytes] <dir> [dir...]\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}
	groups, errs := dupe.Find(roots, *minSize)
	for _, e := range errs {
		fmt.Fprintln(os.Stderr, "warn:", e)
	}
	sort.Slice(groups, func(i, j int) bool { return groups[i].Size > groups[j].Size })
	wasted := int64(0)
	for _, g := range groups {
		fmt.Printf("%d copies of %d bytes (sha256 %s):\n", len(g.Paths), g.Size, g.Hash[:16])
		for _, p := range g.Paths {
			fmt.Println("  ", p)
		}
		wasted += g.Size * int64(len(g.Paths)-1)
	}
	fmt.Printf("%d duplicate group(s), %d bytes reclaimable\n", len(groups), wasted)
}
