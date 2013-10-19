package collections

import (
	"errors"
	"unicode/utf8"
)

// The root and elements of a trie.
//
// Each TrieNode is associated with a rune. For example:
//      a
//     / \
//    b   c
//   / \  |
//  d   e f
//
// In this case, a, b, c, d, e, and f are all TrieNodes.
type trieNode struct {
	children map[rune]*trieNode
	value    rune
}

type Trie struct {
	// I want distinct types for Trie and trieNode. And admittedly
	// have no clue how to cast from (type Integer int) *Integer ->
	// *int
	root trieNode
}

// Makes a trie node for me.
func newTrieNode(r rune) *trieNode {
	return &trieNode{
		children: map[rune]*trieNode{},
		value:    r,
	}
}

// Creates a new Trie for the user
//
// Never returns nil.
func NewTrie() *Trie {
	// TODO: Maybe create children on demand?
	node := *newTrieNode(utf8.RuneError)
	return &Trie{
		root: node,
	}
}

// Searches the given trieNode to a rune.
//
// Returns nil, false on not found.
// Returns non-nil, true on found.
func (t *Trie) SearchRune(r rune) bool {
	_, ok := t.root.children[r]
	return ok
}

// Searches for the given string in the trie.
//
// Returns true on found, false on not found (or error decoding string)
func (t *Trie) Search(s string) bool {
	current := &t.root
	for len(s) != 0 && current != nil {
		r, size := utf8.DecodeRuneInString(s)
		if r == utf8.RuneError {
			// TODO: Maybe return error instead?
			return false
		}
		var ok bool
		current, ok = current.children[r]
		if !ok {
			return false
		}
		s = s[size:]
	}
	return current != nil
}

// Adds a child node and returns the trieNode that 'represents' it.
func (t *trieNode) addChildNode(r rune) *trieNode {
	node, ok := t.children[r]
	if !ok {
		node = newTrieNode(r)
		t.children[r] = node
	}
	return node
}

// Puts a rune in the Trie.
func (t *Trie) PutRune(r rune) {
	t.root.addChildNode(r)
}

// Implementation of `Put`, but assuming everything in the string
// is valid.
func (t *trieNode) putValid(s string) *trieNode {
	if len(s) == 0 {
		return t
	}

	r, size := utf8.DecodeRuneInString(s)
	next := t.addChildNode(r)
	return next.putValid(s[size:])
}

// Puts a full string of runes into the given Trie.
//
// Returns the terminating trieNode and a nil error on success,
// returns nil and an error on failure. Currently, failure only
// happens if s has an invalid utf-8 sequence in it.
func (t *Trie) Put(s string) error {
	if !utf8.ValidString(s) {
		return errors.New("Invalid utf8 in string")
	}

	// TODO: It might be worthwhile to make undos possible, so we can
	// not walk the string twice for this.

	t.root.putValid(s)
	return nil
}
