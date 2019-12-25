This is a bad tool and you shouldn't use it.


This is a multihash-ified file checker which could be used in the same way you might use md5sum or
sha256sum commands.

installing in $GOBIN:

```shell
$ cd /to/this/directory
$ go get
```

list available hash algorithms:

```shell
$ multihash list
...
```

Generating multihash sums:

```shell
$ multihash sum <filename>
```

or from stdin
```shell
$ somecommand | multihash sum
```

Select a different hashing algorithm by passing `-h`.

```shell
multihash sum -h sha3-224 <filename>
```

The output of the sum command can be captured to verfiy files later

```shell
$ echo A > a.txt
$ echo B > b.txt
$ multihash sum *.txt > SUMS
$ multihash check  SUMS
a.txt  PASS
b.txt  PASS

$ echo CORRUPTION > b.txt
a.txt  PASS
b.txt  FAIL
```
