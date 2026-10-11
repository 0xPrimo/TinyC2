package cmd

import (
	"strconv"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func (c *Cli) luaInitialize() {
	c.L.SetGlobal("register_command", c.L.NewFunction(c.luaRegisterCommand))

	apiTable := c.L.NewTable()
	c.L.SetField(apiTable, "inline_execute", c.L.NewFunction(c.luaInlineExecute))
	c.L.SetField(apiTable, "execute_assembly", c.L.NewFunction(c.luaExecuteAssembly))
	c.L.SetGlobal("api", apiTable)
}

func (c *Cli) luaRegisterCommand(L *lua.LState) int {
	var script string

	if dbg, ok := L.GetStack(1); ok {
		L.GetInfo("S", dbg, lua.LNil)
		script = strings.TrimPrefix(dbg.Source, "@")
	}

	name := L.CheckString(1)
	callback := L.CheckFunction(2)
	desc := L.OptString(3, "")

	c.ScriptedCommands[script] = append(c.ScriptedCommands[script], UserCommand{
		Name:        name,
		Callback:    callback,
		Description: desc,
	})

	return 0
}

func (c *Cli) luaExecuteAssembly(L *lua.LState) int {
	return 1
}

func (c *Cli) luaInlineExecute(L *lua.LState) int {
	top := L.GetTop()

	if top < 3 {
		L.Push(lua.LBool(false))
		L.Push(lua.LString("expected at least 3 arguments"))
		return 2
	}

	id := uint32(L.CheckNumber(1))
	bofpath := L.CheckString(2)
	pack := L.CheckString(3)

	var args []string

	args = append(args, bofpath)
	args = append(args, pack)
	for i := 4; i <= top; i++ {
		args = append(args, L.ToStringMeta(L.Get(i)).String())
	}

	err := c.Engine.ImplantExecute(strconv.FormatUint(uint64(id), 16), "inline-execute", args...)
	if err != nil {
		L.Push(lua.LBool(false))
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LBool(true))
	return 1
}
