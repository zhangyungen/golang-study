// templates/dao.tpl
package dao

import (
    "errors"
    "xorm.io/xorm"
    "your-project-path/model"
)

{{range $table := .Tables}}
// {{TableMapper $table.Name}}DAO {{$table.Comment}}
type {{TableMapper $table.Name}}DAO struct {
    *BaseDAO[model.{{TableMapper $table.Name}}, {{GetPrimaryKeyType $table}}]
}

// 全局{{TableMapper $table.Name}}DAO实例
var {{TableMapper $table.Name}}DaoInstance = &{{TableMapper $table.Name}}DAO{&BaseDAO[model.{{TableMapper $table.Name}}, {{GetPrimaryKeyType $table}}]{}}

{{$hasEmail := false}}
{{range $col := $table.Columns}}
{{if eq $col.LowerName "email"}}{{$hasEmail = true}}{{end}}
{{end}}

{{if $hasEmail}}
// Create{{TableMapper $table.Name}} 创建{{$table.Comment}}
func (d *{{TableMapper $table.Name}}DAO) Create{{TableMapper $table.Name}}(session *xorm.Session, {{LowerFirst (TableMapper $table.Name)}} *model.{{TableMapper $table.Name}}) error {
    // 检查邮箱是否已存在
    exist, err := d.ExistByEmail(session, {{LowerFirst (TableMapper $table.Name)}}.Email)
    if err != nil {
        return err
    }
    if exist {
        return errors.New("email already exists")
    }
    return d.BaseDAO.Insert(session, {{LowerFirst (TableMapper $table.Name)}})
}

// GetByEmail 根据邮箱获取{{$table.Comment}}
func (d *{{TableMapper $table.Name}}DAO) GetByEmail(session *xorm.Session, email string) (*model.{{TableMapper $table.Name}}, error) {
    var {{LowerFirst (TableMapper $table.Name)}} model.{{TableMapper $table.Name}}
    has, err := session.Where("email = ?", email).Get(&{{LowerFirst (TableMapper $table.Name)}})
    if err != nil {
        return nil, err
    }
    if !has {
        return nil, errors.New("{{LowerFirst (TableMapper $table.Name)}} not found")
    }
    return &{{LowerFirst (TableMapper $table.Name)}}, nil
}

// Update 更新{{$table.Comment}}
func (d *{{TableMapper $table.Name}}DAO) Update(session *xorm.Session, {{LowerFirst (TableMapper $table.Name)}} *model.{{TableMapper $table.Name}}) error {
    // 检查{{$table.Comment}}是否存在
    existing, err := d.GetByID(session, {{LowerFirst (TableMapper $table.Name)}}.ID)
    if err != nil {
        return err
    }
    // 如果邮箱有变更，检查新邮箱是否被其他{{$table.Comment}}使用
    if {{LowerFirst (TableMapper $table.Name)}}.Email != existing.Email {
        exist, err := d.ExistByEmail(session, {{LowerFirst (TableMapper $table.Name)}}.Email)
        if err != nil {
            return err
        }
        if exist {
            return errors.New("email already used by another {{LowerFirst (TableMapper $table.Name)}}")
        }
    }
    return d.BaseDAO.UpdateById(session, {{LowerFirst (TableMapper $table.Name)}}.ID, {{LowerFirst (TableMapper $table.Name)}})
}

// ExistByEmail 检查邮箱是否存在
func (d *{{TableMapper $table.Name}}DAO) ExistByEmail(session *xorm.Session, email string) (bool, error) {
    return session.Where("email = ?", email).Exist(&model.{{TableMapper $table.Name}}{})
}
{{else}}
// Create{{TableMapper $table.Name}} 创建{{$table.Comment}}
func (d *{{TableMapper $table.Name}}DAO) Create{{TableMapper $table.Name}}(session *xorm.Session, {{LowerFirst (TableMapper $table.Name)}} *model.{{TableMapper $table.Name}}) error {
    return d.BaseDAO.Insert(session, {{LowerFirst (TableMapper $table.Name)}})
}

// Update 更新{{$table.Comment}}
func (d *{{TableMapper $table.Name}}DAO) Update(session *xorm.Session, {{LowerFirst (TableMapper $table.Name)}} *model.{{TableMapper $table.Name}}) error {
    return d.BaseDAO.UpdateById(session, {{LowerFirst (TableMapper $table.Name)}}.ID, {{LowerFirst (TableMapper $table.Name)}})
}
{{end}}
{{end}}