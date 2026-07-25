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

#include <windows.h>
#include "tcg.h"

WINBASEAPI HMODULE WINAPI KERNEL32$GetModuleHandleA (LPCSTR lpModuleName);

/*
 * This is our opt-in Dynamic Function Resolution resolver. It turns MODULE$Function into pointers.
 * See dfr "resolve" "ror13" "KERNEL32, NTDLL" in loader.spec
 */

FARPROC resolve(DWORD modHash, DWORD funcHash) {
	HANDLE hModule = findModuleByHash(modHash);
	return findFunctionByHash(hModule, funcHash);
}

/*
 * This is our default DFR resolver. It resolves Win32 APIs not handled by another resolver.
 */

FARPROC resolve_ext(char * mod, char * func) {
	HANDLE hModule = KERNEL32$GetModuleHandleA(mod);
	if (hModule == NULL)
		hModule = LoadLibraryA(mod);

	return GetProcAddress(hModule, func);
}

/*
 * This is our opt-in function to help fix ptrs in x86 PIC. See fixptrs _caller" in loader.spec
 */

#ifdef WIN_X86
__declspec(noinline) ULONG_PTR caller( VOID ) { return (ULONG_PTR)WIN_GET_CALLER(); }
#endif
