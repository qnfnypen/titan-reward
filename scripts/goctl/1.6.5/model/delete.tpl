func (m *default{{.upperStartCamelObject}}Model) Delete(ctx context.Context, session sqlx.Session, {{.lowerStartCamelPrimaryKey}}s ...{{.dataType}}) error {
	if len({{.lowerStartCamelPrimaryKey}}s) == 0 {
		return nil
	}
	{{if .withCache}}var keys []string{{if .containsIndexCache}}
	for _,{{.lowerStartCamelPrimaryKey}} := range {{.lowerStartCamelPrimaryKey}}s {
		data, _:=m.FindOne(ctx,session,{{.lowerStartCamelPrimaryKey}})
		if data==nil{
			continue
		}
		{{.keys}}
		keys = append(keys, {{.keyValues}})
	}{{else}}
	for _,{{.lowerStartCamelPrimaryKey}} := range {{.lowerStartCamelPrimaryKey}}s {
		{{.keys}}
		keys = append(keys, {{.keyValues}})
	}
{{end}}	
	vb, _ := json.Marshal({{.lowerStartCamelPrimaryKey}}s)
	values := strings.TrimSuffix(strings.TrimPrefix(string(vb), "["), "]")
	query := fmt.Sprintf("delete from %s where {{.originalPrimaryKey}} IN ({{if .postgreSql}}$1{{else}}%s{{end}})", m.table,values)
    _, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		if session != nil {
			return session.ExecCtx(ctx, query)
		}
		return conn.ExecCtx(ctx, query)
	}, keys...){{else}}
		if session != nil {
			_,err := session.ExecCtx(ctx, query)
		}else{
			_,err := m.conn.ExecCtx(ctx, query)
		}{{end}}
	return err
}
