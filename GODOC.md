# RJSON Go Library Documentation

## Overview

This library allows you to serialize (marshal) most Go datatypes, arrays, and maps into a UTF-8 JSON string that 
is formatted/organized per the RJSON specification.

Because the RJSON specification strictly adheres to the JSON spec, there is no need for a corresponding
de-serialization (unmarshal) function. Simply use the 

## Decimal Support

Sadly, the Go language does not support floating point decimal numbers in its standard library, only binary floating
point numbers. Technically, "numbers" in the JSON spec are decimal not binary. 

As a work-around, the library also supports the 128bit Decimal portion of BSON "go.mongodb.org/mongo-driver/v2/bson" 
3rd-party mongodb library. Only the BSON specification decimal part is used. You do NOT need to use the mongodb
driver in any way.

## Usage

Simply import the library and call `Marshal` with a variable. The Marshal function returns with a ([]byte, error) tuple.

For example:

```go
package main

import (
    "fmt"
    "go.mongodb.org/mongo-driver/v2/bson"
    "github.com/JohnAD/RJSON"
)

func main() {
    doc := map[string]interface{}{
        "d_obj_name": map[string]interface{}{
            "111": 44.3,
            "zip": bson.ParseDecimal128("12.00"),
            "222": false,
        },
        "a_str_name": "foo",
    }
    result, _ := RJSON.Marshal(string(doc))
    fmt.PrintLn(result)
}
```

Results in the following:

```json
{
  "a_str_name": "foo",
  "d_obj_name": {
    "111": 44.3,
    "222": false,
    "zip": 12.00
  }
}
```

Please note that object fields (strings) are sorted and do not honor the original documents order. See the RJSON spec
for details.
