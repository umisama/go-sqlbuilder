package sqlbuilder

type LockClause interface {
	NoWait(noWait bool) LockClause
	serialize(bldr *builder)
}

type lockClauseImpl struct {
	strength string
	tables   []Table
	noWait   bool
}

func Lock(strength string, tables ...Table) LockClause {
	return &lockClauseImpl{
		strength: strength,
		tables:   tables,
	}
}

func (l *lockClauseImpl) NoWait(noWait bool) LockClause {
	l.noWait = noWait

	return l
}

func (l *lockClauseImpl) serialize(bldr *builder) {
	t := make([]serializable, len(l.tables))
	for i, v := range l.tables {
		t[i] = v
	}

	bldr.Append(l.strength)
	if len(l.tables) > 0 {
		bldr.Append(" OF ")
		bldr.AppendItems(t, ", ")
	}
	if l.noWait {
		bldr.Append(" NOWAIT")
	}
}
