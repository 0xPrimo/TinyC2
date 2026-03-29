#include "imports.h"
#include "tcg.h"

#define DBG_PRINTF(x, ...) dprintf("PROTOCOL: "x, ##__VA_ARGS__)

typedef struct {
    CHAR* Data;
    DWORD Size;
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

BOOL ProtocolSend(CHANNEL_CONTEXT* Context, CONST CHAR* Data, DWORD Size, BOOL Register) {
	return TRUE;
}

BOOL ProtocolReceive(CHANNEL_CONTEXT* Context, CHAR** Data, DWORD* Size) {
	return TRUE;
}

BOOL ProtocolInitialize(CHANNEL_CONTEXT* Context) {
	return TRUE;
}

BOOL ProtocolCleanup(CHANNEL_CONTEXT* Context) {
	return TRUE;
}

BOOL go(void* implant, IChannel* channel) {

	channel->Context = RtlAllocateHeap(GetProcessHeap(), HEAP_ZERO_MEMORY, sizeof(CHANNEL_CONTEXT));
	if (channel->Context == NULL) {
		return FALSE;
	}

	channel->ID 		= 0xbaadf00d;
	channel->Initialize = ProtocolInitialize;
	channel->Send 		= ProtocolSend;
	channel->Receive 	= ProtocolReceive;
	channel->Cleanup 	= ProtocolCleanup;
	
	return TRUE;
}
