Go Low Level Logging
====================

This is a loging package that lets you log at 3 log levels (with four values of the log level, the fourth being "don't"

I conventionally have "always" log for relatively lower volume things, config values being changes, listen succeeding etc.  
"State" log for logging e.g. accepts or connects.  
"Network" log for logging the network traffic into my server and in response to the server request.  
I might add a "Debug" for doing things like tricky logic errors, but actually testing is better for that.  

It works like:

```
      ml.La('Go Routines for kafka pub', numKafkaHandlers) (logged except when level is "none")
      ml.Ls("Read timeout from", dstAddr)                  (logged when level is "always" or "state")
      ml.Ln('Got a msg:', msg)                             (logged when level is "network" or "always" or "state")
```

THe output is a rotating log provided by github.com/lestrrat-go/file-rotatelogs

Don't use 2fde954 or c6f4c4a or abb58ad. I broke
SetLogPath and made an incompatible APi change.

master is now good.  