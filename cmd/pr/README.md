# pr

Fetch URL of created GitHub pull-requests that at since a week before.

## Usage

option requirements:

```
-r string
    full name of repository (e.g, hatajoe/gh)
```

```
% ./pr -r hatajoe/test
2018-05-30 04:15:35 +0000 UTC hatajoe [#1 pull-request title #1](https://github.com/hatajoe/test/pull/1)
2018-05-30 01:06:12 +0000 UTC hatajoe [#2 pull-request title #2](https://github.com/hatajoe/test/pull/2)
2018-05-30 00:58:13 +0000 UTC hatajoe [#3 pull-request title #3](https://github.com/hatajoe/test/pull/3)
```
