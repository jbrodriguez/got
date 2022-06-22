# got - go time tracker

## Introduction

Simple, lightweight time tracking application built in [Go](https://go.dev/).

_Inspired by [Ultimate Time Tracker](https://github.com/larose/utt)_

## How To Use

```bash
# start the day
got on

# add some work
got add "project1: work on the project"

got add "project2: some other work"

# time for lunch
# break is used to indicate any non-work period
got add "break: lunch"

# back from luch
got add "catchup: emails"

# there's always meetings right ?
got add "meeting: some very productive meeting"

# get a report up to now
$ got report -w

------------------- Jun 06 to Jun 12 (week 23) -------------------------------

             Jun 06  Jun 07  Jun 08  Jun 09  Jun 10  Jun 11  Jun 12     Total

    project   00:00   00:00   00:00   03:30   05:21   00:00   00:00     08:51
       main   00:00   00:06   00:00   01:47   02:16   00:00   00:00     04:09
    catchup   00:00   00:13   00:00   01:12   00:55   00:00   00:00     02:21
    meeting   00:00   00:00   00:00   01:28   00:13   00:00   00:00     01:41
         dm   00:00   00:00   00:00   00:51   00:37   00:00   00:00     01:29
   incident   00:00   00:00   00:00   00:11   00:30   00:00   00:00     00:42
  pr-review   00:00   00:00   00:00   00:18   00:00   00:00   00:00     00:18

    working   00:00   00:19   00:00   09:18   09:53   00:00   00:00     19:30
      break   00:00   00:00   00:00   01:22   00:49   00:00   00:00     02:12
              -----   -----   -----   -----   -----   -----   -----     -----
      total   00:00   00:19   00:00   10:41   10:42   00:00   00:00     21:42

# how to end the day ?
# just do nothing
# the last `got add` you made reflects the end of the day
# the next `got on` will start a new day
```

## Supported flag behavior

```bash
# today
got report
got report -t

# any day
got report 2022-06-09
got report -s 2022-06-09

# current week
got report -w

# any week
got report -w 2022-06-09

# current month
got report -m

# any month
got report -m 2022-05-01

# calendar month (all weeks which include the current month)
got report -c

# any calendar month
got report -c 2022-05-01
```

### Acknowledgements

- [kong](github.com/alecthomas/kong)
- [tabwriter](github.com/Ladicle/tabwriter)
- [aurora](github.com/logrusorgru/aurora/v3)

#
