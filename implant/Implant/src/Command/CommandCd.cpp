#include "Command.h"

BOOL CommandCd(json& args, string artifact, json& result) {
    result["name"] = "cd";

    auto directory = args[0].get<string>();

    if (SetCurrentDirectoryA(directory.c_str())) {
        result["output"] = "[+] directory changed";
        return TRUE;
    } 
    
    result["output"] = "[-] failed to change directory";
    return FALSE;
}