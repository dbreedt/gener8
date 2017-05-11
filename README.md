# gener8
simple golang go:generate code generator

# usage
Create a file with markers
 * $pkg
 * $kw1
 * $kw2
 * $kw3

Then add a go:generate comment to your codebase

```
//go:generate gener8 -in=make.me -out=made.go -pkg=made -kw1=test
```

or just run it from the cmdline `gener8 -in=make.me -out=made.go -pkg=made -kw1=test`
