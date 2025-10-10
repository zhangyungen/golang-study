package sql

import (
	"bytes"
	"log"
	"sync"
	"text/template"
	"zyj.com/golang-study/xorm/param"
)

var (
	userLogIn = `SELECT * FROM user_login_log
                     WHERE 1=1 and 
                         {{if  .Name  }}
                                login_ip = "{{.Name}}"
                          {{end}}`

	tmpl *template.Template
	once sync.Once
)

func GetUserLoginSql(create param.UserCreate) string {
	once.Do(initTemplate)
	var buf bytes.Buffer
	// 4. 执行模板渲染，将结果写入 buffer
	err := tmpl.Execute(&buf, create)
	if err != nil {
		log.Printf("GetUserLoginSql", err)
		return ""
	}
	// 5. 从 buffer 中获取字符串
	renderedString := buf.String()
	return renderedString
}
func initTemplate() {
	tmpl = template.Must(template.New("userLogin").Parse(userLogIn))
}
