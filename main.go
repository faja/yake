package main

import (
  "bytes"
  "flag"
  "fmt"
  "io/ioutil"
  "os"
  "os/exec"
  "regexp"
  "sort"
  "strings"
  "syscall"
  "github.com/smallfish/simpleyaml"
)

func main() {

  // config struct
  type config struct {
    file string
    flags map[string]bool
    vars map[string]string
    steps []string
  }

  // argument parsing
  flagFile := flag.String("file", "Yakefile", "yake file")
  flagKeepgoing := flag.Bool("keepgoing", false, "execute remaining steps even one of them fails (default false)")
  flagStdout := flag.Bool("stdout", false, "prints stdout (default false)")
  flagStderr := flag.Bool("stderr", false, "prints stderr (default false)")
  flagShowcmd := flag.Bool("showcmd", true, "prints executed command")
  flag.Parse()

  c := config{
    file: *flagFile,
    flags: make(map[string]bool),
    vars: make(map[string]string),
  }
  c.flags["keepgoing"] = *flagKeepgoing
  c.flags["stdout"] = *flagStdout
  c.flags["stderr"] = *flagStderr
  c.flags["showcmd"] = *flagShowcmd

  flagsSet := make(map[string]bool)
  for _,v := range os.Args[1:flag.NFlag()+1] {
    flagsSet[strings.Split(v, "=")[0]] = true
  }

  // arguments parsing
  // any argument containing = character it's a variable
  // first argument is a task name
  // all other arguments create CMD variable, which can be used in `steps`
  var task string

  for _, v := range flag.Args() {
    vSplited := strings.Split(v, "=")
    if len(vSplited) < 2 || string(vSplited[0][0]) == "-" {
      if len(task) == 0 {
        task = v
      } else {
        if len(c.vars["CMD"]) > 0 {
          c.vars["CMD"] += " "
        }
        c.vars["CMD"] += v
      }
    } else {
      c.vars[vSplited[0]] = strings.Join(vSplited[1:], "=")
    }
  }
  // task name defined?
  if len(task) == 0 {
    fmt.Println("Please specify task name")
    os.Exit(1)
  }

  // YAML file struct
  type Task struct {
    Steps []string
    Vars map[string]string
    Keepgoing bool
    Stdout bool
    Stderr bool
    Showcmd bool
  }

  // read the yakefile
  yamlfile, err := ioutil.ReadFile(c.file)
  if err != nil {
    fmt.Println(err.Error())
    os.Exit(1)
  }

  // parse yaml file
  yamldata, err := simpleyaml.NewYaml(yamlfile)
  if err != nil {
    fmt.Println(err.Error())
    os.Exit(1)
  }

  // build a yaml map
  yamlmap, err := yamldata.Map()
  if err != nil {
    fmt.Println(err.Error())
    os.Exit(1)
  }

  // does task exist?
  if _, ok := yamlmap[task]; !ok {
    var availableTasks []string
    for k,_ := range yamlmap {
      if k == "_config" {
        continue
      }
      switch k.(type) {
      case string:
        availableTasks = append(availableTasks,k.(string))
      default:
      }
    }
    sort.Strings(availableTasks)
    fmt.Println("Couldn't find task", task, "in the yakefile")
    fmt.Println("Available tasks:", strings.Join(availableTasks, ", "))
    os.Exit(1)
  }

  // parsin _config section
  if _, ok := yamlmap["_config"]; ok {

    // bools parsing
    bools := []string{"keepgoing", "stdout", "stderr", "showcmd"}
    for _,v := range bools {
      vv, err := yamldata.Get("_config").Get(v).Bool()
      if err == nil {
        if ! flagsSet[fmt.Sprintf("-%s",v)] {
          c.flags[v] = vv
        }
      }
    }

    // vars parsing
    vars, err := yamldata.Get("_config").Get("vars").Map()
    if err == nil {
      for k,v := range vars {
        switch k.(type) {
        case string:
          kstring := k.(string)
          // use default variable from _config if not specified on command line
          if _, ok := c.vars[kstring]; ok != true {
            c.vars[kstring] = v.(string)
          }
        default:
        }
      }
    }
  }

  // parsin task
  switch yamlmap[task].(type) {
  case string:
    step, _ := yamldata.Get(task).String()
    c.steps = append(c.steps, step)
  case []interface{}:
    steps, _ := yamldata.Get(task).Array()
    for _, step := range steps {
      switch step.(type) {
      case string:
        c.steps = append(c.steps, step.(string))
      }
    }
  default:
  }

  // variable regexp
  r := regexp.MustCompile("\\$[a-zA-Z0-9-_]+")

  // execute steps
  for _,command := range c.steps {

    // variables could contain other variables
    for r.MatchString(command) {
      m := r.FindString(command)
      command = strings.Replace(command,m,c.vars[strings.TrimPrefix(m,"$")],-1)
    }

    if c.flags["showcmd"] {
      fmt.Println(">>>", command)
    }
    cmd := exec.Command("sh", "-c", command)

    // output buffers
    cmdStdout := &bytes.Buffer{}
    cmdStderr := &bytes.Buffer{}
    cmd.Stdin = os.Stdin
    cmd.Stdout = cmdStdout
    cmd.Stderr = cmdStderr

    err := cmd.Run()

    // print stdout
    if c.flags["stdout"] {
      if len(cmdStdout.Bytes()) > 0 {
        fmt.Printf("%s\n", cmdStdout.Bytes())
      }
    }

    // print stderr
    if c.flags["stderr"] {
      if len(cmdStderr.Bytes()) > 0 {
        os.Stderr.WriteString(fmt.Sprintf("%s\n", cmdStderr.Bytes()))
      }
    }

    // keepgoing?
    if err != nil && ! c.flags["keepgoing"] {
      if c.flags["showcmd"] {
        os.Stderr.WriteString(fmt.Sprintf("%s\n", err.Error()))
      }
      os.Exit(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus())
    }
  }
}
