#include "Command.h"

BOOL CommandCp(json& args, string artifact, json& result) {
    auto src    = args[0].get<string>();
    auto dest   = args[1].get<string>();
    
    if (CopyFileA(src.c_str(), dest.c_str(), FALSE)) {
        result["output"] = "[+] file copied successfully";
        return TRUE;
    } else {
        result["output"] = "[-] failed to copy file";
        return FALSE;
    }
}