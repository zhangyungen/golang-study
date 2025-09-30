// templates/model.tpl
package model

import (
    "time"
)

{{$length := len .Tables}}
{{range $i, $table := .Tables}}
type {{TableMapper $table.Name}} struct {
    {{range $table.ColumnsSeq}}{{$col := $table.GetColumn .}}`xorm:"'{{$col.Name}}' {{$col.SQLType}} {{if $col.IsPrimaryKey}}pk{{end}} {{if $col.IsAutoIncrement}}autoincr{{end}} {{if $col.Default}}default {{$col.Default}}{{end}} {{if not $col.Nullable}}notnull{{end}} {{if $col.Unique}}unique{{end}} {{if $col.IsCreated}}created{{end}} {{if $col.IsUpdated}}updated{{end}}" json:"{{$col.LowerName}}"`
    {{$col.Name}} {{Type $col}} `xorm:"'{{$col.Name}}' {{$col.SQLType}} {{if $col.IsPrimaryKey}}pk{{end}} {{if $col.IsAutoIncrement}}autoincr{{end}} {{if $col.Default}}default {{$col.Default}}{{end}} {{if not $col.Nullable}}notnull{{end}} {{if $col.Unique}}unique{{end}} {{if $col.IsCreated}}created{{end}} {{if $col.IsUpdated}}updated{{end}}" json:"{{$col.LowerName}}"`
    {{end}}
}

func ({{TableMapper $table.Name}}) TableName() string {
    return "{{$table.Name}}"
}
{{end}}