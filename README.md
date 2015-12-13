[![Build Status](https://travis-ci.org/dudang/golt.svg?branch=master)](https://travis-ci.org/dudang/golt)
[![Coverage Status](https://coveralls.io/repos/dudang/golt/badge.svg?branch=master&service=github)](https://coveralls.io/github/dudang/golt?branch=master)
[![Stories in Ready](https://badge.waffle.io/dudang/golt.png?label=ready&title=Ready)](http://waffle.io/dudang/golt)

# Golt
Golt is a load testing tool written in Go currently supporting only HTTP requests.

The goal is to support complex flow of requests replicating a parallel or sequential traffic.

```
$ ./golt -f <path-to-test-file> -l <path-to-log-file>
```

# How it works
Define your test plan with a JSON or YAML file and run golt!
- Example of the syntax can be found in the test folder

- Define your first thread group containing:
    - Amount of threads
    - Amount of repetitions of requests
    - Stage of the group
    - Type of group
    - An array of requests to be executed
    
- Define you first array of requests:
    - URL
    - Method
    - Payload
    - Assertion of the response containing status code, body, timeout.

## What is a stage ?
A stage in a thread group defines at which point this thread group will be executed.

- If 3 thread groups are defined with stage 1, they will all be executed in parallel.
- If a fourth thread groups is defined with stage 2, once the three first groups are finished, it will be executed

## What are type of groups ?
### "sequential"
In a "sequential" thread group, the array of requests in this group will be executed sequentially from top to bottom
### "parallel"
In a "parallel" thread group, the array of requests in this group will be executed in parallel.
