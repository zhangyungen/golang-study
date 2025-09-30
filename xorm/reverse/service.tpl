// templates/service.tpl
package service

import (
    "errors"
    "your-project-path/dao"
    "your-project-path/model"
    "your-project-path/param"
)

{{range $table := .Tables}}
// {{TableMapper $table.Name}}Service {{$table.Comment}}Service
type {{TableMapper $table.Name}}Service struct {
    *BaseService[model.{{TableMapper $table.Name}}, {{GetPrimaryKeyType $table}}]
    {{LowerFirst (TableMapper $table.Name)}}DAO *dao.{{TableMapper $table.Name}}DAO
}

// 全局{{TableMapper $table.Name}}Service实例
var {{TableMapper $table.Name}}ServiceIns = &{{TableMapper $table.Name}}Service{
    &BaseService[model.{{TableMapper $table.Name}}, {{GetPrimaryKeyType $table}}]{},
    dao.{{TableMapper $table.Name}}DaoInstance,
}

{{$hasEmail := false}}
{{range $col := $table.Columns}}
{{if eq $col.LowerName "email"}}{{$hasEmail = true}}{{end}}
{{end}}

// Create{{TableMapper $table.Name}} 创建{{$table.Comment}}
func (s *{{TableMapper $table.Name}}Service) Create{{TableMapper $table.Name}}({{LowerFirst (TableMapper $table.Name)}} *model.{{TableMapper $table.Name}}) error {
    // 数据验证
    if {{LowerFirst (TableMapper $table.Name)}}.Name == "" {
        return errors.New("name is required")
    }
    {{if $hasEmail}}
    if {{LowerFirst (TableMapper $table.Name)}}.Email == "" {
        return errors.New("email is required")
    }
    {{end}}
    session := s.getDBSession()
    defer s.closeDBSession(session)
    return s.{{LowerFirst (TableMapper $table.Name)}}DAO.Create{{TableMapper $table.Name}}(session, {{LowerFirst (TableMapper $table.Name)}})
}

// Get{{TableMapper $table.Name}} 获取{{$table.Comment}}
func (s *{{TableMapper $table.Name}}Service) Get{{TableMapper $table.Name}}(id {{GetPrimaryKeyType $table}}) (*model.{{TableMapper $table.Name}}, error) {
    if id <= 0 {
        return nil, errors.New("invalid {{LowerFirst (TableMapper $table.Name)}} id")
    }
    session := s.getDBSession()
    defer s.closeDBSession(session)
    return s.{{LowerFirst (TableMapper $table.Name)}}DAO.GetByID(session, id)
}


// Update{{TableMapper $table.Name}} 更新{{$table.Comment}}
func (s *{{TableMapper $table.Name}}Service) Update{{TableMapper $table.Name}}({{LowerFirst (TableMapper $table.Name)}} *model.{{TableMapper $table.Name}}) error {
    if {{LowerFirst (TableMapper $table.Name)}}.ID <= 0 {
        return errors.New("invalid {{LowerFirst (TableMapper $table.Name)}} id")
    }
    session := s.getDBSession()
    defer s.closeDBSession(session)
    return s.{{LowerFirst (TableMapper $table.Name)}}DAO.Update(session, {{LowerFirst (TableMapper $table.Name)}})
}

// Delete{{TableMapper $table.Name}} 删除{{$table.Comment}}
func (s *{{TableMapper $table.Name}}Service) Delete{{TableMapper $table.Name}}ById(id {{GetPrimaryKeyType $table}}) error {
    if id <= 0 {
        return errors.New("invalid {{LowerFirst (TableMapper $table.Name)}} id")
    }
    session := s.getDBSession()
    defer s.closeDBSession(session)
    return s.{{LowerFirst (TableMapper $table.Name)}}DAO.DeleteById(session, id, &model.{{TableMapper $table.Name}}{})
}

// List{{TableMapper $table.Name}}s {{$table.Comment}}列表
func (s *{{TableMapper $table.Name}}Service) PageList(param *param.PageParam) ([]*model.{{TableMapper $table.Name}}, error) {
    if param.Page <= 0 {
        param.Page = 1
    }
    if param.PageSize <= 0 || param.PageSize > 100 {
        param.PageSize = 20
    }
    session := s.getDBSession()
    defer s.closeDBSession(session)
    return s.{{LowerFirst (TableMapper $table.Name)}}DAO.PageList(session, param)
}


// Validate{{TableMapper $table.Name}} 验证{{$table.Comment}}数据
func (s *{{TableMapper $table.Name}}Service) Validate{{TableMapper $table.Name}}(name string) error {
    return nil
}

{{end}}