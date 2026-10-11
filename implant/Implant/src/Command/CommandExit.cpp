#include "Command.h"

// CommandExit
//
BOOL CommandExit(json& args, string artifact, json& result) {
    ExitProcess(0);
	return TRUE;
}
