# gener8
simple golang go:generate code generator

# install
```
go get github.com/dbreedt/gener8
```

# usage
Create a file with markers
 * $pkg
 * $kwn where n is a numeric value starting from 1

Then add a go:generate comment to your codebase

```
//go:generate gener8 -in=make.me -out=made.go -pkg=made -kws=test,tPtr
```

or just run it from the cmdline: 
```
gener8 -in=make.me -out=made.go -pkg=made -kws=test,tPtr
```

In these examples `$kw1` will be replaced with `test` and `$kw2` will be replaced with `tPtr`
