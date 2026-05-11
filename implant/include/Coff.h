#pragma once

#include <windows.h>
#include <stdint.h>

// COFF structures
//

#define IS_SYMBOL_DEFINED(x)    (((PCOFF_SYMBOL)x)->SectionNumber > 0)
#define IS_SYMBOL_EXTERNAL(x)   ( ((PCOFF_SYMBOL)x)->StorageClass == IMAGE_SYM_CLASS_EXTERNAL \
                                    || ((PCOFF_SYMBOL)x)->StorageClass == IMAGE_SYM_CLASS_EXTERNAL_DEF ) \
                                    && ( ((PCOFF_SYMBOL)x)->Value == 0 )
#define IS_4GB_AWAY(s,d)        ((DWORD_PTR)(s) - (DWORD_PTR)(d) < 0xffffffff)
#define RELOC_HANDLERS_COUNT    8


typedef struct _COFF_FILE_HEADER 
{
    UINT16  Machine;
    UINT16  NumberOfSections;
    UINT32  TimeDateStamp;
    UINT32  PointerToSymbolTable;
    UINT32  NumberOfSymbols;
    UINT16  SizeOfOptionalHeader;
    UINT16  Characteristics;
} COFF_FILE_HEADER, *PCOFF_FILE_HEADER;

#pragma pack(push,1)

typedef struct _COFF_SECTION
{
    CHAR    Name[ 8 ];
    UINT32  VirtualSize;
    UINT32  VirtualAddress;
    UINT32  SizeOfRawData;
    UINT32  PointerToRawData;
    UINT32  PointerToRelocations;
    UINT32  PointerToLineNumbers;
    UINT16  NumberOfRelocations;
    UINT16  NumberOfLinenumbers;
    UINT32  Characteristics;
} COFF_SECTION_HEADER, *PCOFF_SECTION_HEADER;


typedef struct _COFF_RELOC
{
    UINT32  VirtualAddress;
    UINT32  SymbolTableIndex;
    UINT16  Type;
} COFF_RELOCATION, *PCOFF_RELOCATION;

typedef struct _COFF_SYMBOL
{
    union
    {
        CHAR    Name[8];
        UINT32  Value[2];
    } First;

    UINT32 Value;
    UINT16 SectionNumber;
    UINT16 Type;
    UINT8  StorageClass;
    UINT8  NumberOfAuxSymbols;
} COFF_SYMBOL, *PCOFF_SYMBOL;

#pragma pack(pop)

typedef struct {
    PCOFF_SECTION_HEADER    Header;
    PCOFF_RELOCATION        Relocations;
    PBYTE                   Base;
    DWORD                   Size;
    DWORD                   Protection;
} COFF_SECTION_DATA, *PCOFF_SECTION_DATA;

typedef struct {
    PBYTE                   FileBuffer;
    DWORD                   FileSize;

    PCOFF_FILE_HEADER       FileHeader;
    PCOFF_SECTION_HEADER    SectionHeader;
    PCOFF_RELOCATION        RelocationHeader;
    PVOID*                  SectionBuffer;
    PCOFF_SYMBOL            SymbolTable;
    CHAR*                   StringTable;
    DWORD_PTR*              GlobalOffsetTable;
    DWORD                   GlobalOffsetTableIndex;
} COFF_CONTEXT, *PCOFF_CONTEXT;

typedef DWORD(*COFF_ENTRYPOINT)(PVOID, DWORD);

typedef struct {
    PCOFF_CONTEXT   Context;
    COFF_ENTRYPOINT Entrypoint;
    PVOID           ArgsData;
    DWORD           ArgsSize;
} COFF_ENTRYPOINT_DATA, *PCOFF_ENTRYPOINT_DATA;

typedef DWORD(*COFF_START_ROUTINE)(COFF_ENTRYPOINT_DATA);

typedef struct {
    DWORD   Type;
    BOOL    (*Handler)(
        COFF_CONTEXT*       Ctx,
        DWORD               SectionIndex,
        PCOFF_RELOCATION    Relocation,
        PVOID               SymAddress
    );
} COFF_RELOC_HANDLER, *PCOFF_RELOC_HANDLER;

extern COFF_RELOC_HANDLER COFFRelocationHandlers[RELOC_HANDLERS_COUNT];

BOOL COFFLoader(
    BOOL            Sync,
    PVOID           FileData,
    DWORD           FileSize,
    PVOID           ArgsData,
    DWORD           ArgsSize
);

BOOL COFFCleanup(
    PCOFF_CONTEXT* Ctx
);

// Beacon comptability
//
typedef struct {
    CHAR* Name;
    PVOID Address;
} BEACON_API, *PBEACON_API;

 /* Structures as is in beacon.h */
extern BEACON_API   InternalFunctions[30];
extern char*        beacon_compatibility_output;
extern int          beacon_compatibility_size;
extern int          beacon_compatibility_offset;

typedef struct {
    char * original; /* the original buffer [so we can free it] */
    char * buffer;   /* current pointer into our buffer */
    int    length;   /* remaining length of data */
    int    size;     /* total size of this buffer */
} datap;

typedef struct {
    char * original; /* the original buffer [so we can free it] */
    char * buffer;   /* current pointer into our buffer */
    int    length;   /* remaining length of data */
    int    size;     /* total size of this buffer */
} formatp;

void    BeaconDataParse(datap * parser, char * buffer, int size);
int     BeaconDataInt(datap * parser);
short   BeaconDataShort(datap * parser);
int     BeaconDataLength(datap * parser);
char *  BeaconDataExtract(datap * parser, int * size);

void    BeaconFormatAlloc(formatp * format, int maxsz);
void    BeaconFormatReset(formatp * format);
void    BeaconFormatFree(formatp * format);
void    BeaconFormatAppend(formatp * format, char * text, int len);
void    BeaconFormatPrintf(formatp * format, char * fmt, ...);
char *  BeaconFormatToString(formatp * format, int * size);
void    BeaconFormatInt(formatp * format, int value);

#define CALLBACK_OUTPUT      0x0
#define CALLBACK_OUTPUT_OEM  0x1e
#define CALLBACK_ERROR       0x0d
#define CALLBACK_OUTPUT_UTF8 0x20


void   BeaconPrintf(int type, char * fmt, ...);
void   BeaconOutput(int type, char * data, int len);

/* Token Functions */
BOOL   BeaconUseToken(HANDLE token);
void   BeaconRevertToken();
BOOL   BeaconIsAdmin();

/* Spawn+Inject Functions */
void   BeaconGetSpawnTo(BOOL x86, char * buffer, int length);
BOOL BeaconSpawnTemporaryProcess(BOOL x86, BOOL ignoreToken, STARTUPINFO * sInfo, PROCESS_INFORMATION * pInfo);
void   BeaconInjectProcess(HANDLE hProc, int pid, char * payload, int p_len, int p_offset, char * arg, int a_len);
void   BeaconInjectTemporaryProcess(PROCESS_INFORMATION * pInfo, char * payload, int p_len, int p_offset, char * arg, int a_len);
void   BeaconCleanupProcess(PROCESS_INFORMATION * pInfo);

/* Utility Functions */
BOOL   toWideChar(char * src, wchar_t * dst, int max);
uint32_t swap_endianess(uint32_t indata);

char* BeaconGetOutputData(int *outsize);
