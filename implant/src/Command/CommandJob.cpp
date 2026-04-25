#include "Command.h"
#include "Implant.h"

// CommandJobStop stop running job
//
BOOL CommandJobStop(json& args, string artifact, json& result) {
	DWORD JobID = 0;

	JobID = args[0].get<DWORD>();
	
	LIST_ENTRY* current = g_Implant.JobList.Flink;

	while (current != &g_Implant.JobList) {
		PJOB Job = CONTAINING_RECORD(current, JOB, ListEntry);
        
        if (Job->ID == JobID) {
            TerminateProcess(Job->hProcess, 0);
            CloseHandle(Job->hProcess);
            CloseHandle(Job->hAnonPipe);
            RemoveEntryList(&Job->ListEntry);
            HeapFree(GetProcessHeap(), 0, Job);

            result["output"] = "[+] Job stoped and removed";
            return TRUE;
        }
		
        current = current->Flink;
	}

    result["output"] = "[-] Job not found";

	return FALSE;
}

// CommandJobList list running jobs
//
BOOL CommandJobList(json& args, string artifact, json& result) {
	LIST_ENTRY* current = g_Implant.JobList.Flink;
    json        jblist  = json::array();

	while (current != &g_Implant.JobList) {
        PJOB Job    = CONTAINING_RECORD(current, JOB, ListEntry);
        json JsnJob = json::object();
        
        JsnJob["id"] = Job->ID;
        jblist.push_back(JsnJob);

        current = current->Flink;
	}

    result["artifact"]  = jblist.dump();

	return FALSE;
}