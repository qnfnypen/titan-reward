func (m *default{{.upperStartCamelObject}}Model) Insert(ctx context.Context, session sqlx.Session, datas ...*{{.upperStartCamelObject}}) (sql.Result,error) {
	sq := squirrel.Insert(m.table).Columns({{.lowerStartCamelObject}}RowsExpectAutoSet)
	{{if .withCache}}
	var keys []string
	for _, data := range datas {
		sq = sq.Values({{.expressionValues}})
		{{.keys}}
		keys = append(keys, {{.keyValues}})
	}
	query, args, err := sq.ToSql()
	if err != nil {
		return nil, err
	}
    ret, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		if session != nil {
			return session.ExecCtx(ctx, query, args...)
		}
		return conn.ExecCtx(ctx, query, args...)
	}, keys...){{else}}
	for _, data := range datas {
		sq = sq.Values({{.expressionValues}})
	}
	query, args, err := sq.ToSql()
	if err != nil {
		return nil, err
	}
	if session != nil {
		ret, err := session.ExecCtx(ctx, query, args...)
	}else{
		ret, err := m.conn.ExecCtx(ctx, query, args...)
	}{{end}}
    
	return ret, err
}