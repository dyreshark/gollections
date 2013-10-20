/*
  Copyright 2013 George Burgess IV

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package gollections

import (
	"math/rand"
	"testing"
)

// Indirectly tests PutRune too.
func TestTrieHasAndHasPrefix(t *testing.T) {
	root := NewTrie()

	putIn := []string{"abc", "de", "fghi", "acl", "mlp"}
	notPutIn := []string{"abd", "df", "fhgij", "adl", "fim"}
	prefixes := []string{}

	// Populate implicitlyPutIn
	for _, n := range putIn {
		for i := len(n) - 1; i > 0; i-- {
			prefixes = append(prefixes, n[:i])
		}
	}

	for _, n := range putIn {
		root.Put(n)
		if !root.Has(n) {
			t.Fatal("Expected to find string", n)
		}
		if !root.HasPrefix(n) {
			t.Fatal("Expected to find prefix", n)
		}
	}

	for _, n := range prefixes {
		if root.Has(n) {
			t.Fatal("Expected to not find implicit string through `Has`", n)
		}

		if !root.HasPrefix(n) {
			t.Fatal("Expected to find implicit string through `HasPrefix`", n)
		}
	}

	for _, n := range notPutIn {
		if root.Has(n) {
			t.Fatal("Expected not to find string", n)
		}
		if root.HasPrefix(n) {
			t.Fatal("Expected not to find prefix", n)
		}
	}
}

func TestTrieDelete(t *testing.T) {
	trie := NewTrie()

	// Case one: Just call delete on it; Has and HasPrefix
	// fail on caseOne.
	caseOne := "hello"
	trie.Put(caseOne)

	if !trie.Has(caseOne) {
		t.Fatal("Expected to find", caseOne)
	}

	trie.Delete(caseOne)
	if trie.Has(caseOne) || trie.HasPrefix(caseOne) || trie.HasPrefix(caseOne[:1]) {
		t.Fatal("Expected trie to not have", caseOne)
	}

	// Case two: Delete should *not* delete the node
	// if the string would otherwise be a prefix.
	trie = NewTrie()
	caseTwo := "foo"
	caseTwoPlus := "fooo"

	trie.Put(caseTwoPlus)
	if !trie.Has(caseTwoPlus) {
		t.Fatal("Expected to find", caseTwoPlus)
	}

	if !trie.HasPrefix(caseTwo) {
		t.Fatal("Expected to find prefix", caseTwo)
	}

	trie.Delete(caseTwo)
	if trie.Has(caseTwo) {
		t.Fatal("Expected not to find", caseTwo, "after delete")
	}
	if !trie.Has(caseTwoPlus) {
		t.Fatal("Expected to find", caseTwoPlus, "after delete")
	}

	if !trie.HasPrefix(caseTwo) {
		t.Fatal("Expected to find prefix", caseTwo, "after delete")
	}

	// Case three: Delete should delete *all* nodes up until
	// the most recent one that depends on other strings
	trie = NewTrie()
	caseThree := "foooooo"
	caseThreeMinus := "fooo"

	trie.Put(caseThree)
	trie.Put(caseThreeMinus)
	if !trie.Has(caseThree) || !trie.Has(caseThreeMinus) {
		t.Fatal("Expected to find both", caseThree, "and", caseThreeMinus)
	}

	trie.Delete(caseThree)
	if trie.Has(caseThree) {
		t.Fatal("Expected to not find", caseThree, "after delete")
	}

	if !trie.Has(caseThreeMinus) {
		t.Fatal("Expected to find", caseThree, "after delete")
	}

	if trie.HasPrefix(caseThree[len(caseThreeMinus)+1:]) {
		t.Fatal("Not all caseThree remnants deleted")
	}
}

// --------- Here be benchmarks ------------

func BenchmarkTrieHas(b *testing.B) {
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
		if root.Has(s.s) != s.ok {
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
		if ok := trie.Has(s); ok != expected {
			b.Fatalf("Unexpected result for string %d (%s)", i, s)
		}
	}
}
