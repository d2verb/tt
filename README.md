# tt
Convert time range to regexp. **still buggy**

## Installation
### go get
```
$ go get github.com/d2verb/tt/cmd/tt
```

## Usage
```
$ tt "2020-11-09 17:00:00" "2020-11-09 17:13:59"
2020-11-09 17:[0-1][0-9]:[0-5][0-9]
```