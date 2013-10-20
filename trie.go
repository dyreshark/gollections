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
	isEnd    bool
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
	child := trieNode{
		children: map[rune]*trieNode{},
		value:    utf8.RuneError,
	}

	return &Trie{
		root: child,
	}
}

// Returns the trieNode of the last char in the given string, and its parent.
// If not found (or utf8 decode error), nil is returned for both.
func (t *Trie) searchNode(s string) *trieNode {
	if len(s) == 0 {
		return nil
	}

	current := &t.root
	for len(s) != 0 && current != nil {
		r, size := utf8.DecodeRuneInString(s)
		if r == utf8.RuneError {
			// TODO: Maybe return error instead?
			return nil
		}
		var ok bool
		current, ok = current.children[r]
		if !ok {
			return nil
		}
		s = s[size:]
	}
	return current
}

// Searches for the given string in the trie.
//
// Returns true on found, false on not found (or error decoding string)
func (t *Trie) Has(s string) bool {
	res := t.searchNode(s)
	return res != nil && res.isEnd
}

// Searches for the given string in the trie. This will return true if
// there is either a full string or just the prefix of a string that
// matches the input.
//
// Returns true on found, false on not found (or error decoding string).
func (t *Trie) HasPrefix(s string) bool {
	return t.searchNode(s) != nil
}

func (t *Trie) Delete(s string) {
	if len(s) == 0 {
		return
	}

	// The tradeoff here is to either give each trieNode
	// a parent pointer and use searchNode, or to just memoize
	// the last parent that was marked as !isEnd. More code duplication,
	// but I'd rather that than use extra storage for each node.
	current := &t.root

	var lastNeededRune rune
	var lastNeededNode *trieNode
	for len(s) != 0 {
		var size int
		r, size := utf8.DecodeRuneInString(s)
		if r == utf8.RuneError {
			// TODO: Maybe report error
			return
		}

		// 'Needed' is either the end of a word, or a node that
		// has to support more than 1 child.
		if current.isEnd || len(current.children) > 1 {
			lastNeededRune = r
			lastNeededNode = current
		}

		var ok bool
		current, ok = current.children[r]
		if !ok {
			return
		}

		s = s[size:]
	}

	if len(current.children) != 0 {
		current.isEnd = false
	} else if lastNeededNode == nil {
		// Even root wasn't needed? Sweet. Because this is a special
		// case, it's handled (admittedly) somewhat stupidly.
		if len(t.root.children) != 1 {
			panic("Internal error: t.root.children has length != 1")
		}
		for k, _ := range t.root.children {
			delete(t.root.children, k)
		}
	} else {
		// Nothing depends on current. Delete every node that
		// only current depends on.
		delete(lastNeededNode.children, lastNeededRune)
	}
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

// Implementation of `Put`, but assuming everything in the string
// is valid.
func (t *trieNode) putValid(s string) *trieNode {
	if len(s) == 0 {
		t.isEnd = true
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

	node := &t.root
	for len(s) != 0 {
		r, size := utf8.DecodeRuneInString(s)
		node = node.addChildNode(r)
		s = s[size:]
	}
	node.isEnd = true

	return nil
}
