# Red Black Tree

Go generic implementation of Red Black tree.

## Notes

`New()` Takes a compare function so that the value stored
can be more complicated than only things capable of being
compared with `cmp.Ordered`. Notably, for any type in
`cmp.Ordered` you can construct the tree with
`cmp.Compare` as the comparison function.

```go
t := rbtree.New[int](cmp.Compare)
```

The API of the Red Black tree exposes `Node` and
`Node.Value` so that callers can use nodes themselves,
such as keeping references to the internal data
structure for additional indexing.

## Iterators

Compile with Go 1.22 and `iterator`, and the GOEXPERIMENT=rangefunc build tag to have iterator support.
