#include "Command.h"
#include "Rtlib.h"
// #include "Peer.h"
#include "Debug.h"

BOOL SmbConnect( const CHAR* PipeName, json& PivotResponse );

BOOL CommandPivotConnect( json& args, string artifact, json& result ) {
    //DWORD ConfigSize = 0;
    //CHAR* PipeName   = (CHAR*)Base64Decode( artifact.c_str(), &ConfigSize );

    //json PivotResponse;
    //if (!SmbConnect( PipeName, PivotResponse )) {
    //    printf( "[-] Failed to connect to peer implant\n" );
    //    return FALSE;
    //}

    //result["name"]     = "pivot";
    //result["type"]     = "connect";
    //result["artifact"] = PivotResponse.dump().c_str();

    return TRUE;
}

BOOL SmbWrite( DWORD ID, CHAR* Data, DWORD Size );

BOOL CommandPivotRequest( json& args, string artifact, json& result ) {
    //DWORD Size      = 0;
    //PBYTE Data      = NULL;
    //DWORD ImplantID = args[0].get<DWORD>();
    //json  Request;

    //Data = Base64Decode( artifact.c_str(), &Size );
    //if (Data == NULL) {
    //    printf( "failed to decode artifact" );
    //    return FALSE;
    //}

    //if (!SmbWrite( ImplantID, (CHAR*)Data, Size )) {
    //    printf( "failed to write to pivot" );
    //    return FALSE;
    //}

    return TRUE;
}
