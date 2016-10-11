yake = **Y**et**A**notherma**K**eclon**E**

## example usage:
`yake -file Yakefile.yml -stdout SOMEVAR=foo taskname CMD argumets here`

#### simple_task
```yaml
simple_task:
  steps:
    - echo hello yake
    - uptime
    - date
```
```
% ./yake -file examples/yakefile.yml -stdout -stderr simple_task
>>> echo hello yake
hello yake

>>> uptime
 11:14:02 up  1:12,  9 users,  load average: 0.37, 0.32, 0.35

 >>> date
 Tue 11 Oct 11:14:02 BST 2016

```
#### task_with_variables
```yaml
task_with_variables:
  steps:
    - cp /tmp/$NAME.yml.tar .
    - tar -xvf $NAME.yml.tar
    - ls $NAME.yml
  vars:
    NAME: foo
```
```
./yake -file examples/yakefile.yml task_with_variables
>>> cp /tmp/foo.yml.tar .
>>> tar -xvf foo.yml.tar
>>> ls foo.yml
```
```
./yake -file examples/yakefile.yml task_with_variables NAME=bar
>>> cp /tmp/bar.yml.tar .
>>> tar -xvf bar.yml.tar
>>> ls bar.yml
```

#### task_with_CMD
```yaml
task_with_CMD:
  steps:
    - echo $CMD
```
```
./yake -file examples/yakefile.yml task_with_CMD 42 42 42
>>> echo 42 42 42
```

## flags:

* `-file FILENAME` - specify yake file name, default: `yakefile.yml`
* `-keepgoing` - yake continues to execute all remaining steps even if one of them fails, default: `false`
* `-stdout` - prints STDOUT of executing steps, default: `false`
* `-stderr` - print STDERR  of executing steps, default: `false`
