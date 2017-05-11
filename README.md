# gener8
simple golang go:generate code generator

# install
```
go get github.com/dbreedt/gener8
```

# usage
Create a file with markers
 * $pkg
 * $kw1
 * $kw2
 * $kw3
 * $kw4

Then add a go:generate comment to your codebase

```
//go:generate gener8 -in=make.me -out=made.go -pkg=made -kw1=test
```

or just run it from the cmdline: 
```
gener8 -in=make.me -out=made.go -pkg=made -kw1=test
```
