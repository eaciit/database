package base

import (
	"github.com/eaciit/errorlib"
	_ "strings"
)

type IAdapter interface {
	SetCommand(DB_OP, ICommand)
	Command(DB_OP) ICommand
	Run(DB_OP, interface{}, map[string]interface{}) (ICursor, int, error)
}

type AdapterType int

const (
	AdapterAuto AdapterType = iota
	AdapaterManual
)

type AdapterBase struct {
	Connection    IConnection
	Type          AdapterType
	TableName     string
	SelectCommand ICommand
	InsertCommand ICommand
	UpdateCommand ICommand
	DeleteCommand ICommand
	SaveCommand   ICommand
}

func (a *AdapterBase) SetCommand(commandType DB_OP, command ICommand) {
	command.Prop("Type", commandType)
	if commandType == DB_INSERT {
		a.InsertCommand = command
	} else if commandType == DB_UPDATE {
		a.UpdateCommand = command
	} else if commandType == DB_DELETE {
		a.DeleteCommand = command
	} else if commandType == DB_SELECT {
		a.SelectCommand = command
	} else if commandType == DB_SAVE {
		a.SaveCommand = command
	}
}

func (a *AdapterBase) Command(commandType DB_OP) ICommand {
	if commandType == DB_INSERT {
		return a.InsertCommand
	} else if commandType == DB_UPDATE {
		return a.UpdateCommand
	} else if commandType == DB_DELETE {
		return a.DeleteCommand
	} else if commandType == DB_SELECT {
		return a.SelectCommand
	} else if commandType == DB_SAVE {
		return a.SaveCommand
	}
	return nil
}

func (a *AdapterBase) Run(commandType DB_OP, result interface{}, parms map[string]interface{}) (ICursor, int, error) {
	cmd := a.Command(commandType)
	if cmd == nil {
		return nil, 0, errorlib.Error(packageName, modAdapter, "Run", "Command "+string(commandType)+" is not initialized")
	}
	return cmd.Run(result, parms)
}
