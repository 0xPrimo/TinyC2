rule LibTCG_DLLLoader_a3d8be97 {
	meta:
		description = "LibTCG: DLL loading and parsing functions derived from ReflectiveDllInjection"
		author = "Raphael Mudge"
		date = "2026-02-02"
		reference = "https://tradecraftgarden.org/libtcg.html"
		arch_context = "x64"
		scan_context = "file, memory"
		os = "windows"
		license = "BSD"
		generator = "Crystal Palace"
	strings:
		// ----------------------------------------
		// Function: ProcessRelocation
		// ----------------------------------------
		/*
		 * 49 01 C0                      add r8, rax
		 * 41 8B 41 04                   mov eax, dword ptr [r9+4]
		 * 48 83 E8 08                   sub rax, 8
		 * 48 D1 E8                      shr rax, 1
		 * 49 8D 51 08                   lea rdx, [r9+8]
		 * (Score: 56)
		 */
		$r0_ProcessRelocation = { 49 01 C0 41 8B 41 04 48 83 E8 08 48 D1 E8 49 8D 51 08 }

		/*
		 * 4D 89 D3                      mov r11, r10
		 * 49 C1 EB 10                   shr r11, 0x10
		 * 8D 40 FF                      lea eax, [rax-1]
		 * 49 8D 4C 41 0A                lea rcx, [r9+rax*2+0xA]
		 * (Score: 77)
		 */
		$r1_ProcessRelocation = { 4D 89 D3 49 C1 EB 10 8D 40 FF 49 8D 4C 41 0A }

		/*
		 * 0F B7 02                      movzx eax, word ptr [rdx]
		 * 25 FF 0F 00 00                and eax, 0xFFF
		 * 4D 01 14 00                   add qword ptr [r8+rax], r10
		 * (Score: 28)
		 */
		$r2_ProcessRelocation = { 0F B7 02 25 FF 0F 00 00 4D 01 14 00 }

		/*
		 * 0F B7 02                      movzx eax, word ptr [rdx]
		 * 25 FF 0F 00 00                and eax, 0xFFF
		 * 66 45 01 14 00                add word ptr [r8+rax], r10w
		 * (Score: 28)
		 */
		$r3_ProcessRelocation = { 0F B7 02 25 FF 0F 00 00 66 45 01 14 00 }

		// ----------------------------------------
		// Function: ParseDLL
		// ----------------------------------------
		/*
		 * 48 01 C1                      add rcx, rax
		 * 48 89 4A 08                   mov qword ptr [rdx+8], rcx
		 * 48 83 C1 18                   add rcx, 0x18
		 * 48 89 4A 10                   mov qword ptr [rdx+0x10], rcx
		 * (Score: 74)
		 */
		$r4_ParseDLL = { 48 01 C1 48 89 4A 08 48 83 C1 18 48 89 4A 10 }

	condition:
		all of them
}

rule LibTCG_PICOLoader_533d91eb {
	meta:
		description = "LibTCG: PICO loading and parsing functions. NOTE: wildcard non-volatile registers for a more robust rule"
		author = "Raphael Mudge"
		date = "2026-02-02"
		reference = "https://tradecraftgarden.org/libtcg.html"
		arch_context = "x64"
		scan_context = "file, memory"
		os = "windows"
		license = "BSD"
		generator = "Crystal Palace"
	strings:
		// ----------------------------------------
		// Function: PicoEntryPoint
		// ----------------------------------------
		/*
		 * 8B 41 0C                      mov eax, dword ptr [rcx+0xC]
		 * 48 63 C8                      movsxd rcx, eax
		 * 48 01 CA                      add rdx, rcx
		 * 85 C0                         test eax, eax
		 * B8 00 00 00 00                mov eax, 0
		 * (Score: 24)
		 */
		$r0_PicoEntryPoint = { 8B 41 0C 48 63 C8 48 01 CA 85 C0 B8 00 00 00 00 }

		// ----------------------------------------
		// Function: PicoLoad
		// ----------------------------------------
		/*
		 * 49 89 CD                      mov r13, rcx
		 * 49 89 D6                      mov r14, rdx
		 * 4C 89 C5                      mov rbp, r8
		 * 4D 89 CC                      mov r12, r9
		 * 48 8D 5A 10                   lea rbx, [rdx+0x10]
		 * (Score: 46)
		 */
		$r1_PicoLoad = { 49 89 CD 49 89 D6 4C 89 C5 4D 89 CC 48 8D 5A 10 }

		/*
		 * 48 63 43 08                   movsxd rax, dword ptr [rbx+8]
		 * 48 01 C7                      add rdi, rax
		 * 49 63 76 08                   movsxd rsi, dword ptr [r14+8]
		 * 48 63 43 04                   movsxd rax, dword ptr [rbx+4]
		 * (Score: 66)
		 */
		$r2_PicoLoad = { 48 63 43 08 48 01 C7 49 63 76 08 48 63 43 04 }

		/*
		 * 48 8D 4B 04                   lea rcx, [rbx+4]
		 * 41 FF 55 00                   call qword ptr [r13]
		 * 48 89 44 24 28                mov qword ptr [rsp+0x28], rax
		 * (Score: 37)
		 */
		$r3_PicoLoad = { 48 8D 4B 04 41 FF 55 00 48 89 44 24 28 }

		/*
		 * 48 8D 53 04                   lea rdx, [rbx+4]
		 * 48 8B 4C 24 28                mov rcx, qword ptr [rsp+0x28]
		 * 41 FF 55 08                   call qword ptr [r13+8]
		 * (Score: 50)
		 */
		$r4_PicoLoad = { 48 8D 53 04 48 8B 4C 24 28 41 FF 55 08 }

	condition:
		all of them
}

