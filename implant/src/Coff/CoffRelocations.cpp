#include "Coff.h"
#include <stdio.h>

BOOL Handle_AMD64_ADDR64    (COFF_CONTEXT* Ctx, DWORD SectionIndex, PCOFF_RELOCATION RelocationHeader, PVOID SymAddress);
BOOL Handle_AMD64_ADDR32NB  (COFF_CONTEXT* Ctx, DWORD SectionIndex, PCOFF_RELOCATION RelocationHeader, PVOID SymAddress);
BOOL Handle_AMD64_REL32     (COFF_CONTEXT* Ctx, DWORD SectionIndex, PCOFF_RELOCATION RelocationHeader, PVOID SymAddress);
BOOL Handle_AMD64_REL32_1   (COFF_CONTEXT* Ctx, DWORD SectionIndex, PCOFF_RELOCATION RelocationHeader, PVOID SymAddress);
BOOL Handle_AMD64_REL32_2   (COFF_CONTEXT* Ctx, DWORD SectionIndex, PCOFF_RELOCATION RelocationHeader, PVOID SymAddress);
BOOL Handle_AMD64_REL32_3   (COFF_CONTEXT* Ctx, DWORD SectionIndex, PCOFF_RELOCATION RelocationHeader, PVOID SymAddress);
BOOL Handle_AMD64_REL32_4   (COFF_CONTEXT* Ctx, DWORD SectionIndex, PCOFF_RELOCATION RelocationHeader, PVOID SymAddress);
BOOL Handle_AMD64_REL32_5   (COFF_CONTEXT* Ctx, DWORD SectionIndex, PCOFF_RELOCATION RelocationHeader, PVOID SymAddress);


COFF_RELOC_HANDLER COFFRelocationHandlers[RELOC_HANDLERS_COUNT] = {
    { .Type = IMAGE_REL_AMD64_ADDR64,   .Handler = Handle_AMD64_ADDR64      }, 
    { .Type = IMAGE_REL_AMD64_ADDR32NB, .Handler = Handle_AMD64_ADDR32NB    }, 
    { .Type = IMAGE_REL_AMD64_REL32,    .Handler = Handle_AMD64_REL32       }, 
    { .Type = IMAGE_REL_AMD64_REL32_1,  .Handler = Handle_AMD64_REL32_1     }, 
    { .Type = IMAGE_REL_AMD64_REL32_2,  .Handler = Handle_AMD64_REL32_2     }, 
    { .Type = IMAGE_REL_AMD64_REL32_3,  .Handler = Handle_AMD64_REL32_3     }, 
    { .Type = IMAGE_REL_AMD64_REL32_4,  .Handler = Handle_AMD64_REL32_4     }, 
    { .Type = IMAGE_REL_AMD64_REL32_5,  .Handler = Handle_AMD64_REL32_5     }, 
};

BOOL 
Handle_AMD64_ADDR64(
    COFF_CONTEXT*           Ctx,
    DWORD                   SectionIndex,
    PCOFF_RELOCATION        RelocationHeader, 
    PVOID                   SymAddress
) {
    UINT64                  Offset          = 0;
    PCOFF_SECTION_HEADER    Section         = Ctx->SectionHeader + SectionIndex;
    PBYTE                   SectionBuffer   = (PBYTE)(Ctx->SectionBuffer[SectionIndex]);

    memcpy(
        &Offset,
        SectionBuffer + RelocationHeader->VirtualAddress,
        sizeof(UINT64)
    );
    
    Offset += (UINT64)SymAddress;
    
    memcpy(
        SectionBuffer + RelocationHeader->VirtualAddress,
        &Offset,
        sizeof(UINT64)
    );
    
    return TRUE;
}


