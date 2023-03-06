# Coverage Table

## Overview

CLI tool to aggregate coverage profiles.

## Input

use coverprofile data which is exported by go test, like the following commands

```
go test -coverprofile=coverage.txt ./...
```

example
```
mode: atomic
github.com/sample/a/b/b-file.go:5.47,10.2 4 1
github.com/sample/a/b/b-file.go:11.47,16.2 4 1
github.com/sample/a/b/b-file.go:18.47,22.2 2 0
github.com/sample/a/a-file1.go:5.47,22.2 6 0
```

## Output 

* output the following table in csv
    * aggregate the number of lines as well as coverage
    * aggregate by directory or file unit

example
| filepath                        | coverage | total statements | tested | not tested | 
| :-----------------------------: | :------: | :--------------: | :----: | :--------: | 
| github.com/sample               | 50.00    | 16               | 8      | 8          | 
| github.com/sample/a             | 50.00    | 16               | 8      | 8          | 
| github.com/sample/a/a-file1.go  | 0.00     | 6                | 0      | 6          | 
| github.com/sample/a/b           | 80.00    | 10               | 8      | 2          | 
| github.com/sample/a/b/b-file.go | 80.00    | 10               | 8      | 2          | 

## usage
```
go run ./*.go <rootPath> <InputFilePath> <OutputPath> 

example:
go run ./*.go github.com/sample ./sample_cover.txt sample_result.csv
```


