# GoMoji
<p align="center">work with emoji in the most convenient way</p>

GoMoji is a Go package that provides a [fast](#performance) and [simple](#check-string-contains-emoji) way to work with emojis in a string.
It has features such as:
 * [check whether string contains emoji](#check-string-contains-emoji),
 * [find all emojis in string](#find-all),
 * [get all emojis](#get-all), 

Getting Started
===============

## Installing

To start using GoMoji, install Go and run `go get`:

```sh
$ go get -u github.com/forPelevin/gomoji
```

This will retrieve the package.

## Check string contains emoji
```go
package main

import (
    "github.com/forPelevin/gomoji"
)

func main() {
    res := gomoji.ContainsEmoji("hello world")
    println(res) // false
    
    res = gomoji.ContainsEmoji("hello world 🤗")
    println(res) // true
}
```

## Find all
The function searches for all emoji occurrences in a string. It returns a nil slice if there are no emojis.
```go
package main

import (
    "github.com/forPelevin/gomoji"
)

func main() {
    res := gomoji.FindAll("🧖 hello 🦋 world")
    println(res)
}
```

## Get all
The function returns all existed emojis. You can do whatever you need with the list.
 ```go
 package main
 
 import (
     "github.com/forPelevin/gomoji"
 )
 
 func main() {
     emojis := gomoji.AllEmojis()
     println(emojis)
 }
 ```

## Emoji entity
All searching methods return the Emoji entity which contains comprehensive info about emoji.
```go
type Emoji struct {
    Slug        string `json:"slug"`
    Character   string `json:"character"`
    UnicodeName string `json:"unicode_name"`
    CodePoint   string `json:"code_point"`
    Group       string `json:"group"`
    SubGroup    string `json:"sub_group"`
}
 ```
Example:
```go
[]gomoji.Emoji{
    {
        Slug:        "e3-0-butterfly",
        Character:   "🦋",
        UnicodeName: "E3.0 butterfly",
        CodePoint:   "1F98B",
        Group:       "animals-nature",
        SubGroup:    "animal-bug",
    },
    {
        Slug:        "roll-of-paper",
        Character:   "🧻",
        UnicodeName: "roll of paper",
        CodePoint:   "1F9FB",
        Group:       "objects",
        SubGroup:    "household",
    },
}
 ```

## Performance

Benchmarks of GoMoji

```
BenchmarkContainsEmojiParallel-8   	94079461	        13.1 ns/op	       0 B/op	       0 allocs/op
BenchmarkContainsEmoji-8           	23728635	        49.8 ns/op	       0 B/op	       0 allocs/op
BenchmarkFindAllParallel-8         	10220854	       115 ns/op	     288 B/op	       2 allocs/op
BenchmarkFindAll-8                 	 4023626	       294 ns/op	     288 B/op	       2 allocs/op
```

## Contact
Vlad Gukasov [@vgukasov](https://www.facebook.com/vgukasov)

## License

GoMoji source code is available under the MIT [License](/LICENSE).