BOOL Handle_AMD64_ADDR32NB(
    COFF_CONTEXT*           Ctx,
    DWORD                   SectionIndex,
    PCOFF_RELOCATION        RelocationHeader, 
    PVOID                   SymAddress
) {
    UINT32                  Offset          = 0;
    PCOFF_SECTION_HEADER    Section         = Ctx->SectionHeader + SectionIndex;
    PBYTE                   SectionBuffer   = (PBYTE)(Ctx->SectionBuffer[SectionIndex]);
    PBYTE                   JumpSection     = NULL;

    memcpy(
        &Offset,
        SectionBuffer + RelocationHeader->VirtualAddress, 
        sizeof(UINT32)
    );

    JumpSection = (PBYTE)(Ctx->SectionBuffer[(Ctx->SymbolTable[RelocationHeader->SymbolTableIndex].SectionNumber - 1)]);
    if (!IS_4GB_AWAY(JumpSection, SectionBuffer + RelocationHeader->VirtualAddress + 4)) {
        return FALSE;
    }
    
    Offset = (DWORD_PTR)(JumpSection + Offset ) - (DWORD_PTR)(SectionBuffer + RelocationHeader->VirtualAddress + 4);
    Offset += Ctx->SymbolTable[RelocationHeader->SymbolTableIndex].Value;

    memcpy(
        SectionBuffer + RelocationHeader->VirtualAddress,
        &Offset,
        sizeof(UINT32)
    );

    return TRUE;
}


BOOL Handle_AMD64_REL32(
    COFF_CONTEXT*           Ctx,
    DWORD                   SectionIndex,
    PCOFF_RELOCATION        RelocationHeader, 
    PVOID                   SymAddress
) {
    UINT32                  Offset          = 0;
    PCOFF_SECTION_HEADER    Section         = Ctx->SectionHeader + SectionIndex;
    PBYTE                   SectionBuffer   = (PBYTE)(Ctx->SectionBuffer[SectionIndex]);
    PBYTE                   JumpSection     = NULL;

    memcpy(&Offset, SectionBuffer + RelocationHeader->VirtualAddress, sizeof(UINT32));
    
    if (llabs((DWORD_PTR)SymAddress - (DWORD_PTR)(SectionBuffer + RelocationHeader->VirtualAddress + 4)) > UINT_MAX) {
        printf("[+] address 4gb away\n");
        return FALSE;
    }
    

    Offset += ((size_t)SymAddress - ((size_t)SectionBuffer + RelocationHeader->VirtualAddress + 4));
    
    memcpy(SectionBuffer + RelocationHeader->VirtualAddress, &Offset, sizeof(UINT32));
    return TRUE;
}

BOOL Handle_AMD64_REL32_1(
    COFF_CONTEXT*           Ctx,
    DWORD                   SectionIndex,
    PCOFF_RELOCATION        RelocationHeader, 
    PVOID                   SymAddress
) {
    UINT32                  Offset = 0;
    PCOFF_SECTION_HEADER    Section         = Ctx->SectionHeader + SectionIndex;
    PBYTE                   SectionBuffer   = (PBYTE)(Ctx->SectionBuffer[SectionIndex]);
    PBYTE                   JumpSection     = NULL;


    memcpy(
        &Offset,
        SectionBuffer + RelocationHeader->VirtualAddress,
        sizeof(UINT32)
    );

    if (llabs((DWORD_PTR)SymAddress - (DWORD_PTR)(SectionBuffer + RelocationHeader->VirtualAddress + 4 + 1)) > UINT_MAX) {
        return FALSE;
    }
    
    Offset += ((size_t)SymAddress - ((size_t)SectionBuffer + RelocationHeader->VirtualAddress + 4 + 1));
    
    memcpy(
        SectionBuffer + RelocationHeader->VirtualAddress,
        &Offset,
        sizeof(UINT32)
    );

    return TRUE;
}

BOOL Handle_AMD64_REL32_2(
    COFF_CONTEXT*           Ctx,
    DWORD                   SectionIndex,
    PCOFF_RELOCATION        RelocationHeader, 
    PVOID                   SymAddress
) {
    UINT32                  Offset = 0;
    PCOFF_SECTION_HEADER    Section         = Ctx->SectionHeader + SectionIndex;
    PBYTE                   SectionBuffer   = (PBYTE)(Ctx->SectionBuffer[SectionIndex]);
    PBYTE                   JumpSection     = NULL;


    memcpy(
        &Offset,
        SectionBuffer + RelocationHeader->VirtualAddress,
        sizeof(UINT32)
    );

    if (llabs((DWORD_PTR)SymAddress - (DWORD_PTR)(SectionBuffer + RelocationHeader->VirtualAddress + 4 + 2)) > UINT_MAX) {
        return FALSE;
    }

    Offset += ((size_t)SymAddress - ((size_t)SectionBuffer + RelocationHeader->VirtualAddress + 4 + 2));
    
    memcpy(
        SectionBuffer + RelocationHeader->VirtualAddress,
        &Offset,
        sizeof(UINT32)
    );

    return TRUE;
}

