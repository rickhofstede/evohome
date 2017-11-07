## Library for connecting to Honeywell Evohome using Go

Copyright (c) 2017, Rick Hofstede

All rights reserved. This software is distributed under a BSD-style
license. For more information, see LICENSE.

### 1. Introduction

This library can be used for building applications that utilize
Honeywell's public APIs to contol Evohome systems.

### 2. How to use

```golang
client = evohome.NewEvohome(username, password)
if client == nil || !client.Initialized() {
    fmt.Println("\nConnection/authentication error")
    os.Exit(0)
    return
}

t := client.TemperatureControlSystem()
zones := t.ZoneNames()
```

### 3. Support

Please request support by creating an 'issue' [here](https://github.com/rickhofstede/evohome/issues).
