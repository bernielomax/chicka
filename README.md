# chicka
A pluggable Go monitoring system




# Writing custom plugins

As long as your plugin outputs conforms to the following JSON standard your plugin will work:

```
{
    "description": "the total number of docker images must not exceed 5",  // describes the test
    "result":      false,                                                  // did the test pass?
    "data":        28,                                                     // the result data that the result logic was based on.
    "expect":      true                                                    // what result were you expecting?
}
```