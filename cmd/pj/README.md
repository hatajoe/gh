# pj

Fetch URL of GitHub issues or pull-requests that are referred to project card of `Done` column.

## Usage

option requirements:

```
-p string
    project name
-r string
    full name of repository (e.g, hatajoe/gh)
```

```
% ./pj -p test -r hatajoe/test
https://github.com/dev-cloverlab/leeap/pull/1 pull request title#1
https://github.com/dev-cloverlab/leeap/pull/2 pull request title#2
https://github.com/dev-cloverlab/leeap/pull/3 pull request title#3
https://github.com/dev-cloverlab/leeap/issues/4 issue title#4
```
