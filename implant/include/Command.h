#pragma once
#define _CRT_SECURE_NO_WARNINGS

#include "Ntdll.h"
#include <vector>
#include <string>
#include <nlohmann/json.hpp>

using json = nlohmann::json;
using string = std::string;

// ImplantCommand implant command handler
//
// g_ImplantCommandRegistry implant command handlers
//
typedef struct {
	string Name;
	BOOL(*Invoke)(json&, string, json&);
} ImplantCommand;

extern std::vector<ImplantCommand>	g_CommandRegistry;