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
