// templates/biz.tpl
package biz

import (
    "your-project-path/param"
    "your-project-path/service"
)

{{range $table := .Tables}}
// {{TableMapper $table.Name}}CmdBiz {{$table.Comment}}CmdBiz
type {{TableMapper $table.Name}}CmdBiz struct {
    *BaseCmdBiz
    *service.{{TableMapper $table.Name}}Service
}

// 全局{{TableMapper $table.Name}}CmdBizIns实例
var {{TableMapper $table.Name}}CmdBizIns = &{{TableMapper $table.Name}}CmdBiz{BaseCmdBizIns, service.{{TableMapper $table.Name}}ServiceIns}



{{end}}