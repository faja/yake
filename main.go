package main

import (
  "bytes"
  "flag"
  "fmt"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "os"
  "os/exec"
  "strings"
)

func main() {

  // argument parsing
  flagTask := flag.String("task", "default", "task to execute")
  flagFile := flag.String("file", "yakefile.yml", "yake file")
  flagKeepgoing := flag.Bool("keepgoing", false, "execute remaining steps even one of them fails")
  flagStdout := flag.Bool("stdout", false, "prints stdout")
  flagStderr := flag.Bool("stderr", false, "prints stderr")
  flag.Parse()

  // variables parsing
  variables := make(map[string]string)
  for _, v := range flag.Args() {
    vSplited := strings.Split(v, "=")
    if len(vSplited) < 2 {
      fmt.Println("[ERROR]",v,"variable syntax incorrect, try NAME=VALUE")
    }
    variables[vSplited[0]] = vSplited[1]
  }

  // YAML file structs
  type Task struct {
    Steps []string
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
  if _, ok := tasks[*flagTask]; ok != true {
    fmt.Println("Couldn't find task", *flagTask, "in the yakefile")
    os.Exit(1)
  }
  // execute steps
  for _,command := range tasks[*flagTask].Steps {
    for k,v := range variables {
      command = strings.Replace(command,"$"+k,v,-1)
    }
    taskSplitted := strings.Split(command, " ")
    fmt.Println(">>>", command)
    cmd := exec.Command(taskSplitted[0], taskSplitted[1:]...)

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
      os.Exit(1)
    }
  }
}
