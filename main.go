package main

import (
  "bytes"
  "flag"
  "fmt"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "os"
  "os/exec"
  "regexp"
  "strings"
  "syscall"
)

func main() {

  // argument parsing
  flagFile := flag.String("file", "Yakefile", "yake file")
  flagKeepgoing := flag.Bool("keepgoing", false, "execute remaining steps even one of them fails (default false)")
  flagStdout := flag.Bool("stdout", false, "prints stdout (default false)")
  flagStderr := flag.Bool("stderr", false, "prints stderr (default false)")
  flag.Parse()

  flagsSet := make(map[string]bool)
  for _,v := range os.Args[1:flag.NFlag()+1] {
    flagsSet[strings.Split(v, "=")[0]] = true
  }

  // arguments parsing
  // any argument containing = character it's a variable
  // first argument is a task name
  // all other arguments create CMD variable, which can be used in `steps`
  var task string
  var defaultCmd string
  variables := make(map[string]string)

  for _, v := range flag.Args() {
    vSplited := strings.Split(v, "=")
    if len(vSplited) < 2 || string(vSplited[0][0]) == "-" {
      if len(task) == 0 {
        task = v
      } else {
        if len(defaultCmd) > 0 {
          defaultCmd += " "
        }
        defaultCmd += v
      }
    } else {
      variables[vSplited[0]] = strings.Join(vSplited[1:], "=")
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
  }

  // read the yakefile
  yamlfile, err := ioutil.ReadFile(*flagFile)
  if err != nil {
    fmt.Println(err.Error())
    os.Exit(1)
  }

  // parse yaml file
  var tasks map[string]Task
  if err := yaml.Unmarshal(yamlfile, &tasks); err != nil {
    fmt.Println(err.Error())
    os.Exit(1)
  }

  // does task exist?
  if _, ok := tasks[task]; ok != true {
    fmt.Println("Couldn't find task", task, "in the yakefile")
    os.Exit(1)
  }

  // use default variable value from task definition if not specified in command line
  for k,v := range tasks[task].Vars {
    if _, ok := variables[k]; ok != true {
      variables[k] = v
    }
  }

  // use default variable from _config if not specified anywhere else
  for k,v := range tasks["_config"].Vars {
    if _, ok := variables[k]; ok != true {
      variables[k] = v
    }
  }

  // _config parsing
  if tasks["_config"].Keepgoing && ! flagsSet["-keepgoing"] {
    *flagKeepgoing = tasks["_config"].Keepgoing
  }
  if tasks["_config"].Stdout && ! flagsSet["-stdout"] {
    *flagStdout = tasks["_config"].Stdout
  }
  if tasks["_config"].Stderr && ! flagsSet["-stderr"] {
    *flagStderr = tasks["_config"].Stderr
  }

  // variable regexp
  r := regexp.MustCompile("\\$[a-zA-Z0-9-_]+")

  // execute steps
  for _,command := range tasks[task].Steps {

    // variables could contain other variables
    for r.MatchString(command) {
      m := r.FindString(command)
      command = strings.Replace(command,m,variables[strings.TrimPrefix(m,"$")],-1)
    }

    // CMD variable
    command = strings.Replace(command,"$CMD",defaultCmd,-1)
    fmt.Println(">>>", command)
    cmd := exec.Command("sh", "-c", command)

    // output buffers
    cmdStdout := &bytes.Buffer{}
    cmdStderr := &bytes.Buffer{}
    cmd.Stdout = cmdStdout
    cmd.Stderr = cmdStderr

    err := cmd.Run()

    // print stdout
    if *flagStdout {
      if len(cmdStdout.Bytes()) > 0 {
        fmt.Printf("%s\n", cmdStdout.Bytes())
      }
    }

    // print stderr
    if *flagStderr {
      if len(cmdStderr.Bytes()) > 0 {
        os.Stderr.WriteString(fmt.Sprintf("%s\n", cmdStderr.Bytes()))
      }
    }

    // keepgoing?
    if err != nil && ! *flagKeepgoing {
      os.Stderr.WriteString(fmt.Sprintf("%s\n", err.Error()))
      os.Exit(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus())
    }
  }
}