rule LibTCG_ResolveEAT_1bfccb0e {
	meta:
		description = "LibTCG: Export Address Table Win32 API resolution derived from ReflectiveDllInjection"
		author = "Raphael Mudge"
		date = "2026-02-02"
		reference = "https://tradecraftgarden.org/libtcg.html"
		arch_context = "x64"
		scan_context = "file, memory"
		os = "windows"
		license = "BSD"
		generator = "Crystal Palace"
	strings:
		// ----------------------------------------
		// Function: findModuleByHash
		// ----------------------------------------
		/*
		 * 65 48 8B 04 25 60 00 00 00    mov rax, qword ptr gs:[0x60]
		 * 48 8B 40 18                   mov rax, qword ptr [rax+0x18]
		 * 4C 8B 50 20                   mov r10, qword ptr [rax+0x20]
		 * (Score: 19)
		 */
		$r0_findModuleByHash = { 65 48 8B 04 25 60 00 00 00 48 8B 40 18 4C 8B 50 20 }

		/*
		 * C1 C8 0D                      ror eax, 0xD
		 * 89 C2                         mov edx, eax
		 * 41 0F B6 00                   movzx eax, byte ptr [r8]
		 * 3C 60                         cmp al, 0x60
		 * (Score: 18)
		 */
		$r1_findModuleByHash = { C1 C8 0D 89 C2 41 0F B6 00 3C 60 }

		/*
		 * 4D 8B 42 50                   mov r8, qword ptr [r10+0x50]
		 * 41 0F B7 42 48                movzx eax, word ptr [r10+0x48]
		 * 83 E8 01                      sub eax, 1
		 * 0F B7 C0                      movzx eax, ax
		 * (Score: 52)
		 */
		$r2_findModuleByHash = { 4D 8B 42 50 41 0F B7 42 48 83 E8 01 0F B7 C0 }

		// ----------------------------------------
		// Function: findFunctionByHash
		// ----------------------------------------
		/*
		 * C1 CA 0D                      ror edx, 0xD
		 * 0F BE 08                      movsx ecx, byte ptr [rax]
		 * 01 CA                         add edx, ecx
		 * 48 83 C0 01                   add rax, 1
		 * 80 38 00                      cmp byte ptr [rax], 0
		 * (Score: 17)
		 */
		$r3_findFunctionByHash = { C1 CA 0D 0F BE 08 01 CA 48 83 C0 01 80 38 00 }

		// ----------------------------------------
		// Function: GetDataDirectory
		// ----------------------------------------
		/*
		 * 48 8B 41 10                   mov rax, qword ptr [rcx+0x10]
		 * 89 D2                         mov edx, edx
		 * 48 8D 44 D0 70                lea rax, [rax+rdx*8+0x70]
		 * (Score: 21)
		 */
		$r4_GetDataDirectory = { 48 8B 41 10 89 D2 48 8D 44 D0 70 }

		// ----------------------------------------
		// Function: ParseDLL
		// ----------------------------------------
		/*
		 * 48 01 C1                      add rcx, rax
		 * 48 89 4A 08                   mov qword ptr [rdx+8], rcx
		 * 48 83 C1 18                   add rcx, 0x18
		 * 48 89 4A 10                   mov qword ptr [rdx+0x10], rcx
		 * (Score: 74)
		 */
		$r5_ParseDLL = { 48 01 C1 48 89 4A 08 48 83 C1 18 48 89 4A 10 }

	condition:
		all of them
}

