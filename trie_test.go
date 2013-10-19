package collections

import (
	"math/rand"
	"testing"
)

// Indirectly tests PutRune too for free
func TestTrieSearchRune(t *testing.T) {
	root := NewTrie()

	putIn := []rune{'h', 'i'}
	testFail := []rune{'w', 'a', 't'}

	for _, n := range putIn {
		root.PutRune(n)
		if ok := root.SearchRune(n); !ok {
			t.Fatal("Didn't find node")
		}
	}

	for _, n := range testFail {
		if ok := root.SearchRune(n); ok {
			t.Fatal("Not expecting to find rune. Found it.")
		}
	}
}

// Indirectly tests PutRune too.
func TestTrieSearch(t *testing.T) {
	root := NewTrie()

	putIn := []string{"abc", "de", "fghi", "acl", "mlp"}
	notPutIn := []string{"abd", "df", "fhgij", "adl", "fim"}
	implicitlyPutIn := []string{}

	// Populate implicitlyPutIn
	for _, n := range putIn {
		for i := len(n); i > 0; i-- {
			implicitlyPutIn = append(implicitlyPutIn, n[:i])
		}
	}

	for _, n := range putIn {
		root.Put(n)
		if ok := root.Search(n); !ok {
			t.Fatal("Expected to find string", n)
		}
	}

	for _, n := range implicitlyPutIn {
		if ok := root.Search(n); !ok {
			t.Fatal("Expected to find implicit string", n)
		}
	}

	for _, n := range notPutIn {
		if ok := root.Search(n); ok {
			t.Fatal("Expected not to find string", n)
		}
	}
}

// --------- Here be benchmarks ------------

// Admittedly, this is just here to test memory isn't alloced for some
// obscure reason.
func BenchmarkTrieSearchRune(b *testing.B) {
	root := NewTrie()
	runes := []struct {
		r  rune
		ok bool
	}{{'a', false}, {'b', true}, {'c', false}, {'d', true}}

	for _, r := range runes {
		if r.ok {
			root.PutRune(r.r)
		}
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := runes[i%len(runes)]
		if root.SearchRune(s.r) != s.ok {
			b.Fatal("Got unexpected result for", s.r)
		}
	}
}

func BenchmarkTrieSearch(b *testing.B) {
	root := NewTrie()

	strings := []struct {
		s  string
		ok bool
	}{
		{"hi", false},
		{"ha", true},
		{"whatdidyousay", false},
		{"whatdidyousai", true},
		{"Invariant", false},
		{"invariant", true},
		{"some other string", false},
		{"some other strin", true},
	}

	for _, r := range strings {
		if r.ok {
			root.Put(r.s)
		}
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := strings[i%len(strings)]
		if root.Search(s.s) != s.ok {
			b.Fatal("Got unexpected result for", s.s)
		}
	}
}

func BenchmarkLargeTrieSearch(b *testing.B) {
	const NUM_STRINGS = 1000000
	const STR_LEN = 20
	// Number of possible chars our strings can have
	const NUM_CHRS = 94
	// Offset of the char values
	const OFFSET = 32

	rand.Seed(0) // Arbitrary seed

	trie := NewTrie()
	strings := make([]string, NUM_STRINGS)
	buf := make([]rune, STR_LEN)

	// There's a chance of two strings being identical. If one is marked
	// as "not in trie" and the other is marked as "in trie", we'll get
	// incorrect output. Need a way to quickly see if we've used a string before.
	// Don't want this benchmark to be polynomial time.
	stringSet := make(map[string]bool)

	// Make length-10 strings
	// Side-note: Even-indexed strings are marked as not in the trie.
	for i := 0; i < NUM_STRINGS; i++ {
		for x := 0; x < STR_LEN; x++ {
			buf[x] = rune(rand.Int31n(NUM_CHRS) + OFFSET)
		}
		s := string(buf)
		if len(s) != STR_LEN {
			b.Fatal("Unexpected string size:", len(s))
		}
		// If randomness has cursed us, try again
		if _, ok := stringSet[s]; ok {
			i--
			continue
		}
		strings[i] = s
		stringSet[s] = true
		if i%2 != 0 {
			trie.Put(s)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := strings[i%len(strings)]
		expected := i%2 != 0
		if ok := trie.Search(s); ok != expected {
			b.Fatalf("Unexpected result for string %d (%s)", i, s)
		}
	}
}
