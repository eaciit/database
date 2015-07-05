package base

import (
	//"github.com/eaciit/errorlib"
	"strings"
)

type ICommand interface {
	Prop(string, interface{})
	Run(interface{}, map[string]interface{}) (ICursor, int, error)
}

type CommandBase struct {
	Connection IConnection
	Type       string
	Text       string
	Settings   map[string]interface{}
}

func (c *CommandBase) Run(result interface{}, parms map[string]interface{}) (ICursor, int, error) {
	return nil, 0, nil
}

func (c *CommandBase) Prop(fieldname string, result interface{}) {
	fieldname = strings.ToLower(fieldname)
	if fieldname == "type" {
		c.Type = result.(string)
	} else if fieldname == "text" {
		c.Text = result.(string)
	} else if fieldname == "settings" {
		c.Settings = result.(map[string]interface{})
	}
}
