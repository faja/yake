yake = **Y**et**A**notherma**K**eclon**E**

## usage:
`yake [flags] [variable_definition] taskname [CMD_arguments]`
eg:
* `yake taskename`
* `yake -file Yakefile.yml -stdout SOMEVAR=foo taskname foo bar`

#### simple_task
```yaml
simple_task:
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
  - cp /tmp/$NAME.yml.tar .
  - tar -xvf $NAME.yml.tar
  - ls $NAME.yml

_config:
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
./yake -file examples/yakefile.yml NAME=bar task_with_variables
>>> cp /tmp/bar.yml.tar .
>>> tar -xvf bar.yml.tar
>>> ls bar.yml
```

#### task_with_CMD
```yaml
task_with_CMD:
  - echo $CMD
```
```
./yake -file examples/yakefile.yml task_with_CMD 42 42 42
>>> echo 42 42 42
```

## flags:

* `-file FILENAME` - specify yake file name, default: `Yakefile`
* `-showcmd` - prints executed commands, default: `true`
* `-keepgoing` - yake continues to execute all remaining steps even if one of them fails, default: `false`
* `-stdout` - prints STDOUT of executing steps, default: `false`
* `-stderr` - prints STDERR  of executing steps, default: `false`
