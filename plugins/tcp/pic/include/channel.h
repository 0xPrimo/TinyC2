#ifndef _CHANNEL_H_
#define _CHANNEL_H_

#include "imports.h"
#include "tcg.h"

typedef struct {
	WSADATA WSAData;
    SOCKET 	Socket;
} CHANNEL_CONTEXT;

// IChannel channel interface
//
typedef struct {
	DWORD               ID;
	CHANNEL_CONTEXT*    Context;

	/**
     * @brief used by the implant when registring or switching channels.
     * @param Context channel context.
     * @return TRUE if the operation succeeds, otherwise FALSE.
     */
	BOOL(*Initialize)(CHANNEL_CONTEXT* Context);
	
	/**
     * @brief used by the implant to send data to the listener.
     * @param Context 	channel context.
	 * @param Data 		data implant want to send
	 * @param Size 		data size
	 * @param Register 	TRUE if this a registration request, otherwise FALSE.
     * @return TRUE if the operation succeeds, otherwise FALSE.
     */
	BOOL(*Send)(CHANNEL_CONTEXT* Context, CONST CHAR* Data, DWORD Size, BOOL Register);
	
	/**
     * @brief used by the implant to receive data from the listener.
     * @param Context 	channel context.
	 * @param Data 		received data
	 * @param Size 		received data size
     * @return TRUE if the operation succeeds, otherwise FALSE.
     */
	BOOL(*Receive)(CHANNEL_CONTEXT* Context, CHAR** Data, DWORD* Size);

	/**
     * @brief used by the implant when removing or switching channels.
     * @param Context channel context.
     * @return TRUE if the operation succeeds, otherwise FALSE.
     */
	BOOL(*Cleanup)(CHANNEL_CONTEXT* Context);
} IChannel;

#endif // _CHANNEL_H_
