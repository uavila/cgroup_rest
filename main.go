/*
   FastCGI program that provides a restful API to manage the cgroups of a server. It should support the following:
   list available cgroups
   list the tasks (PIDs) for a given cgroup
   place a process into a cgroup
   */

package main

import (
    "net"
    "net/http"
    "net/http/fcgi"
    "fmt"
    "os"
    "encoding/csv"
    "encoding/json"
    "io/ioutil"
    "strings"
)

type FastCGIServer struct{}
type Subsys struct{
    Subsys_name string `json:"subsys_name,omitempty"`
    Hierarchy string `json:"hierarchy,omitempty"`
    Num_cgroups string `json:"num_cgroups,omitempty"`
    Enabled string `json:"enabled,omitempty"`
}
type Cgroup struct{
    Subsys_name string
    Cgroup_name string
}

func list_subsys() ([]Subsys, error){
        file, err := os.Open("/proc/cgroups")
        if err != nil {
            fmt.Println("Failed to open %s", err)
            return nil, err
        }
        reader := csv.NewReader(file)
        reader.Comma = '\t'
        lines, err := reader.ReadAll()
        if err != nil{
            fmt.Println("Error reading all lines: %v", err)
        }
        fmt.Printf("%v",lines)
        subsys := make([]Subsys, len(lines))

        for i, line := range lines {
            if line[0] !="#subsys_name" && line[0] !="" { // Make sure we're not hitting the header row
              t_subsys := Subsys{
                    Subsys_name: line[0],
                    Hierarchy: line[1],
                    Num_cgroups: line[2],
                    Enabled: line[3],
               }
               subsys[i-1] = t_subsys
            }else{
            }
        }
        return subsys, err 
}

func list_cgroups () []Cgroup{
       files, _ := ioutil.ReadDir("/sys/fs/cgroup")
       cgroups := make([]Cgroup, len(files))
       fmt.Println(len(files))
           for i, f := range files {
               if !f.IsDir() && f.Name() != "notify_on_release" && f.Name() != "release_agent"  && f.Name() != "tasks"{
                   filename := strings.SplitN(f.Name(),".",2)
                   //fmt.Printf("%i subsys %s - cgroup %s \n", i, filename[0], filename [1])
                   cgroups[i].Subsys_name = filename[0] 
                   cgroups[i].Cgroup_name = filename[1] 
                }
          }
       return cgroups
}
func h_list_subsys(resp http.ResponseWriter, req *http.Request){
        //http://localhost/zalora/subsys
        resp.Header().Set("Content-Type","application/json; charset=utf-8")

        subsys, err := list_subsys()
        if err != nil{
            fmt.Println("Error getting subsystems: ", err)
        }
        b, err := json.Marshal(subsys)
        resp.Write(b)
}

func h_list_cgroups(resp http.ResponseWriter, req *http.Request){
        //http://localhost/zalora/cgroups
        resp.Header().Set("Content-Type","application/json; charset=utf-8")
        cgroups := list_cgroups()
        b,err  := json.Marshal(cgroups)
        if err != nil{
            fmt.Println("Error getting cgroups: ", err)
        }
        resp.Write(b)
}

func h_list_cg_tasks(resp http.ResponseWriter, req *http.Request){
    //PUT /gists/:id/star
    //cat /sys/fs/cgroup/cpuset/Charlie/tasks
}

func h_put_task_cg(resp http.ResponseWriter, req *http.Request){
    //PUT /cgroups/:id/tasks/:tid
    //echo PID > /sys/fs/cgroup/cpuset/Charlie/tasks
}

func (s FastCGIServer) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
    baseurl := "/zalora"
    fmt.Println(req.URL.Path)
    params := req.URL.Query()
    fmt.Printf("Params: %#v\n",params)
    group := params.Get("cgroup")
    fmt.Println("Group: ", group)
//    id, err := strconv.Atoi(params.Get(":id"))
    if group != "" {
            h_list_cg_tasks(resp, req)
    } else {
        if req.URL.Path == baseurl + "/subsys" {
            h_list_subsys(resp, req)
    }  else {
        if req.URL.Path == baseurl + "/cgroups" {
            h_list_cgroups(resp, req)
        }
    }
    }
}


func main() {
    listener, _ := net.Listen("tcp", "127.0.0.1:8000")
    srv := new(FastCGIServer)
    fcgi.Serve(listener, srv)
}
