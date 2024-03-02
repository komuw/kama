# Release Notes

Most recent version is listed first.  

# v0.0.24
- Add example: https://github.com/komuw/kama/pull/77

# v0.0.23
- Make stacktrace colorization configurable: https://github.com/komuw/kama/pull/76

# v0.0.22
- Bugfix

# v0.0.21
- Bugfix

# v0.0.20
- Support older Go versions: https://github.com/komuw/kama/pull/71
- Stackp can now write to any io.writer: https://github.com/komuw/kama/pull/74
- Display stacktrace when dumping variables: https://github.com/komuw/kama/pull/73
                                             https://github.com/komuw/kama/pull/75

# v0.0.19
- Update docs: https://github.com/komuw/kama/commit/3dfdef76fbdd55a8ebed2d01dfb1cbb7c1cf2fc2

# v0.0.18
- Add ability to diff items: https://github.com/komuw/kama/pull/67
- Add legend to stack trace: https://github.com/komuw/kama/pull/68

# v0.0.17
- Add ability to dump items with circular references: https://github.com/komuw/kama/pull/66

# v0.0.16
- Update Go version: https://github.com/komuw/kama/pull/62
- Add ability to dump private struct fields via config option: https://github.com/komuw/kama/pull/64

# v0.0.15
- Retracts v0.0.14 which had a panic bug.

# v0.0.14
- Add config option for how many slices/maps to display: https://github.com/komuw/kama/pull/58

# v0.0.13
- Update Go version: https://github.com/komuw/kama/pull/56

# v0.0.12
- Pretty print errors: https://github.com/komuw/kama/pull/54
- Handle more interfaces: https://github.com/komuw/kama/pull/55

# v0.0.11
- Fix formatting of custom string types: https://github.com/komuw/kama/pull/52

# v0.0.10
- Implement better indentation for maps and slices: https://github.com/komuw/kama/pull/51

# v0.0.9
- Use testdata in more testcases: https://github.com/komuw/kama/pull/46
- Implement better indentation: https://github.com/komuw/kama/pull/47
- Pretty print context.Context: https://github.com/komuw/kama/pull/48
- Fix Static analysis: https://github.com/komuw/kama/pull/49

# v0.0.8
- Zero value struct should not indent: https://github.com/komuw/kama/pull/44

# v0.0.7
- Update to Go 1.20: https://github.com/komuw/kama/pull/34
- Better formatting for zero-length slices and maps: https://github.com/komuw/kama/pull/36
- Zero value pointers should not panic: https://github.com/komuw/kama/pull/37
- Add tests for embedded structs with tags: https://github.com/komuw/kama/pull/41

# v0.0.6
- Update dependencies: https://github.com/komuw/kama/pull/32
- Update to Go 1.19:   https://github.com/komuw/kama/pull/33

## v0.0.5
- Check if terminal supports color before printing stack traces: https://github.com/komuw/kama/pull/31

## v0.0.4
- Add ability to print stack traces: https://github.com/komuw/kama/pull/29
  The stack traces are colorized with different colors for your code, third-party libs & the Go stdlib/runtime.
  Also code snippet for the most recent call is shown.
- Stop compacting data structures: https://github.com/komuw/kama/pull/28
- Add errcheck to CI: https://github.com/komuw/kama/pull/23

## v0.0.3
- Update CI: https://github.com/komuw/kama/pull/17   
- Dump more information about variables/types: https://github.com/komuw/kama/pull/18      
                                             : https://github.com/komuw/kama/pull/21       
- Implement own `dump` functionality: https://github.com/komuw/kama/pull/22     
  We used to use `sanity-io/litter` to do dumping.      
  We however, decided to implement our own dump functionality.       
  The main reason precipitating we are doing this is because sanity-io/litter has no way to compact       
  arrays/slices/maps that are inside structs.        

## v0.0.2
- Add test example: https://github.com/komuw/kama/pull/13
- Add types to the fields of a struct: https://github.com/komuw/kama/pull/16

## v0.0.1
- Pretty print variables and packages: https://github.com/komuw/kama/pull/10
- Add cli: https://github.com/komuw/kama/pull/11
- Add pretty printing for data structures: https://github.com/komuw/kama/pull/12
