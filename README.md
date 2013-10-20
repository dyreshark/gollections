Gollections
==========

Random data structures/algorithms that are implemented in Go. At the moment,
all that I have is a Trie. I plan to add more as time goes on.

Trie
----------

Usage consists of five functions:
      NewTrie
      *Trie.Put
      *Trie.Delete
      *Trie.Has
      *Trie.HasPrefix

      trie := gollections.NewTrie()
      trie.Put("FooBarBaz")
      trie.Has("FooBarBaz")       // true
      trie.HasPrefix("FooBarBaz") // true
      trie.Has("FooBa")           // false
      trie.HasPrefix("FooBa")     // true
      trie.Has("foobarbaz")       // false
      trie.HasPrefix("foobarbaz") // false
      trie.Put("FooBar")
      trie.Delete("FooBarBaz")
      trie.HasPrefix("FooBarBa")  // false
      trie.Has("FooBar")          // true
      trie.Has("FooBarBaz")       // false

License
----------

Licensed under Apache v2.
