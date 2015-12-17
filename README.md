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
    - An array of requests to be executed
    
- Define you first array of requests:
    - URL
    - Method
    - Payload
    - Assertion of the response containing status code, timeout.
    - Extraction

## What is a stage ?
A stage in a thread group defines at which point this thread group will be executed.

- If 3 thread groups are defined with stage 1, they will all be executed in parallel.
- If a fourth thread groups is defined with stage 2, once the three first groups are finished, it will be executed

## What is an extraction ?
An extraction is a way to dynamically get results from the response and inject it in the further requests.

It needs the following fields:

- **var**: The variable name to store the result into
- **field**: Either "headers" or "body". Field to search into for a value
- **regex**: Any regular expression to find a dynamic value

Once the value is extracted, it's possible to use it in the further requests. Example below:

- *var*: "OAUTH_TOKEN"
- *field*: "headers"
- *regex*: "bearer (.*)?"

The variable $(OAUTH_TOKEN) can be injected in further requests afterwards
