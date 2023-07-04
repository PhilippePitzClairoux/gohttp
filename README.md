# Go Http Server
This library tries to make golang http servers
"clean" and easily readable. There's still alot of
things to implement but there is currently the bare minimum.

Please refer to dummy.go and Example.go in order to learn
how this works.

## TODO LIST
- Give access to request headers when calling endpoint method
- If the output of a method is a struct, parse it to JSON (same thing for body of a HTTP call)
- Handle various header related stuff (mostly for authentication purposes)
- Optimize code (use more pointers instead of copying most data)
- Reduce cognitive complexity of functions
- document exported methods