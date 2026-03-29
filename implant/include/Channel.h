#pragma once

#include "phnt.h"

// IChannel channel interface
//
typedef struct {
	DWORD ID;
	PVOID Context;

	/**
     * @brief used by the implant when registring or switching channels.
     * @param Context channel context.
     * @return TRUE if the operation succeeds, otherwise FALSE.
     */
	BOOL(*Initialize)(PVOID Context);
	
	/**
     * @brief used by the implant to send data to the listener.
     * @param Context 	channel context.
	 * @param Data 		data implant want to send
	 * @param Size 		data size
	 * @param Register 	TRUE if this a registration request, otherwise FALSE.
     * @return TRUE if the operation succeeds, otherwise FALSE.
     */
	BOOL(*Send)(PVOID Context, CONST CHAR* Data, DWORD Size, BOOL Register);
	
	/**
     * @brief used by the implant to receive data from the listener.
     * @param Context 	channel context.
	 * @param Data 		received data
	 * @param Size 		received data size
     * @return TRUE if the operation succeeds, otherwise FALSE.
     */
	BOOL(*Receive)(PVOID Context, CHAR** Data, DWORD* Size);

	/**
     * @brief used by the implant when removing or switching channels.
     * @param Context channel context.
     * @return TRUE if the operation succeeds, otherwise FALSE.
     */
	BOOL(*Cleanup)(PVOID Context);
} IChannel;

// CHANNEL channel node
//
typedef struct {
	LIST_ENTRY	ListEntry;
	IChannel*	Interface;
	struct {
		PVOID Base;
		DWORD Size;
	} Memory;
} CHANNEL, *PCHANNEL;

BOOL ChannelInitialize();
BOOL ChannelLoad(PVOID Destination, PVOID Source, DWORD Size);
BOOL ChannelRegister(PVOID BaseAddr, DWORD Size);
BOOL ChannelSwitch(DWORD ID);
BOOL ChannelRemove(DWORD ID);
VOID ChannelFree(PCHANNEL Channel);

extern LIST_ENTRY	g_ChannelList;
extern PCHANNEL		g_Channel;