BOOL Handle_AMD64_REL32_3(
    COFF_CONTEXT*           Ctx,
    DWORD                   SectionIndex,
    PCOFF_RELOCATION        RelocationHeader, 
    PVOID                   SymAddress
) {
    UINT32                  Offset = 0;
    PCOFF_SECTION_HEADER    Section         = Ctx->SectionHeader + SectionIndex;
    PBYTE                   SectionBuffer   = (PBYTE)(Ctx->SectionBuffer[SectionIndex]);
    PBYTE                   JumpSection     = NULL;


    memcpy(
        &Offset,
        SectionBuffer + RelocationHeader->VirtualAddress,
        sizeof(UINT32)
    );

    if (llabs((DWORD_PTR)SymAddress - (DWORD_PTR)(SectionBuffer + RelocationHeader->VirtualAddress + 4 + 3)) > UINT_MAX) {
        return FALSE;
    }

    Offset += ((size_t)SymAddress - ((size_t)SectionBuffer + RelocationHeader->VirtualAddress + 4 + 3));
    
    memcpy(
        SectionBuffer + RelocationHeader->VirtualAddress,
        &Offset,
        sizeof(UINT32)
    );

    return TRUE;
}

BOOL Handle_AMD64_REL32_4(
    COFF_CONTEXT*           Ctx,
    DWORD                   SectionIndex,
    PCOFF_RELOCATION        RelocationHeader, 
    PVOID                   SymAddress
) {
    UINT32                  Offset = 0;
    PCOFF_SECTION_HEADER    Section         = Ctx->SectionHeader + SectionIndex;
    PBYTE                   SectionBuffer   = (PBYTE)(Ctx->SectionBuffer[SectionIndex]);
    PBYTE                   JumpSection     = NULL;


    memcpy(
        &Offset,
        SectionBuffer + RelocationHeader->VirtualAddress,
        sizeof(UINT32)
    );

    if (llabs((DWORD_PTR)SymAddress - (DWORD_PTR)(SectionBuffer + RelocationHeader->VirtualAddress + 4 + 4)) > UINT_MAX) {
        return FALSE;
    }

    Offset += ((size_t)SymAddress - ((size_t)SectionBuffer + RelocationHeader->VirtualAddress + 4 + 4));
    
    memcpy(
        SectionBuffer + RelocationHeader->VirtualAddress, 
        &Offset,
        sizeof(UINT32)
    );

    return TRUE;
}

BOOL Handle_AMD64_REL32_5(
    COFF_CONTEXT*           Ctx,
    DWORD                   SectionIndex,
    PCOFF_RELOCATION        RelocationHeader, 
    PVOID                   SymAddress
) {
    UINT32                  Offset = 0;
    PCOFF_SECTION_HEADER    Section         = Ctx->SectionHeader + SectionIndex;
    PBYTE                   SectionBuffer   = (PBYTE)(Ctx->SectionBuffer[SectionIndex]);
    PBYTE                   JumpSection     = NULL;


    memcpy(
        &Offset,
        SectionBuffer + RelocationHeader->VirtualAddress,
        sizeof(UINT32)
    );

    if (llabs((DWORD_PTR)SymAddress - (DWORD_PTR)(SectionBuffer + RelocationHeader->VirtualAddress + 4 + 5)) > UINT_MAX) {
        return FALSE;
    }

    Offset += ((size_t)SymAddress - ((size_t)SectionBuffer + RelocationHeader->VirtualAddress + 4 + 5));
    
    memcpy(
        SectionBuffer + RelocationHeader->VirtualAddress,
        &Offset, 
        sizeof(UINT32)
    );

    return TRUE;
}
