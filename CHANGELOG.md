# Release Notes

Most recent version is listed first.  


# v0.0.6
- update dependencies: https://github.com/komuw/kama/pull/32
- update to Go 1.19:   https://github.com/komuw/kama/pull/33

## v0.0.5
- check if terminal supports color before printing stack traces: https://github.com/komuw/kama/pull/31

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
- add test example: https://github.com/komuw/kama/pull/13
- add types to the fields of a struct: https://github.com/komuw/kama/pull/16

## v0.0.1
- pretty print variables and packages: https://github.com/komuw/kama/pull/10
- add cli: https://github.com/komuw/kama/pull/11
- add pretty printing for data structures: https://github.com/komuw/kama/pull/12
