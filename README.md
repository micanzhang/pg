# pg [![Build Status](https://travis-ci.org/micanzhang/pg.svg?branch=master)](https://travis-ci.org/micanzhang/pg)
command line tools  for password generation.

## required 

1. go  (lastest version recommended)

## install 

```sh
$ go install github.com/micanzhang/pg
```

## basic usage

### generate password

```sh
$ pg #  pg gen
```

### list entry

```sh
$pg list
```
### new entry 

```sh
$pg new -d mysql -u root -p root -k passwordphrase
```

if **-p** is empty, pg will generate new one.

### update entry 

```sh
$pg update -d mysql -u root -k passwordphrase
```

### remove entry 

```sh
$pg new -d mysql -u root
```
### get  password 

```sh
$pg info -d mysql -u root -k passwordphrase
```
