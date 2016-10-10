yake = **Y**et**A**notherma**K**eclon**E**


example usage:
```
% ./yake -file examples/yakefile.yml -task task1 NAME=foo
>>> cp /tmp/foo.yml.tar .
>>> tar -xvf foo.yml.tar
>>> ls foo.yml
```
flags:

* `-task TASKNAME` - specify task name, default: `default`
* `-file FILENAME` - specify yake file name, default: `yakefile.yml`
* `-keepgoing` - yake continues to execute all remaining steps even if one of them fails, default: `false`
* `-stdout` - prints STDOUT of executing steps, default: `false`
* `-stderr` - print STDERR  of executing steps, default: `false`
