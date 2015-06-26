package base

const (
	packageName   = "database/base"
	modAdapter    = "AdapterBase"
	modCommand    = "CommandBase"
	modConnection = "ConnectionBase"
	modQuery      = "QueryBase"
	modCursor     = "CursorBase"
)

const (
	DB_INSERT = "insert"
	DB_UPDATE = "update"
	DB_DELETE = "delete"
	DB_SELECT = "select"
	DB_SAVE   = "save"
)

type M map[string]interface{}
