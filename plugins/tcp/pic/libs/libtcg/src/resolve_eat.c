/*
 * Copyright 2025 Raphael Mudge, Adversary Fan Fiction Writers Guild
 *
 * Redistribution and use in source and binary forms, with or without modification, are
 * permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this list of
 * conditions and the following disclaimer.
 *
 * 2. Redistributions in binary form must reproduce the above copyright notice, this list of
 * conditions and the following disclaimer in the documentation and/or other materials provided
 * with the distribution.
 *
 * 3. Neither the name of the copyright holder nor the names of its contributors may be used to
 * endorse or promote products derived from this software without specific prior written
 * permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS “AS IS” AND ANY EXPRESS
 * OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE
 * COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL,
 * EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR
 * TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE,
 * EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

#include "loaderdefs.h"
#include "tcg.h"

/*
 * Walk the Export Address Table to resolve functions by hash
 */

HANDLE findModuleByHash(DWORD moduleHash) {
	_PEB                 * pPEB;
	LDR_DATA_TABLE_ENTRY * pEntry;
	char                 * name;
	DWORD                  hashValue;
	USHORT                 counter;

	/* get the Process Enviroment Block */
#if defined WIN_X64
	pPEB = (_PEB *)__readgsqword( 0x60 );
#elif defined WIN_X86
	pPEB = (_PEB *)__readfsdword( 0x30 );
#else
#error "Neither WIN_X64 or WIN_X86 is defined"
#endif

	/* walk the module list */
	pEntry = (LDR_DATA_TABLE_ENTRY *)pPEB->pLdr->InMemoryOrderModuleList.Flink;

	while (pEntry) {
		/* pEntry->BaseDllName is a UNICODE_STRING, pBuffer is wchar_t*, and Length is IN bytes.
		   We are walking and hashing this string, one byte at a time */
		name      = (char *)pEntry->BaseDllName.pBuffer;
		counter   = pEntry->BaseDllName.Length;

		/* calculate the hash of our DLL name */
		hashValue = 0;
		do {
			hashValue = ror(hashValue);
			if (*name >= 'a')
				hashValue += (BYTE)*name - 0x20;
			else
				hashValue += (BYTE)*name;

			name++;
		} while (--counter);

		/* if we have a match, return it */
		if (hashValue == moduleHash)
			return (HANDLE)pEntry->DllBase;

		/* next entry */
		pEntry = (LDR_DATA_TABLE_ENTRY *)pEntry->InMemoryOrderModuleList.Flink;
	}

	return NULL;
}

FARPROC findFunctionByHash(HANDLE src, DWORD wantedFunction) {
	DLLDATA                  data;
	IMAGE_DATA_DIRECTORY   * exportTableHdr;
	IMAGE_EXPORT_DIRECTORY * exportDir;
	DWORD                  * exportName;
	WORD                   * exportOrdinal;
	DWORD                  * exportAddress;
	DWORD                    hashValue;

	/* parse our DLL! */
	ParseDLL(src, &data);

	/* grab our export directory */
	exportTableHdr = GetDataDirectory(&data, IMAGE_DIRECTORY_ENTRY_EXPORT);
	exportDir      = (IMAGE_EXPORT_DIRECTORY *)PTR_OFFSET(src, exportTableHdr->VirtualAddress);

	/* walk the array of exported names/address ordinals */
	exportName    = (DWORD *)PTR_OFFSET(src, exportDir->AddressOfNames);
	exportOrdinal = (WORD *) PTR_OFFSET(src, exportDir->AddressOfNameOrdinals);

	while (TRUE) {
		hashValue = hash( (char *)PTR_OFFSET(src, *exportName) );
		if (hashValue == wantedFunction) {
			/* figure out the base of our AddressOfFunctions array */
			exportAddress   = PTR_OFFSET(src, exportDir->AddressOfFunctions);

			/* increment it by the current value of our exportOrdinal array */
			exportAddress  += *exportOrdinal;

			/* and... there-in is our virtual address to the actual ptr we want */
			return (FARPROC)PTR_OFFSET(src, *exportAddress);
		}

		exportName++;
		exportOrdinal++;
	}
}
