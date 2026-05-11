package cli

import (
	"strings"

	"github.com/0xPrimo/TinyC2/server/internal/pkg/pack"
	lua "github.com/yuin/gopher-lua"
)

func (c *Cli) luaRegisterCommand(L *lua.LState) int {
	var script_name string

	if dbg, ok := L.GetStack(1); ok {
		L.GetInfo("S", dbg, lua.LNil)
		script_name = strings.TrimPrefix(dbg.Source, "@")
	}

	name := L.CheckString(1)
	callback := L.CheckFunction(2)
	desc := L.OptString(3, "")

	c.UserCommands[script_name] = append(c.UserCommands[script_name], UserCommand{
		Name:        name,
		Callback:    callback,
		Description: desc,
	})

	return 0
}

func (c *Cli) luaExecuteAssembly(L *lua.LState) int {
	id := uint32(L.CheckNumber(1))
	dotnet := L.CheckString(2)
	args := L.CheckString(3)

	c.Engine.ImplantExecuteAssembly(id, dotnet, args)

	L.Push(lua.LBool(true))
	return 1
}

func (c *Cli) luaInlineExecute(L *lua.LState) int {
	id := uint32(L.CheckNumber(1))
	bofpath := L.CheckString(2)
	bofargs := []byte(L.OptString(4, ""))

	c.Engine.ImplantInlineExecute(id, bofpath, bofargs)

	L.Push(lua.LBool(true))
	return 1
}

func (c *Cli) luaBofPack(l *lua.LState) int {
	var stringArgs []string

	formatStr := c.L.CheckString(1)
	numArgs := c.L.GetTop()

	for i := 2; i <= numArgs; i++ {
		lval := c.L.Get(i)
		stringArgs = append(stringArgs, lval.String())
	}

	packedBytes, _ := pack.BofPack(formatStr, stringArgs)
	c.L.Push(lua.LString(string(packedBytes)))
	return 1
}
