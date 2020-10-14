# Log with Context and rpcid

The goal is to have very simple log function to call and help us analyze the output by link all the function call with rpcid accross system

## Basic use
This is the basic default use of the log out of the box
```
package main

import (
  "context"

  "github.com/mirzaakhena/log"
)

func main() {

  ctx := log.ContextWithRpcID(context.Background())

  log.Info(ctx, "hello")
  log.Warn(ctx, "world")
  log.Error(ctx, "my name is %s", "mirza")

}
```

The context can be from different previous service that call this service. This is sample code when using `https://github.com/gin-gonic/gin` framework

```
func TheController(c *gin.Context) {

  previousCtx := c.Request.Context()
  ctx := log.ContextWithRpcID(previousCtx)
  log.Info(ctx, "hello")

}
```

This is sample when using builtin http go
```
func TheController(w http.ResponseWriter, req *http.Request) {

  previousCtx := req.Context()
  ctx := log.ContextWithRpcID(previousCtx)
  log.Info(ctx, "hello")

}
```

This is the sample output format
```
{"func":"main.main:18","level":"info","msg":"hello","rpcid":"1iqnii541bXcCIXkbJ3OMvxrx6R","time":"1014 095518.829"}
{"func":"main.main:19","level":"warning","msg":"world","rpcid":"1iqnii541bXcCIXkbJ3OMvxrx6R","time":"1014 095518.829"}
{"func":"main.main:20","level":"error","msg":"my name is mirza","rpcid":"1iqnii541bXcCIXkbJ3OMvxrx6R","time":"1014 095518.829"}
```
From that output we can see the rpcid is same. We can use this rpcid information to do grep from console and by collect the same rpcid we can trace it easily

## Change the output format
Currently we have two format. Nested format and JSON format. Call this method before we call the first log. The default one is in JSON format
```
log.UseNestedFormat()
// or
log.UseJSONFormat()
```

## Replace The RPCID Generator
Currently rpcid is generate by `https://github.com/segmentio/ksuid`

You can replace the function generation by call this method at the first place before any log is called. In this example we replace it with uuid from  `https://github.com/satori/go.uuid`

```
log.SetRpcIDFunc(func() string {
  x, _ := uuid.NewV4()
  return x.String()
})
```

## Use the rotate file 
If we want to have log file we can enable it by call this before first log is called
```
// the current directory
path := "." 

// log filename
filename := "logfilename"

// after n max day. the file will deleted automatically
maxAgeOfLogInDays := 7 

log.UseRotateFile(path, filename, maxAgeOfLogInDays)
log.Info(context.Background(), "hello")

```
The output file will be like this
```
projectdir
  +-logs
  | +-logfilename.log.20201013
  | +-logfilename.log.20201014
  +-logfilename.log
  +-main.go
```





