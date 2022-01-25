package strfind

// ContainsCtrlChar returns true if s contains any control character,
// otherwise returns false.
func ContainsCtrlChar(s string) bool {
	const sc = 0x20

	if len(s) < 4 {
		for i := 0; i < len(s); i++ {
			if s[i] < sc {
				return true
			}
		}
		return false
	}
	if len(s) < 8 {
		if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc {
			return true
		}
		s = s[4:]
		for i := 0; i < len(s); i++ {
			if s[i] < sc {
				return true
			}
		}
		return false
	}
	if len(s) < 16 {
		if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
			s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc {
			return true
		}
		for i := 8; i < len(s); i++ {
			if s[i] < sc {
				return true
			}
		}
		return false
	}
	if len(s) < 32 {
		if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
			s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc ||
			s[8] < sc || s[9] < sc || s[10] < sc || s[11] < sc ||
			s[12] < sc || s[13] < sc || s[14] < sc || s[15] < sc {
			return true
		}
		s = s[16:]
		for len(s) >= 8 {
			if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
				s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc {
				return true
			}
			s = s[8:]
		}
		for i := 0; i < len(s); i++ {
			if s[i] < sc {
				return true
			}
		}
		return false
	}
	if len(s) < 64 {
		if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
			s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc ||
			s[8] < sc || s[9] < sc || s[10] < sc || s[11] < sc ||
			s[12] < sc || s[13] < sc || s[14] < sc || s[15] < sc ||
			s[16] < sc || s[17] < sc || s[18] < sc || s[19] < sc ||
			s[20] < sc || s[21] < sc || s[22] < sc || s[23] < sc ||
			s[24] < sc || s[25] < sc || s[26] < sc || s[27] < sc ||
			s[28] < sc || s[29] < sc || s[30] < sc || s[31] < sc {
			return true
		}
		s = s[32:]
		for len(s) >= 16 {
			if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
				s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc ||
				s[8] < sc || s[9] < sc || s[10] < sc || s[11] < sc ||
				s[12] < sc || s[13] < sc || s[14] < sc || s[15] < sc {
				return true
			}
			s = s[16:]
		}
		for len(s) >= 8 {
			if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
				s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc {
				return true
			}
			s = s[8:]
		}
		for i := 0; i < len(s); i++ {
			if s[i] < sc {
				return true
			}
		}
		return false
	}
	if len(s) < 256 {
		if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
			s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc ||
			s[8] < sc || s[9] < sc || s[10] < sc || s[11] < sc ||
			s[12] < sc || s[13] < sc || s[14] < sc || s[15] < sc ||
			s[16] < sc || s[17] < sc || s[18] < sc || s[19] < sc ||
			s[20] < sc || s[21] < sc || s[22] < sc || s[23] < sc ||
			s[24] < sc || s[25] < sc || s[26] < sc || s[27] < sc ||
			s[28] < sc || s[29] < sc || s[30] < sc || s[31] < sc ||
			s[32] < sc || s[33] < sc || s[34] < sc || s[35] < sc ||
			s[36] < sc || s[37] < sc || s[38] < sc || s[39] < sc ||
			s[40] < sc || s[41] < sc || s[42] < sc || s[43] < sc ||
			s[44] < sc || s[45] < sc || s[46] < sc || s[47] < sc ||
			s[48] < sc || s[49] < sc || s[50] < sc || s[51] < sc ||
			s[52] < sc || s[53] < sc || s[54] < sc || s[55] < sc ||
			s[56] < sc || s[57] < sc || s[58] < sc || s[59] < sc ||
			s[60] < sc || s[61] < sc || s[62] < sc || s[63] < sc {
			return true
		}
		s = s[64:]
		for len(s) >= 32 {
			if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
				s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc ||
				s[8] < sc || s[9] < sc || s[10] < sc || s[11] < sc ||
				s[12] < sc || s[13] < sc || s[14] < sc || s[15] < sc ||
				s[16] < sc || s[17] < sc || s[18] < sc || s[19] < sc ||
				s[20] < sc || s[21] < sc || s[22] < sc || s[23] < sc ||
				s[24] < sc || s[25] < sc || s[26] < sc || s[27] < sc ||
				s[28] < sc || s[29] < sc || s[30] < sc || s[31] < sc {
				return true
			}
			s = s[32:]
		}
		for len(s) >= 16 {
			if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
				s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc ||
				s[8] < sc || s[9] < sc || s[10] < sc || s[11] < sc ||
				s[12] < sc || s[13] < sc || s[14] < sc || s[15] < sc {
				return true
			}
			s = s[16:]
		}
		for len(s) >= 8 {
			if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
				s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc {
				return true
			}
			s = s[8:]
		}
		for i := 0; i < len(s); i++ {
			if s[i] < sc {
				return true
			}
		}
		return false
	}
	if len(s) < 1024 {
		if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
			s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc ||
			s[8] < sc || s[9] < sc || s[10] < sc || s[11] < sc ||
			s[12] < sc || s[13] < sc || s[14] < sc || s[15] < sc ||
			s[16] < sc || s[17] < sc || s[18] < sc || s[19] < sc ||
			s[20] < sc || s[21] < sc || s[22] < sc || s[23] < sc ||
			s[24] < sc || s[25] < sc || s[26] < sc || s[27] < sc ||
			s[28] < sc || s[29] < sc || s[30] < sc || s[31] < sc ||
			s[32] < sc || s[33] < sc || s[34] < sc || s[35] < sc ||
			s[36] < sc || s[37] < sc || s[38] < sc || s[39] < sc ||
			s[40] < sc || s[41] < sc || s[42] < sc || s[43] < sc ||
			s[44] < sc || s[45] < sc || s[46] < sc || s[47] < sc ||
			s[48] < sc || s[49] < sc || s[50] < sc || s[51] < sc ||
			s[52] < sc || s[53] < sc || s[54] < sc || s[55] < sc ||
			s[56] < sc || s[57] < sc || s[58] < sc || s[59] < sc ||
			s[60] < sc || s[61] < sc || s[62] < sc || s[63] < sc ||
			s[64] < sc || s[65] < sc || s[66] < sc || s[67] < sc ||
			s[68] < sc || s[69] < sc || s[70] < sc || s[71] < sc ||
			s[72] < sc || s[73] < sc || s[74] < sc || s[75] < sc ||
			s[76] < sc || s[77] < sc || s[78] < sc || s[79] < sc ||
			s[80] < sc || s[81] < sc || s[82] < sc || s[83] < sc ||
			s[84] < sc || s[85] < sc || s[86] < sc || s[87] < sc ||
			s[88] < sc || s[89] < sc || s[90] < sc || s[91] < sc ||
			s[92] < sc || s[93] < sc || s[94] < sc || s[95] < sc ||
			s[96] < sc || s[97] < sc || s[98] < sc || s[99] < sc ||
			s[100] < sc || s[101] < sc || s[102] < sc || s[103] < sc ||
			s[104] < sc || s[105] < sc || s[106] < sc || s[107] < sc ||
			s[108] < sc || s[109] < sc || s[110] < sc || s[111] < sc ||
			s[112] < sc || s[113] < sc || s[114] < sc || s[115] < sc ||
			s[116] < sc || s[117] < sc || s[118] < sc || s[119] < sc ||
			s[120] < sc || s[121] < sc || s[122] < sc || s[123] < sc ||
			s[124] < sc || s[125] < sc || s[126] < sc || s[127] < sc ||
			s[128] < sc || s[129] < sc || s[130] < sc || s[131] < sc ||
			s[132] < sc || s[133] < sc || s[134] < sc || s[135] < sc ||
			s[136] < sc || s[137] < sc || s[138] < sc || s[139] < sc ||
			s[140] < sc || s[141] < sc || s[142] < sc || s[143] < sc ||
			s[144] < sc || s[145] < sc || s[146] < sc || s[147] < sc ||
			s[148] < sc || s[149] < sc || s[150] < sc || s[151] < sc ||
			s[152] < sc || s[153] < sc || s[154] < sc || s[155] < sc ||
			s[156] < sc || s[157] < sc || s[158] < sc || s[159] < sc ||
			s[160] < sc || s[161] < sc || s[162] < sc || s[163] < sc ||
			s[164] < sc || s[165] < sc || s[166] < sc || s[167] < sc ||
			s[168] < sc || s[169] < sc || s[170] < sc || s[171] < sc ||
			s[172] < sc || s[173] < sc || s[174] < sc || s[175] < sc ||
			s[176] < sc || s[177] < sc || s[178] < sc || s[179] < sc ||
			s[180] < sc || s[181] < sc || s[182] < sc || s[183] < sc ||
			s[184] < sc || s[185] < sc || s[186] < sc || s[187] < sc ||
			s[188] < sc || s[189] < sc || s[190] < sc || s[191] < sc ||
			s[192] < sc || s[193] < sc || s[194] < sc || s[195] < sc ||
			s[196] < sc || s[197] < sc || s[198] < sc || s[199] < sc ||
			s[200] < sc || s[201] < sc || s[202] < sc || s[203] < sc ||
			s[204] < sc || s[205] < sc || s[206] < sc || s[207] < sc ||
			s[208] < sc || s[209] < sc || s[210] < sc || s[211] < sc ||
			s[212] < sc || s[213] < sc || s[214] < sc || s[215] < sc ||
			s[216] < sc || s[217] < sc || s[218] < sc || s[219] < sc ||
			s[220] < sc || s[221] < sc || s[222] < sc || s[223] < sc ||
			s[224] < sc || s[225] < sc || s[226] < sc || s[227] < sc ||
			s[228] < sc || s[229] < sc || s[230] < sc || s[231] < sc ||
			s[232] < sc || s[233] < sc || s[234] < sc || s[235] < sc ||
			s[236] < sc || s[237] < sc || s[238] < sc || s[239] < sc ||
			s[240] < sc || s[241] < sc || s[242] < sc || s[243] < sc ||
			s[244] < sc || s[245] < sc || s[246] < sc || s[247] < sc ||
			s[248] < sc || s[249] < sc || s[250] < sc || s[251] < sc ||
			s[252] < sc || s[253] < sc || s[254] < sc || s[255] < sc {
			return true
		}
		s = s[256:]
		for len(s) >= 64 {
			if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
				s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc ||
				s[8] < sc || s[9] < sc || s[10] < sc || s[11] < sc ||
				s[12] < sc || s[13] < sc || s[14] < sc || s[15] < sc ||
				s[16] < sc || s[17] < sc || s[18] < sc || s[19] < sc ||
				s[20] < sc || s[21] < sc || s[22] < sc || s[23] < sc ||
				s[24] < sc || s[25] < sc || s[26] < sc || s[27] < sc ||
				s[28] < sc || s[29] < sc || s[30] < sc || s[31] < sc ||
				s[32] < sc || s[33] < sc || s[34] < sc || s[35] < sc ||
				s[36] < sc || s[37] < sc || s[38] < sc || s[39] < sc ||
				s[40] < sc || s[41] < sc || s[42] < sc || s[43] < sc ||
				s[44] < sc || s[45] < sc || s[46] < sc || s[47] < sc ||
				s[48] < sc || s[49] < sc || s[50] < sc || s[51] < sc ||
				s[52] < sc || s[53] < sc || s[54] < sc || s[55] < sc ||
				s[56] < sc || s[57] < sc || s[58] < sc || s[59] < sc ||
				s[60] < sc || s[61] < sc || s[62] < sc || s[63] < sc {
				return true
			}
			s = s[64:]
		}
		for len(s) >= 32 {
			if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
				s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc ||
				s[8] < sc || s[9] < sc || s[10] < sc || s[11] < sc ||
				s[12] < sc || s[13] < sc || s[14] < sc || s[15] < sc ||
				s[16] < sc || s[17] < sc || s[18] < sc || s[19] < sc ||
				s[20] < sc || s[21] < sc || s[22] < sc || s[23] < sc ||
				s[24] < sc || s[25] < sc || s[26] < sc || s[27] < sc ||
				s[28] < sc || s[29] < sc || s[30] < sc || s[31] < sc {
				return true
			}
			s = s[32:]
		}
		for len(s) >= 16 {
			if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
				s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc ||
				s[8] < sc || s[9] < sc || s[10] < sc || s[11] < sc ||
				s[12] < sc || s[13] < sc || s[14] < sc || s[15] < sc {
				return true
			}
			s = s[16:]
		}
		for len(s) >= 8 {
			if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
				s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc {
				return true
			}
			s = s[8:]
		}
		for i := 0; i < len(s); i++ {
			if s[i] < sc {
				return true
			}
		}
		return false
	}

	// Worst case
	for len(s) >= 1024 {
		if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
			s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc ||
			s[8] < sc || s[9] < sc || s[10] < sc || s[11] < sc ||
			s[12] < sc || s[13] < sc || s[14] < sc || s[15] < sc ||
			s[16] < sc || s[17] < sc || s[18] < sc || s[19] < sc ||
			s[20] < sc || s[21] < sc || s[22] < sc || s[23] < sc ||
			s[24] < sc || s[25] < sc || s[26] < sc || s[27] < sc ||
			s[28] < sc || s[29] < sc || s[30] < sc || s[31] < sc ||
			s[32] < sc || s[33] < sc || s[34] < sc || s[35] < sc ||
			s[36] < sc || s[37] < sc || s[38] < sc || s[39] < sc ||
			s[40] < sc || s[41] < sc || s[42] < sc || s[43] < sc ||
			s[44] < sc || s[45] < sc || s[46] < sc || s[47] < sc ||
			s[48] < sc || s[49] < sc || s[50] < sc || s[51] < sc ||
			s[52] < sc || s[53] < sc || s[54] < sc || s[55] < sc ||
			s[56] < sc || s[57] < sc || s[58] < sc || s[59] < sc ||
			s[60] < sc || s[61] < sc || s[62] < sc || s[63] < sc ||
			s[64] < sc || s[65] < sc || s[66] < sc || s[67] < sc ||
			s[68] < sc || s[69] < sc || s[70] < sc || s[71] < sc ||
			s[72] < sc || s[73] < sc || s[74] < sc || s[75] < sc ||
			s[76] < sc || s[77] < sc || s[78] < sc || s[79] < sc ||
			s[80] < sc || s[81] < sc || s[82] < sc || s[83] < sc ||
			s[84] < sc || s[85] < sc || s[86] < sc || s[87] < sc ||
			s[88] < sc || s[89] < sc || s[90] < sc || s[91] < sc ||
			s[92] < sc || s[93] < sc || s[94] < sc || s[95] < sc ||
			s[96] < sc || s[97] < sc || s[98] < sc || s[99] < sc ||
			s[100] < sc || s[101] < sc || s[102] < sc || s[103] < sc ||
			s[104] < sc || s[105] < sc || s[106] < sc || s[107] < sc ||
			s[108] < sc || s[109] < sc || s[110] < sc || s[111] < sc ||
			s[112] < sc || s[113] < sc || s[114] < sc || s[115] < sc ||
			s[116] < sc || s[117] < sc || s[118] < sc || s[119] < sc ||
			s[120] < sc || s[121] < sc || s[122] < sc || s[123] < sc ||
			s[124] < sc || s[125] < sc || s[126] < sc || s[127] < sc ||
			s[128] < sc || s[129] < sc || s[130] < sc || s[131] < sc ||
			s[132] < sc || s[133] < sc || s[134] < sc || s[135] < sc ||
			s[136] < sc || s[137] < sc || s[138] < sc || s[139] < sc ||
			s[140] < sc || s[141] < sc || s[142] < sc || s[143] < sc ||
			s[144] < sc || s[145] < sc || s[146] < sc || s[147] < sc ||
			s[148] < sc || s[149] < sc || s[150] < sc || s[151] < sc ||
			s[152] < sc || s[153] < sc || s[154] < sc || s[155] < sc ||
			s[156] < sc || s[157] < sc || s[158] < sc || s[159] < sc ||
			s[160] < sc || s[161] < sc || s[162] < sc || s[163] < sc ||
			s[164] < sc || s[165] < sc || s[166] < sc || s[167] < sc ||
			s[168] < sc || s[169] < sc || s[170] < sc || s[171] < sc ||
			s[172] < sc || s[173] < sc || s[174] < sc || s[175] < sc ||
			s[176] < sc || s[177] < sc || s[178] < sc || s[179] < sc ||
			s[180] < sc || s[181] < sc || s[182] < sc || s[183] < sc ||
			s[184] < sc || s[185] < sc || s[186] < sc || s[187] < sc ||
			s[188] < sc || s[189] < sc || s[190] < sc || s[191] < sc ||
			s[192] < sc || s[193] < sc || s[194] < sc || s[195] < sc ||
			s[196] < sc || s[197] < sc || s[198] < sc || s[199] < sc ||
			s[200] < sc || s[201] < sc || s[202] < sc || s[203] < sc ||
			s[204] < sc || s[205] < sc || s[206] < sc || s[207] < sc ||
			s[208] < sc || s[209] < sc || s[210] < sc || s[211] < sc ||
			s[212] < sc || s[213] < sc || s[214] < sc || s[215] < sc ||
			s[216] < sc || s[217] < sc || s[218] < sc || s[219] < sc ||
			s[220] < sc || s[221] < sc || s[222] < sc || s[223] < sc ||
			s[224] < sc || s[225] < sc || s[226] < sc || s[227] < sc ||
			s[228] < sc || s[229] < sc || s[230] < sc || s[231] < sc ||
			s[232] < sc || s[233] < sc || s[234] < sc || s[235] < sc ||
			s[236] < sc || s[237] < sc || s[238] < sc || s[239] < sc ||
			s[240] < sc || s[241] < sc || s[242] < sc || s[243] < sc ||
			s[244] < sc || s[245] < sc || s[246] < sc || s[247] < sc ||
			s[248] < sc || s[249] < sc || s[250] < sc || s[251] < sc ||
			s[252] < sc || s[253] < sc || s[254] < sc || s[255] < sc ||
			s[256] < sc || s[257] < sc || s[258] < sc || s[259] < sc ||
			s[260] < sc || s[261] < sc || s[262] < sc || s[263] < sc ||
			s[264] < sc || s[265] < sc || s[266] < sc || s[267] < sc ||
			s[268] < sc || s[269] < sc || s[270] < sc || s[271] < sc ||
			s[272] < sc || s[273] < sc || s[274] < sc || s[275] < sc ||
			s[276] < sc || s[277] < sc || s[278] < sc || s[279] < sc ||
			s[280] < sc || s[281] < sc || s[282] < sc || s[283] < sc ||
			s[284] < sc || s[285] < sc || s[286] < sc || s[287] < sc ||
			s[288] < sc || s[289] < sc || s[290] < sc || s[291] < sc ||
			s[292] < sc || s[293] < sc || s[294] < sc || s[295] < sc ||
			s[296] < sc || s[297] < sc || s[298] < sc || s[299] < sc ||
			s[300] < sc || s[301] < sc || s[302] < sc || s[303] < sc ||
			s[304] < sc || s[305] < sc || s[306] < sc || s[307] < sc ||
			s[308] < sc || s[309] < sc || s[310] < sc || s[311] < sc ||
			s[312] < sc || s[313] < sc || s[314] < sc || s[315] < sc ||
			s[316] < sc || s[317] < sc || s[318] < sc || s[319] < sc ||
			s[320] < sc || s[321] < sc || s[322] < sc || s[323] < sc ||
			s[324] < sc || s[325] < sc || s[326] < sc || s[327] < sc ||
			s[328] < sc || s[329] < sc || s[330] < sc || s[331] < sc ||
			s[332] < sc || s[333] < sc || s[334] < sc || s[335] < sc ||
			s[336] < sc || s[337] < sc || s[338] < sc || s[339] < sc ||
			s[340] < sc || s[341] < sc || s[342] < sc || s[343] < sc ||
			s[344] < sc || s[345] < sc || s[346] < sc || s[347] < sc ||
			s[348] < sc || s[349] < sc || s[350] < sc || s[351] < sc ||
			s[352] < sc || s[353] < sc || s[354] < sc || s[355] < sc ||
			s[356] < sc || s[357] < sc || s[358] < sc || s[359] < sc ||
			s[360] < sc || s[361] < sc || s[362] < sc || s[363] < sc ||
			s[364] < sc || s[365] < sc || s[366] < sc || s[367] < sc ||
			s[368] < sc || s[369] < sc || s[370] < sc || s[371] < sc ||
			s[372] < sc || s[373] < sc || s[374] < sc || s[375] < sc ||
			s[376] < sc || s[377] < sc || s[378] < sc || s[379] < sc ||
			s[380] < sc || s[381] < sc || s[382] < sc || s[383] < sc ||
			s[384] < sc || s[385] < sc || s[386] < sc || s[387] < sc ||
			s[388] < sc || s[389] < sc || s[390] < sc || s[391] < sc ||
			s[392] < sc || s[393] < sc || s[394] < sc || s[395] < sc ||
			s[396] < sc || s[397] < sc || s[398] < sc || s[399] < sc ||
			s[400] < sc || s[401] < sc || s[402] < sc || s[403] < sc ||
			s[404] < sc || s[405] < sc || s[406] < sc || s[407] < sc ||
			s[408] < sc || s[409] < sc || s[410] < sc || s[411] < sc ||
			s[412] < sc || s[413] < sc || s[414] < sc || s[415] < sc ||
			s[416] < sc || s[417] < sc || s[418] < sc || s[419] < sc ||
			s[420] < sc || s[421] < sc || s[422] < sc || s[423] < sc ||
			s[424] < sc || s[425] < sc || s[426] < sc || s[427] < sc ||
			s[428] < sc || s[429] < sc || s[430] < sc || s[431] < sc ||
			s[432] < sc || s[433] < sc || s[434] < sc || s[435] < sc ||
			s[436] < sc || s[437] < sc || s[438] < sc || s[439] < sc ||
			s[440] < sc || s[441] < sc || s[442] < sc || s[443] < sc ||
			s[444] < sc || s[445] < sc || s[446] < sc || s[447] < sc ||
			s[448] < sc || s[449] < sc || s[450] < sc || s[451] < sc ||
			s[452] < sc || s[453] < sc || s[454] < sc || s[455] < sc ||
			s[456] < sc || s[457] < sc || s[458] < sc || s[459] < sc ||
			s[460] < sc || s[461] < sc || s[462] < sc || s[463] < sc ||
			s[464] < sc || s[465] < sc || s[466] < sc || s[467] < sc ||
			s[468] < sc || s[469] < sc || s[470] < sc || s[471] < sc ||
			s[472] < sc || s[473] < sc || s[474] < sc || s[475] < sc ||
			s[476] < sc || s[477] < sc || s[478] < sc || s[479] < sc ||
			s[480] < sc || s[481] < sc || s[482] < sc || s[483] < sc ||
			s[484] < sc || s[485] < sc || s[486] < sc || s[487] < sc ||
			s[488] < sc || s[489] < sc || s[490] < sc || s[491] < sc ||
			s[492] < sc || s[493] < sc || s[494] < sc || s[495] < sc ||
			s[496] < sc || s[497] < sc || s[498] < sc || s[499] < sc ||
			s[500] < sc || s[501] < sc || s[502] < sc || s[503] < sc ||
			s[504] < sc || s[505] < sc || s[506] < sc || s[507] < sc ||
			s[508] < sc || s[509] < sc || s[510] < sc || s[511] < sc ||
			s[512] < sc || s[513] < sc || s[514] < sc || s[515] < sc ||
			s[516] < sc || s[517] < sc || s[518] < sc || s[519] < sc ||
			s[520] < sc || s[521] < sc || s[522] < sc || s[523] < sc ||
			s[524] < sc || s[525] < sc || s[526] < sc || s[527] < sc ||
			s[528] < sc || s[529] < sc || s[530] < sc || s[531] < sc ||
			s[532] < sc || s[533] < sc || s[534] < sc || s[535] < sc ||
			s[536] < sc || s[537] < sc || s[538] < sc || s[539] < sc ||
			s[540] < sc || s[541] < sc || s[542] < sc || s[543] < sc ||
			s[544] < sc || s[545] < sc || s[546] < sc || s[547] < sc ||
			s[548] < sc || s[549] < sc || s[550] < sc || s[551] < sc ||
			s[552] < sc || s[553] < sc || s[554] < sc || s[555] < sc ||
			s[556] < sc || s[557] < sc || s[558] < sc || s[559] < sc ||
			s[560] < sc || s[561] < sc || s[562] < sc || s[563] < sc ||
			s[564] < sc || s[565] < sc || s[566] < sc || s[567] < sc ||
			s[568] < sc || s[569] < sc || s[570] < sc || s[571] < sc ||
			s[572] < sc || s[573] < sc || s[574] < sc || s[575] < sc ||
			s[576] < sc || s[577] < sc || s[578] < sc || s[579] < sc ||
			s[580] < sc || s[581] < sc || s[582] < sc || s[583] < sc ||
			s[584] < sc || s[585] < sc || s[586] < sc || s[587] < sc ||
			s[588] < sc || s[589] < sc || s[590] < sc || s[591] < sc ||
			s[592] < sc || s[593] < sc || s[594] < sc || s[595] < sc ||
			s[596] < sc || s[597] < sc || s[598] < sc || s[599] < sc ||
			s[600] < sc || s[601] < sc || s[602] < sc || s[603] < sc ||
			s[604] < sc || s[605] < sc || s[606] < sc || s[607] < sc ||
			s[608] < sc || s[609] < sc || s[610] < sc || s[611] < sc ||
			s[612] < sc || s[613] < sc || s[614] < sc || s[615] < sc ||
			s[616] < sc || s[617] < sc || s[618] < sc || s[619] < sc ||
			s[620] < sc || s[621] < sc || s[622] < sc || s[623] < sc ||
			s[624] < sc || s[625] < sc || s[626] < sc || s[627] < sc ||
			s[628] < sc || s[629] < sc || s[630] < sc || s[631] < sc ||
			s[632] < sc || s[633] < sc || s[634] < sc || s[635] < sc ||
			s[636] < sc || s[637] < sc || s[638] < sc || s[639] < sc ||
			s[640] < sc || s[641] < sc || s[642] < sc || s[643] < sc ||
			s[644] < sc || s[645] < sc || s[646] < sc || s[647] < sc ||
			s[648] < sc || s[649] < sc || s[650] < sc || s[651] < sc ||
			s[652] < sc || s[653] < sc || s[654] < sc || s[655] < sc ||
			s[656] < sc || s[657] < sc || s[658] < sc || s[659] < sc ||
			s[660] < sc || s[661] < sc || s[662] < sc || s[663] < sc ||
			s[664] < sc || s[665] < sc || s[666] < sc || s[667] < sc ||
			s[668] < sc || s[669] < sc || s[670] < sc || s[671] < sc ||
			s[672] < sc || s[673] < sc || s[674] < sc || s[675] < sc ||
			s[676] < sc || s[677] < sc || s[678] < sc || s[679] < sc ||
			s[680] < sc || s[681] < sc || s[682] < sc || s[683] < sc ||
			s[684] < sc || s[685] < sc || s[686] < sc || s[687] < sc ||
			s[688] < sc || s[689] < sc || s[690] < sc || s[691] < sc ||
			s[692] < sc || s[693] < sc || s[694] < sc || s[695] < sc ||
			s[696] < sc || s[697] < sc || s[698] < sc || s[699] < sc ||
			s[700] < sc || s[701] < sc || s[702] < sc || s[703] < sc ||
			s[704] < sc || s[705] < sc || s[706] < sc || s[707] < sc ||
			s[708] < sc || s[709] < sc || s[710] < sc || s[711] < sc ||
			s[712] < sc || s[713] < sc || s[714] < sc || s[715] < sc ||
			s[716] < sc || s[717] < sc || s[718] < sc || s[719] < sc ||
			s[720] < sc || s[721] < sc || s[722] < sc || s[723] < sc ||
			s[724] < sc || s[725] < sc || s[726] < sc || s[727] < sc ||
			s[728] < sc || s[729] < sc || s[730] < sc || s[731] < sc ||
			s[732] < sc || s[733] < sc || s[734] < sc || s[735] < sc ||
			s[736] < sc || s[737] < sc || s[738] < sc || s[739] < sc ||
			s[740] < sc || s[741] < sc || s[742] < sc || s[743] < sc ||
			s[744] < sc || s[745] < sc || s[746] < sc || s[747] < sc ||
			s[748] < sc || s[749] < sc || s[750] < sc || s[751] < sc ||
			s[752] < sc || s[753] < sc || s[754] < sc || s[755] < sc ||
			s[756] < sc || s[757] < sc || s[758] < sc || s[759] < sc ||
			s[760] < sc || s[761] < sc || s[762] < sc || s[763] < sc ||
			s[764] < sc || s[765] < sc || s[766] < sc || s[767] < sc ||
			s[768] < sc || s[769] < sc || s[770] < sc || s[771] < sc ||
			s[772] < sc || s[773] < sc || s[774] < sc || s[775] < sc ||
			s[776] < sc || s[777] < sc || s[778] < sc || s[779] < sc ||
			s[780] < sc || s[781] < sc || s[782] < sc || s[783] < sc ||
			s[784] < sc || s[785] < sc || s[786] < sc || s[787] < sc ||
			s[788] < sc || s[789] < sc || s[790] < sc || s[791] < sc ||
			s[792] < sc || s[793] < sc || s[794] < sc || s[795] < sc ||
			s[796] < sc || s[797] < sc || s[798] < sc || s[799] < sc ||
			s[800] < sc || s[801] < sc || s[802] < sc || s[803] < sc ||
			s[804] < sc || s[805] < sc || s[806] < sc || s[807] < sc ||
			s[808] < sc || s[809] < sc || s[810] < sc || s[811] < sc ||
			s[812] < sc || s[813] < sc || s[814] < sc || s[815] < sc ||
			s[816] < sc || s[817] < sc || s[818] < sc || s[819] < sc ||
			s[820] < sc || s[821] < sc || s[822] < sc || s[823] < sc ||
			s[824] < sc || s[825] < sc || s[826] < sc || s[827] < sc ||
			s[828] < sc || s[829] < sc || s[830] < sc || s[831] < sc ||
			s[832] < sc || s[833] < sc || s[834] < sc || s[835] < sc ||
			s[836] < sc || s[837] < sc || s[838] < sc || s[839] < sc ||
			s[840] < sc || s[841] < sc || s[842] < sc || s[843] < sc ||
			s[844] < sc || s[845] < sc || s[846] < sc || s[847] < sc ||
			s[848] < sc || s[849] < sc || s[850] < sc || s[851] < sc ||
			s[852] < sc || s[853] < sc || s[854] < sc || s[855] < sc ||
			s[856] < sc || s[857] < sc || s[858] < sc || s[859] < sc ||
			s[860] < sc || s[861] < sc || s[862] < sc || s[863] < sc ||
			s[864] < sc || s[865] < sc || s[866] < sc || s[867] < sc ||
			s[868] < sc || s[869] < sc || s[870] < sc || s[871] < sc ||
			s[872] < sc || s[873] < sc || s[874] < sc || s[875] < sc ||
			s[876] < sc || s[877] < sc || s[878] < sc || s[879] < sc ||
			s[880] < sc || s[881] < sc || s[882] < sc || s[883] < sc ||
			s[884] < sc || s[885] < sc || s[886] < sc || s[887] < sc ||
			s[888] < sc || s[889] < sc || s[890] < sc || s[891] < sc ||
			s[892] < sc || s[893] < sc || s[894] < sc || s[895] < sc ||
			s[896] < sc || s[897] < sc || s[898] < sc || s[899] < sc ||
			s[900] < sc || s[901] < sc || s[902] < sc || s[903] < sc ||
			s[904] < sc || s[905] < sc || s[906] < sc || s[907] < sc ||
			s[908] < sc || s[909] < sc || s[910] < sc || s[911] < sc ||
			s[912] < sc || s[913] < sc || s[914] < sc || s[915] < sc ||
			s[916] < sc || s[917] < sc || s[918] < sc || s[919] < sc ||
			s[920] < sc || s[921] < sc || s[922] < sc || s[923] < sc ||
			s[924] < sc || s[925] < sc || s[926] < sc || s[927] < sc ||
			s[928] < sc || s[929] < sc || s[930] < sc || s[931] < sc ||
			s[932] < sc || s[933] < sc || s[934] < sc || s[935] < sc ||
			s[936] < sc || s[937] < sc || s[938] < sc || s[939] < sc ||
			s[940] < sc || s[941] < sc || s[942] < sc || s[943] < sc ||
			s[944] < sc || s[945] < sc || s[946] < sc || s[947] < sc ||
			s[948] < sc || s[949] < sc || s[950] < sc || s[951] < sc ||
			s[952] < sc || s[953] < sc || s[954] < sc || s[955] < sc ||
			s[956] < sc || s[957] < sc || s[958] < sc || s[959] < sc ||
			s[960] < sc || s[961] < sc || s[962] < sc || s[963] < sc ||
			s[964] < sc || s[965] < sc || s[966] < sc || s[967] < sc ||
			s[968] < sc || s[969] < sc || s[970] < sc || s[971] < sc ||
			s[972] < sc || s[973] < sc || s[974] < sc || s[975] < sc ||
			s[976] < sc || s[977] < sc || s[978] < sc || s[979] < sc ||
			s[980] < sc || s[981] < sc || s[982] < sc || s[983] < sc ||
			s[984] < sc || s[985] < sc || s[986] < sc || s[987] < sc ||
			s[988] < sc || s[989] < sc || s[990] < sc || s[991] < sc ||
			s[992] < sc || s[993] < sc || s[994] < sc || s[995] < sc ||
			s[996] < sc || s[997] < sc || s[998] < sc || s[999] < sc ||
			s[1000] < sc || s[1001] < sc || s[1002] < sc || s[1003] < sc ||
			s[1004] < sc || s[1005] < sc || s[1006] < sc || s[1007] < sc ||
			s[1008] < sc || s[1009] < sc || s[1010] < sc || s[1011] < sc ||
			s[1012] < sc || s[1013] < sc || s[1014] < sc || s[1015] < sc ||
			s[1016] < sc || s[1017] < sc || s[1018] < sc || s[1019] < sc ||
			s[1020] < sc || s[1021] < sc || s[1022] < sc || s[1023] < sc {
			return true
		}
		s = s[1024:]
	}
	for len(s) >= 256 {
		if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
			s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc ||
			s[8] < sc || s[9] < sc || s[10] < sc || s[11] < sc ||
			s[12] < sc || s[13] < sc || s[14] < sc || s[15] < sc ||
			s[16] < sc || s[17] < sc || s[18] < sc || s[19] < sc ||
			s[20] < sc || s[21] < sc || s[22] < sc || s[23] < sc ||
			s[24] < sc || s[25] < sc || s[26] < sc || s[27] < sc ||
			s[28] < sc || s[29] < sc || s[30] < sc || s[31] < sc ||
			s[32] < sc || s[33] < sc || s[34] < sc || s[35] < sc ||
			s[36] < sc || s[37] < sc || s[38] < sc || s[39] < sc ||
			s[40] < sc || s[41] < sc || s[42] < sc || s[43] < sc ||
			s[44] < sc || s[45] < sc || s[46] < sc || s[47] < sc ||
			s[48] < sc || s[49] < sc || s[50] < sc || s[51] < sc ||
			s[52] < sc || s[53] < sc || s[54] < sc || s[55] < sc ||
			s[56] < sc || s[57] < sc || s[58] < sc || s[59] < sc ||
			s[60] < sc || s[61] < sc || s[62] < sc || s[63] < sc ||
			s[64] < sc || s[65] < sc || s[66] < sc || s[67] < sc ||
			s[68] < sc || s[69] < sc || s[70] < sc || s[71] < sc ||
			s[72] < sc || s[73] < sc || s[74] < sc || s[75] < sc ||
			s[76] < sc || s[77] < sc || s[78] < sc || s[79] < sc ||
			s[80] < sc || s[81] < sc || s[82] < sc || s[83] < sc ||
			s[84] < sc || s[85] < sc || s[86] < sc || s[87] < sc ||
			s[88] < sc || s[89] < sc || s[90] < sc || s[91] < sc ||
			s[92] < sc || s[93] < sc || s[94] < sc || s[95] < sc ||
			s[96] < sc || s[97] < sc || s[98] < sc || s[99] < sc ||
			s[100] < sc || s[101] < sc || s[102] < sc || s[103] < sc ||
			s[104] < sc || s[105] < sc || s[106] < sc || s[107] < sc ||
			s[108] < sc || s[109] < sc || s[110] < sc || s[111] < sc ||
			s[112] < sc || s[113] < sc || s[114] < sc || s[115] < sc ||
			s[116] < sc || s[117] < sc || s[118] < sc || s[119] < sc ||
			s[120] < sc || s[121] < sc || s[122] < sc || s[123] < sc ||
			s[124] < sc || s[125] < sc || s[126] < sc || s[127] < sc ||
			s[128] < sc || s[129] < sc || s[130] < sc || s[131] < sc ||
			s[132] < sc || s[133] < sc || s[134] < sc || s[135] < sc ||
			s[136] < sc || s[137] < sc || s[138] < sc || s[139] < sc ||
			s[140] < sc || s[141] < sc || s[142] < sc || s[143] < sc ||
			s[144] < sc || s[145] < sc || s[146] < sc || s[147] < sc ||
			s[148] < sc || s[149] < sc || s[150] < sc || s[151] < sc ||
			s[152] < sc || s[153] < sc || s[154] < sc || s[155] < sc ||
			s[156] < sc || s[157] < sc || s[158] < sc || s[159] < sc ||
			s[160] < sc || s[161] < sc || s[162] < sc || s[163] < sc ||
			s[164] < sc || s[165] < sc || s[166] < sc || s[167] < sc ||
			s[168] < sc || s[169] < sc || s[170] < sc || s[171] < sc ||
			s[172] < sc || s[173] < sc || s[174] < sc || s[175] < sc ||
			s[176] < sc || s[177] < sc || s[178] < sc || s[179] < sc ||
			s[180] < sc || s[181] < sc || s[182] < sc || s[183] < sc ||
			s[184] < sc || s[185] < sc || s[186] < sc || s[187] < sc ||
			s[188] < sc || s[189] < sc || s[190] < sc || s[191] < sc ||
			s[192] < sc || s[193] < sc || s[194] < sc || s[195] < sc ||
			s[196] < sc || s[197] < sc || s[198] < sc || s[199] < sc ||
			s[200] < sc || s[201] < sc || s[202] < sc || s[203] < sc ||
			s[204] < sc || s[205] < sc || s[206] < sc || s[207] < sc ||
			s[208] < sc || s[209] < sc || s[210] < sc || s[211] < sc ||
			s[212] < sc || s[213] < sc || s[214] < sc || s[215] < sc ||
			s[216] < sc || s[217] < sc || s[218] < sc || s[219] < sc ||
			s[220] < sc || s[221] < sc || s[222] < sc || s[223] < sc ||
			s[224] < sc || s[225] < sc || s[226] < sc || s[227] < sc ||
			s[228] < sc || s[229] < sc || s[230] < sc || s[231] < sc ||
			s[232] < sc || s[233] < sc || s[234] < sc || s[235] < sc ||
			s[236] < sc || s[237] < sc || s[238] < sc || s[239] < sc ||
			s[240] < sc || s[241] < sc || s[242] < sc || s[243] < sc ||
			s[244] < sc || s[245] < sc || s[246] < sc || s[247] < sc ||
			s[248] < sc || s[249] < sc || s[250] < sc || s[251] < sc ||
			s[252] < sc || s[253] < sc || s[254] < sc || s[255] < sc {
			return true
		}
		s = s[256:]
	}
	for len(s) >= 64 {
		if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
			s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc ||
			s[8] < sc || s[9] < sc || s[10] < sc || s[11] < sc ||
			s[12] < sc || s[13] < sc || s[14] < sc || s[15] < sc ||
			s[16] < sc || s[17] < sc || s[18] < sc || s[19] < sc ||
			s[20] < sc || s[21] < sc || s[22] < sc || s[23] < sc ||
			s[24] < sc || s[25] < sc || s[26] < sc || s[27] < sc ||
			s[28] < sc || s[29] < sc || s[30] < sc || s[31] < sc ||
			s[32] < sc || s[33] < sc || s[34] < sc || s[35] < sc ||
			s[36] < sc || s[37] < sc || s[38] < sc || s[39] < sc ||
			s[40] < sc || s[41] < sc || s[42] < sc || s[43] < sc ||
			s[44] < sc || s[45] < sc || s[46] < sc || s[47] < sc ||
			s[48] < sc || s[49] < sc || s[50] < sc || s[51] < sc ||
			s[52] < sc || s[53] < sc || s[54] < sc || s[55] < sc ||
			s[56] < sc || s[57] < sc || s[58] < sc || s[59] < sc ||
			s[60] < sc || s[61] < sc || s[62] < sc || s[63] < sc {
			return true
		}
		s = s[64:]
	}
	for len(s) >= 32 {
		if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
			s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc ||
			s[8] < sc || s[9] < sc || s[10] < sc || s[11] < sc ||
			s[12] < sc || s[13] < sc || s[14] < sc || s[15] < sc ||
			s[16] < sc || s[17] < sc || s[18] < sc || s[19] < sc ||
			s[20] < sc || s[21] < sc || s[22] < sc || s[23] < sc ||
			s[24] < sc || s[25] < sc || s[26] < sc || s[27] < sc ||
			s[28] < sc || s[29] < sc || s[30] < sc || s[31] < sc {
			return true
		}
		s = s[32:]
	}
	for len(s) >= 16 {
		if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
			s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc ||
			s[8] < sc || s[9] < sc || s[10] < sc || s[11] < sc ||
			s[12] < sc || s[13] < sc || s[14] < sc || s[15] < sc {
			return true
		}
		s = s[16:]
	}
	for len(s) >= 8 {
		if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
			s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc {
			return true
		}
		s = s[8:]
	}
	for i := 0; i < len(s); i++ {
		if s[i] < sc {
			return true
		}
	}
	return false
}

// ContainsCtrlCharBytes returns true if s contains any control character,
// otherwise returns false.
func ContainsCtrlCharBytes(s []byte) bool {
	const sc = 0x20

	if len(s) < 4 {
		for i := 0; i < len(s); i++ {
			if s[i] < sc {
				return true
			}
		}
		return false
	}
	if len(s) < 8 {
		if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc {
			return true
		}
		s = s[4:]
		for i := 0; i < len(s); i++ {
			if s[i] < sc {
				return true
			}
		}
		return false
	}
	if len(s) < 16 {
		if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
			s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc {
			return true
		}
		for i := 8; i < len(s); i++ {
			if s[i] < sc {
				return true
			}
		}
		return false
	}
	if len(s) < 32 {
		if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
			s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc ||
			s[8] < sc || s[9] < sc || s[10] < sc || s[11] < sc ||
			s[12] < sc || s[13] < sc || s[14] < sc || s[15] < sc {
			return true
		}
		s = s[16:]
		for len(s) >= 8 {
			if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
				s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc {
				return true
			}
			s = s[8:]
		}
		for i := 0; i < len(s); i++ {
			if s[i] < sc {
				return true
			}
		}
		return false
	}
	if len(s) < 64 {
		if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
			s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc ||
			s[8] < sc || s[9] < sc || s[10] < sc || s[11] < sc ||
			s[12] < sc || s[13] < sc || s[14] < sc || s[15] < sc ||
			s[16] < sc || s[17] < sc || s[18] < sc || s[19] < sc ||
			s[20] < sc || s[21] < sc || s[22] < sc || s[23] < sc ||
			s[24] < sc || s[25] < sc || s[26] < sc || s[27] < sc ||
			s[28] < sc || s[29] < sc || s[30] < sc || s[31] < sc {
			return true
		}
		s = s[32:]
		for len(s) >= 16 {
			if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
				s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc ||
				s[8] < sc || s[9] < sc || s[10] < sc || s[11] < sc ||
				s[12] < sc || s[13] < sc || s[14] < sc || s[15] < sc {
				return true
			}
			s = s[16:]
		}
		for len(s) >= 8 {
			if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
				s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc {
				return true
			}
			s = s[8:]
		}
		for i := 0; i < len(s); i++ {
			if s[i] < sc {
				return true
			}
		}
		return false
	}
	if len(s) < 256 {
		if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
			s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc ||
			s[8] < sc || s[9] < sc || s[10] < sc || s[11] < sc ||
			s[12] < sc || s[13] < sc || s[14] < sc || s[15] < sc ||
			s[16] < sc || s[17] < sc || s[18] < sc || s[19] < sc ||
			s[20] < sc || s[21] < sc || s[22] < sc || s[23] < sc ||
			s[24] < sc || s[25] < sc || s[26] < sc || s[27] < sc ||
			s[28] < sc || s[29] < sc || s[30] < sc || s[31] < sc ||
			s[32] < sc || s[33] < sc || s[34] < sc || s[35] < sc ||
			s[36] < sc || s[37] < sc || s[38] < sc || s[39] < sc ||
			s[40] < sc || s[41] < sc || s[42] < sc || s[43] < sc ||
			s[44] < sc || s[45] < sc || s[46] < sc || s[47] < sc ||
			s[48] < sc || s[49] < sc || s[50] < sc || s[51] < sc ||
			s[52] < sc || s[53] < sc || s[54] < sc || s[55] < sc ||
			s[56] < sc || s[57] < sc || s[58] < sc || s[59] < sc ||
			s[60] < sc || s[61] < sc || s[62] < sc || s[63] < sc {
			return true
		}
		s = s[64:]
		for len(s) >= 32 {
			if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
				s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc ||
				s[8] < sc || s[9] < sc || s[10] < sc || s[11] < sc ||
				s[12] < sc || s[13] < sc || s[14] < sc || s[15] < sc ||
				s[16] < sc || s[17] < sc || s[18] < sc || s[19] < sc ||
				s[20] < sc || s[21] < sc || s[22] < sc || s[23] < sc ||
				s[24] < sc || s[25] < sc || s[26] < sc || s[27] < sc ||
				s[28] < sc || s[29] < sc || s[30] < sc || s[31] < sc {
				return true
			}
			s = s[32:]
		}
		for len(s) >= 16 {
			if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
				s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc ||
				s[8] < sc || s[9] < sc || s[10] < sc || s[11] < sc ||
				s[12] < sc || s[13] < sc || s[14] < sc || s[15] < sc {
				return true
			}
			s = s[16:]
		}
		for len(s) >= 8 {
			if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
				s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc {
				return true
			}
			s = s[8:]
		}
		for i := 0; i < len(s); i++ {
			if s[i] < sc {
				return true
			}
		}
		return false
	}
	if len(s) < 1024 {
		if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
			s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc ||
			s[8] < sc || s[9] < sc || s[10] < sc || s[11] < sc ||
			s[12] < sc || s[13] < sc || s[14] < sc || s[15] < sc ||
			s[16] < sc || s[17] < sc || s[18] < sc || s[19] < sc ||
			s[20] < sc || s[21] < sc || s[22] < sc || s[23] < sc ||
			s[24] < sc || s[25] < sc || s[26] < sc || s[27] < sc ||
			s[28] < sc || s[29] < sc || s[30] < sc || s[31] < sc ||
			s[32] < sc || s[33] < sc || s[34] < sc || s[35] < sc ||
			s[36] < sc || s[37] < sc || s[38] < sc || s[39] < sc ||
			s[40] < sc || s[41] < sc || s[42] < sc || s[43] < sc ||
			s[44] < sc || s[45] < sc || s[46] < sc || s[47] < sc ||
			s[48] < sc || s[49] < sc || s[50] < sc || s[51] < sc ||
			s[52] < sc || s[53] < sc || s[54] < sc || s[55] < sc ||
			s[56] < sc || s[57] < sc || s[58] < sc || s[59] < sc ||
			s[60] < sc || s[61] < sc || s[62] < sc || s[63] < sc ||
			s[64] < sc || s[65] < sc || s[66] < sc || s[67] < sc ||
			s[68] < sc || s[69] < sc || s[70] < sc || s[71] < sc ||
			s[72] < sc || s[73] < sc || s[74] < sc || s[75] < sc ||
			s[76] < sc || s[77] < sc || s[78] < sc || s[79] < sc ||
			s[80] < sc || s[81] < sc || s[82] < sc || s[83] < sc ||
			s[84] < sc || s[85] < sc || s[86] < sc || s[87] < sc ||
			s[88] < sc || s[89] < sc || s[90] < sc || s[91] < sc ||
			s[92] < sc || s[93] < sc || s[94] < sc || s[95] < sc ||
			s[96] < sc || s[97] < sc || s[98] < sc || s[99] < sc ||
			s[100] < sc || s[101] < sc || s[102] < sc || s[103] < sc ||
			s[104] < sc || s[105] < sc || s[106] < sc || s[107] < sc ||
			s[108] < sc || s[109] < sc || s[110] < sc || s[111] < sc ||
			s[112] < sc || s[113] < sc || s[114] < sc || s[115] < sc ||
			s[116] < sc || s[117] < sc || s[118] < sc || s[119] < sc ||
			s[120] < sc || s[121] < sc || s[122] < sc || s[123] < sc ||
			s[124] < sc || s[125] < sc || s[126] < sc || s[127] < sc ||
			s[128] < sc || s[129] < sc || s[130] < sc || s[131] < sc ||
			s[132] < sc || s[133] < sc || s[134] < sc || s[135] < sc ||
			s[136] < sc || s[137] < sc || s[138] < sc || s[139] < sc ||
			s[140] < sc || s[141] < sc || s[142] < sc || s[143] < sc ||
			s[144] < sc || s[145] < sc || s[146] < sc || s[147] < sc ||
			s[148] < sc || s[149] < sc || s[150] < sc || s[151] < sc ||
			s[152] < sc || s[153] < sc || s[154] < sc || s[155] < sc ||
			s[156] < sc || s[157] < sc || s[158] < sc || s[159] < sc ||
			s[160] < sc || s[161] < sc || s[162] < sc || s[163] < sc ||
			s[164] < sc || s[165] < sc || s[166] < sc || s[167] < sc ||
			s[168] < sc || s[169] < sc || s[170] < sc || s[171] < sc ||
			s[172] < sc || s[173] < sc || s[174] < sc || s[175] < sc ||
			s[176] < sc || s[177] < sc || s[178] < sc || s[179] < sc ||
			s[180] < sc || s[181] < sc || s[182] < sc || s[183] < sc ||
			s[184] < sc || s[185] < sc || s[186] < sc || s[187] < sc ||
			s[188] < sc || s[189] < sc || s[190] < sc || s[191] < sc ||
			s[192] < sc || s[193] < sc || s[194] < sc || s[195] < sc ||
			s[196] < sc || s[197] < sc || s[198] < sc || s[199] < sc ||
			s[200] < sc || s[201] < sc || s[202] < sc || s[203] < sc ||
			s[204] < sc || s[205] < sc || s[206] < sc || s[207] < sc ||
			s[208] < sc || s[209] < sc || s[210] < sc || s[211] < sc ||
			s[212] < sc || s[213] < sc || s[214] < sc || s[215] < sc ||
			s[216] < sc || s[217] < sc || s[218] < sc || s[219] < sc ||
			s[220] < sc || s[221] < sc || s[222] < sc || s[223] < sc ||
			s[224] < sc || s[225] < sc || s[226] < sc || s[227] < sc ||
			s[228] < sc || s[229] < sc || s[230] < sc || s[231] < sc ||
			s[232] < sc || s[233] < sc || s[234] < sc || s[235] < sc ||
			s[236] < sc || s[237] < sc || s[238] < sc || s[239] < sc ||
			s[240] < sc || s[241] < sc || s[242] < sc || s[243] < sc ||
			s[244] < sc || s[245] < sc || s[246] < sc || s[247] < sc ||
			s[248] < sc || s[249] < sc || s[250] < sc || s[251] < sc ||
			s[252] < sc || s[253] < sc || s[254] < sc || s[255] < sc {
			return true
		}
		s = s[256:]
		for len(s) >= 64 {
			if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
				s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc ||
				s[8] < sc || s[9] < sc || s[10] < sc || s[11] < sc ||
				s[12] < sc || s[13] < sc || s[14] < sc || s[15] < sc ||
				s[16] < sc || s[17] < sc || s[18] < sc || s[19] < sc ||
				s[20] < sc || s[21] < sc || s[22] < sc || s[23] < sc ||
				s[24] < sc || s[25] < sc || s[26] < sc || s[27] < sc ||
				s[28] < sc || s[29] < sc || s[30] < sc || s[31] < sc ||
				s[32] < sc || s[33] < sc || s[34] < sc || s[35] < sc ||
				s[36] < sc || s[37] < sc || s[38] < sc || s[39] < sc ||
				s[40] < sc || s[41] < sc || s[42] < sc || s[43] < sc ||
				s[44] < sc || s[45] < sc || s[46] < sc || s[47] < sc ||
				s[48] < sc || s[49] < sc || s[50] < sc || s[51] < sc ||
				s[52] < sc || s[53] < sc || s[54] < sc || s[55] < sc ||
				s[56] < sc || s[57] < sc || s[58] < sc || s[59] < sc ||
				s[60] < sc || s[61] < sc || s[62] < sc || s[63] < sc {
				return true
			}
			s = s[64:]
		}
		for len(s) >= 32 {
			if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
				s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc ||
				s[8] < sc || s[9] < sc || s[10] < sc || s[11] < sc ||
				s[12] < sc || s[13] < sc || s[14] < sc || s[15] < sc ||
				s[16] < sc || s[17] < sc || s[18] < sc || s[19] < sc ||
				s[20] < sc || s[21] < sc || s[22] < sc || s[23] < sc ||
				s[24] < sc || s[25] < sc || s[26] < sc || s[27] < sc ||
				s[28] < sc || s[29] < sc || s[30] < sc || s[31] < sc {
				return true
			}
			s = s[32:]
		}
		for len(s) >= 16 {
			if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
				s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc ||
				s[8] < sc || s[9] < sc || s[10] < sc || s[11] < sc ||
				s[12] < sc || s[13] < sc || s[14] < sc || s[15] < sc {
				return true
			}
			s = s[16:]
		}
		for len(s) >= 8 {
			if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
				s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc {
				return true
			}
			s = s[8:]
		}
		for i := 0; i < len(s); i++ {
			if s[i] < sc {
				return true
			}
		}
		return false
	}

	// Worst case
	for len(s) >= 1024 {
		if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
			s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc ||
			s[8] < sc || s[9] < sc || s[10] < sc || s[11] < sc ||
			s[12] < sc || s[13] < sc || s[14] < sc || s[15] < sc ||
			s[16] < sc || s[17] < sc || s[18] < sc || s[19] < sc ||
			s[20] < sc || s[21] < sc || s[22] < sc || s[23] < sc ||
			s[24] < sc || s[25] < sc || s[26] < sc || s[27] < sc ||
			s[28] < sc || s[29] < sc || s[30] < sc || s[31] < sc ||
			s[32] < sc || s[33] < sc || s[34] < sc || s[35] < sc ||
			s[36] < sc || s[37] < sc || s[38] < sc || s[39] < sc ||
			s[40] < sc || s[41] < sc || s[42] < sc || s[43] < sc ||
			s[44] < sc || s[45] < sc || s[46] < sc || s[47] < sc ||
			s[48] < sc || s[49] < sc || s[50] < sc || s[51] < sc ||
			s[52] < sc || s[53] < sc || s[54] < sc || s[55] < sc ||
			s[56] < sc || s[57] < sc || s[58] < sc || s[59] < sc ||
			s[60] < sc || s[61] < sc || s[62] < sc || s[63] < sc ||
			s[64] < sc || s[65] < sc || s[66] < sc || s[67] < sc ||
			s[68] < sc || s[69] < sc || s[70] < sc || s[71] < sc ||
			s[72] < sc || s[73] < sc || s[74] < sc || s[75] < sc ||
			s[76] < sc || s[77] < sc || s[78] < sc || s[79] < sc ||
			s[80] < sc || s[81] < sc || s[82] < sc || s[83] < sc ||
			s[84] < sc || s[85] < sc || s[86] < sc || s[87] < sc ||
			s[88] < sc || s[89] < sc || s[90] < sc || s[91] < sc ||
			s[92] < sc || s[93] < sc || s[94] < sc || s[95] < sc ||
			s[96] < sc || s[97] < sc || s[98] < sc || s[99] < sc ||
			s[100] < sc || s[101] < sc || s[102] < sc || s[103] < sc ||
			s[104] < sc || s[105] < sc || s[106] < sc || s[107] < sc ||
			s[108] < sc || s[109] < sc || s[110] < sc || s[111] < sc ||
			s[112] < sc || s[113] < sc || s[114] < sc || s[115] < sc ||
			s[116] < sc || s[117] < sc || s[118] < sc || s[119] < sc ||
			s[120] < sc || s[121] < sc || s[122] < sc || s[123] < sc ||
			s[124] < sc || s[125] < sc || s[126] < sc || s[127] < sc ||
			s[128] < sc || s[129] < sc || s[130] < sc || s[131] < sc ||
			s[132] < sc || s[133] < sc || s[134] < sc || s[135] < sc ||
			s[136] < sc || s[137] < sc || s[138] < sc || s[139] < sc ||
			s[140] < sc || s[141] < sc || s[142] < sc || s[143] < sc ||
			s[144] < sc || s[145] < sc || s[146] < sc || s[147] < sc ||
			s[148] < sc || s[149] < sc || s[150] < sc || s[151] < sc ||
			s[152] < sc || s[153] < sc || s[154] < sc || s[155] < sc ||
			s[156] < sc || s[157] < sc || s[158] < sc || s[159] < sc ||
			s[160] < sc || s[161] < sc || s[162] < sc || s[163] < sc ||
			s[164] < sc || s[165] < sc || s[166] < sc || s[167] < sc ||
			s[168] < sc || s[169] < sc || s[170] < sc || s[171] < sc ||
			s[172] < sc || s[173] < sc || s[174] < sc || s[175] < sc ||
			s[176] < sc || s[177] < sc || s[178] < sc || s[179] < sc ||
			s[180] < sc || s[181] < sc || s[182] < sc || s[183] < sc ||
			s[184] < sc || s[185] < sc || s[186] < sc || s[187] < sc ||
			s[188] < sc || s[189] < sc || s[190] < sc || s[191] < sc ||
			s[192] < sc || s[193] < sc || s[194] < sc || s[195] < sc ||
			s[196] < sc || s[197] < sc || s[198] < sc || s[199] < sc ||
			s[200] < sc || s[201] < sc || s[202] < sc || s[203] < sc ||
			s[204] < sc || s[205] < sc || s[206] < sc || s[207] < sc ||
			s[208] < sc || s[209] < sc || s[210] < sc || s[211] < sc ||
			s[212] < sc || s[213] < sc || s[214] < sc || s[215] < sc ||
			s[216] < sc || s[217] < sc || s[218] < sc || s[219] < sc ||
			s[220] < sc || s[221] < sc || s[222] < sc || s[223] < sc ||
			s[224] < sc || s[225] < sc || s[226] < sc || s[227] < sc ||
			s[228] < sc || s[229] < sc || s[230] < sc || s[231] < sc ||
			s[232] < sc || s[233] < sc || s[234] < sc || s[235] < sc ||
			s[236] < sc || s[237] < sc || s[238] < sc || s[239] < sc ||
			s[240] < sc || s[241] < sc || s[242] < sc || s[243] < sc ||
			s[244] < sc || s[245] < sc || s[246] < sc || s[247] < sc ||
			s[248] < sc || s[249] < sc || s[250] < sc || s[251] < sc ||
			s[252] < sc || s[253] < sc || s[254] < sc || s[255] < sc ||
			s[256] < sc || s[257] < sc || s[258] < sc || s[259] < sc ||
			s[260] < sc || s[261] < sc || s[262] < sc || s[263] < sc ||
			s[264] < sc || s[265] < sc || s[266] < sc || s[267] < sc ||
			s[268] < sc || s[269] < sc || s[270] < sc || s[271] < sc ||
			s[272] < sc || s[273] < sc || s[274] < sc || s[275] < sc ||
			s[276] < sc || s[277] < sc || s[278] < sc || s[279] < sc ||
			s[280] < sc || s[281] < sc || s[282] < sc || s[283] < sc ||
			s[284] < sc || s[285] < sc || s[286] < sc || s[287] < sc ||
			s[288] < sc || s[289] < sc || s[290] < sc || s[291] < sc ||
			s[292] < sc || s[293] < sc || s[294] < sc || s[295] < sc ||
			s[296] < sc || s[297] < sc || s[298] < sc || s[299] < sc ||
			s[300] < sc || s[301] < sc || s[302] < sc || s[303] < sc ||
			s[304] < sc || s[305] < sc || s[306] < sc || s[307] < sc ||
			s[308] < sc || s[309] < sc || s[310] < sc || s[311] < sc ||
			s[312] < sc || s[313] < sc || s[314] < sc || s[315] < sc ||
			s[316] < sc || s[317] < sc || s[318] < sc || s[319] < sc ||
			s[320] < sc || s[321] < sc || s[322] < sc || s[323] < sc ||
			s[324] < sc || s[325] < sc || s[326] < sc || s[327] < sc ||
			s[328] < sc || s[329] < sc || s[330] < sc || s[331] < sc ||
			s[332] < sc || s[333] < sc || s[334] < sc || s[335] < sc ||
			s[336] < sc || s[337] < sc || s[338] < sc || s[339] < sc ||
			s[340] < sc || s[341] < sc || s[342] < sc || s[343] < sc ||
			s[344] < sc || s[345] < sc || s[346] < sc || s[347] < sc ||
			s[348] < sc || s[349] < sc || s[350] < sc || s[351] < sc ||
			s[352] < sc || s[353] < sc || s[354] < sc || s[355] < sc ||
			s[356] < sc || s[357] < sc || s[358] < sc || s[359] < sc ||
			s[360] < sc || s[361] < sc || s[362] < sc || s[363] < sc ||
			s[364] < sc || s[365] < sc || s[366] < sc || s[367] < sc ||
			s[368] < sc || s[369] < sc || s[370] < sc || s[371] < sc ||
			s[372] < sc || s[373] < sc || s[374] < sc || s[375] < sc ||
			s[376] < sc || s[377] < sc || s[378] < sc || s[379] < sc ||
			s[380] < sc || s[381] < sc || s[382] < sc || s[383] < sc ||
			s[384] < sc || s[385] < sc || s[386] < sc || s[387] < sc ||
			s[388] < sc || s[389] < sc || s[390] < sc || s[391] < sc ||
			s[392] < sc || s[393] < sc || s[394] < sc || s[395] < sc ||
			s[396] < sc || s[397] < sc || s[398] < sc || s[399] < sc ||
			s[400] < sc || s[401] < sc || s[402] < sc || s[403] < sc ||
			s[404] < sc || s[405] < sc || s[406] < sc || s[407] < sc ||
			s[408] < sc || s[409] < sc || s[410] < sc || s[411] < sc ||
			s[412] < sc || s[413] < sc || s[414] < sc || s[415] < sc ||
			s[416] < sc || s[417] < sc || s[418] < sc || s[419] < sc ||
			s[420] < sc || s[421] < sc || s[422] < sc || s[423] < sc ||
			s[424] < sc || s[425] < sc || s[426] < sc || s[427] < sc ||
			s[428] < sc || s[429] < sc || s[430] < sc || s[431] < sc ||
			s[432] < sc || s[433] < sc || s[434] < sc || s[435] < sc ||
			s[436] < sc || s[437] < sc || s[438] < sc || s[439] < sc ||
			s[440] < sc || s[441] < sc || s[442] < sc || s[443] < sc ||
			s[444] < sc || s[445] < sc || s[446] < sc || s[447] < sc ||
			s[448] < sc || s[449] < sc || s[450] < sc || s[451] < sc ||
			s[452] < sc || s[453] < sc || s[454] < sc || s[455] < sc ||
			s[456] < sc || s[457] < sc || s[458] < sc || s[459] < sc ||
			s[460] < sc || s[461] < sc || s[462] < sc || s[463] < sc ||
			s[464] < sc || s[465] < sc || s[466] < sc || s[467] < sc ||
			s[468] < sc || s[469] < sc || s[470] < sc || s[471] < sc ||
			s[472] < sc || s[473] < sc || s[474] < sc || s[475] < sc ||
			s[476] < sc || s[477] < sc || s[478] < sc || s[479] < sc ||
			s[480] < sc || s[481] < sc || s[482] < sc || s[483] < sc ||
			s[484] < sc || s[485] < sc || s[486] < sc || s[487] < sc ||
			s[488] < sc || s[489] < sc || s[490] < sc || s[491] < sc ||
			s[492] < sc || s[493] < sc || s[494] < sc || s[495] < sc ||
			s[496] < sc || s[497] < sc || s[498] < sc || s[499] < sc ||
			s[500] < sc || s[501] < sc || s[502] < sc || s[503] < sc ||
			s[504] < sc || s[505] < sc || s[506] < sc || s[507] < sc ||
			s[508] < sc || s[509] < sc || s[510] < sc || s[511] < sc ||
			s[512] < sc || s[513] < sc || s[514] < sc || s[515] < sc ||
			s[516] < sc || s[517] < sc || s[518] < sc || s[519] < sc ||
			s[520] < sc || s[521] < sc || s[522] < sc || s[523] < sc ||
			s[524] < sc || s[525] < sc || s[526] < sc || s[527] < sc ||
			s[528] < sc || s[529] < sc || s[530] < sc || s[531] < sc ||
			s[532] < sc || s[533] < sc || s[534] < sc || s[535] < sc ||
			s[536] < sc || s[537] < sc || s[538] < sc || s[539] < sc ||
			s[540] < sc || s[541] < sc || s[542] < sc || s[543] < sc ||
			s[544] < sc || s[545] < sc || s[546] < sc || s[547] < sc ||
			s[548] < sc || s[549] < sc || s[550] < sc || s[551] < sc ||
			s[552] < sc || s[553] < sc || s[554] < sc || s[555] < sc ||
			s[556] < sc || s[557] < sc || s[558] < sc || s[559] < sc ||
			s[560] < sc || s[561] < sc || s[562] < sc || s[563] < sc ||
			s[564] < sc || s[565] < sc || s[566] < sc || s[567] < sc ||
			s[568] < sc || s[569] < sc || s[570] < sc || s[571] < sc ||
			s[572] < sc || s[573] < sc || s[574] < sc || s[575] < sc ||
			s[576] < sc || s[577] < sc || s[578] < sc || s[579] < sc ||
			s[580] < sc || s[581] < sc || s[582] < sc || s[583] < sc ||
			s[584] < sc || s[585] < sc || s[586] < sc || s[587] < sc ||
			s[588] < sc || s[589] < sc || s[590] < sc || s[591] < sc ||
			s[592] < sc || s[593] < sc || s[594] < sc || s[595] < sc ||
			s[596] < sc || s[597] < sc || s[598] < sc || s[599] < sc ||
			s[600] < sc || s[601] < sc || s[602] < sc || s[603] < sc ||
			s[604] < sc || s[605] < sc || s[606] < sc || s[607] < sc ||
			s[608] < sc || s[609] < sc || s[610] < sc || s[611] < sc ||
			s[612] < sc || s[613] < sc || s[614] < sc || s[615] < sc ||
			s[616] < sc || s[617] < sc || s[618] < sc || s[619] < sc ||
			s[620] < sc || s[621] < sc || s[622] < sc || s[623] < sc ||
			s[624] < sc || s[625] < sc || s[626] < sc || s[627] < sc ||
			s[628] < sc || s[629] < sc || s[630] < sc || s[631] < sc ||
			s[632] < sc || s[633] < sc || s[634] < sc || s[635] < sc ||
			s[636] < sc || s[637] < sc || s[638] < sc || s[639] < sc ||
			s[640] < sc || s[641] < sc || s[642] < sc || s[643] < sc ||
			s[644] < sc || s[645] < sc || s[646] < sc || s[647] < sc ||
			s[648] < sc || s[649] < sc || s[650] < sc || s[651] < sc ||
			s[652] < sc || s[653] < sc || s[654] < sc || s[655] < sc ||
			s[656] < sc || s[657] < sc || s[658] < sc || s[659] < sc ||
			s[660] < sc || s[661] < sc || s[662] < sc || s[663] < sc ||
			s[664] < sc || s[665] < sc || s[666] < sc || s[667] < sc ||
			s[668] < sc || s[669] < sc || s[670] < sc || s[671] < sc ||
			s[672] < sc || s[673] < sc || s[674] < sc || s[675] < sc ||
			s[676] < sc || s[677] < sc || s[678] < sc || s[679] < sc ||
			s[680] < sc || s[681] < sc || s[682] < sc || s[683] < sc ||
			s[684] < sc || s[685] < sc || s[686] < sc || s[687] < sc ||
			s[688] < sc || s[689] < sc || s[690] < sc || s[691] < sc ||
			s[692] < sc || s[693] < sc || s[694] < sc || s[695] < sc ||
			s[696] < sc || s[697] < sc || s[698] < sc || s[699] < sc ||
			s[700] < sc || s[701] < sc || s[702] < sc || s[703] < sc ||
			s[704] < sc || s[705] < sc || s[706] < sc || s[707] < sc ||
			s[708] < sc || s[709] < sc || s[710] < sc || s[711] < sc ||
			s[712] < sc || s[713] < sc || s[714] < sc || s[715] < sc ||
			s[716] < sc || s[717] < sc || s[718] < sc || s[719] < sc ||
			s[720] < sc || s[721] < sc || s[722] < sc || s[723] < sc ||
			s[724] < sc || s[725] < sc || s[726] < sc || s[727] < sc ||
			s[728] < sc || s[729] < sc || s[730] < sc || s[731] < sc ||
			s[732] < sc || s[733] < sc || s[734] < sc || s[735] < sc ||
			s[736] < sc || s[737] < sc || s[738] < sc || s[739] < sc ||
			s[740] < sc || s[741] < sc || s[742] < sc || s[743] < sc ||
			s[744] < sc || s[745] < sc || s[746] < sc || s[747] < sc ||
			s[748] < sc || s[749] < sc || s[750] < sc || s[751] < sc ||
			s[752] < sc || s[753] < sc || s[754] < sc || s[755] < sc ||
			s[756] < sc || s[757] < sc || s[758] < sc || s[759] < sc ||
			s[760] < sc || s[761] < sc || s[762] < sc || s[763] < sc ||
			s[764] < sc || s[765] < sc || s[766] < sc || s[767] < sc ||
			s[768] < sc || s[769] < sc || s[770] < sc || s[771] < sc ||
			s[772] < sc || s[773] < sc || s[774] < sc || s[775] < sc ||
			s[776] < sc || s[777] < sc || s[778] < sc || s[779] < sc ||
			s[780] < sc || s[781] < sc || s[782] < sc || s[783] < sc ||
			s[784] < sc || s[785] < sc || s[786] < sc || s[787] < sc ||
			s[788] < sc || s[789] < sc || s[790] < sc || s[791] < sc ||
			s[792] < sc || s[793] < sc || s[794] < sc || s[795] < sc ||
			s[796] < sc || s[797] < sc || s[798] < sc || s[799] < sc ||
			s[800] < sc || s[801] < sc || s[802] < sc || s[803] < sc ||
			s[804] < sc || s[805] < sc || s[806] < sc || s[807] < sc ||
			s[808] < sc || s[809] < sc || s[810] < sc || s[811] < sc ||
			s[812] < sc || s[813] < sc || s[814] < sc || s[815] < sc ||
			s[816] < sc || s[817] < sc || s[818] < sc || s[819] < sc ||
			s[820] < sc || s[821] < sc || s[822] < sc || s[823] < sc ||
			s[824] < sc || s[825] < sc || s[826] < sc || s[827] < sc ||
			s[828] < sc || s[829] < sc || s[830] < sc || s[831] < sc ||
			s[832] < sc || s[833] < sc || s[834] < sc || s[835] < sc ||
			s[836] < sc || s[837] < sc || s[838] < sc || s[839] < sc ||
			s[840] < sc || s[841] < sc || s[842] < sc || s[843] < sc ||
			s[844] < sc || s[845] < sc || s[846] < sc || s[847] < sc ||
			s[848] < sc || s[849] < sc || s[850] < sc || s[851] < sc ||
			s[852] < sc || s[853] < sc || s[854] < sc || s[855] < sc ||
			s[856] < sc || s[857] < sc || s[858] < sc || s[859] < sc ||
			s[860] < sc || s[861] < sc || s[862] < sc || s[863] < sc ||
			s[864] < sc || s[865] < sc || s[866] < sc || s[867] < sc ||
			s[868] < sc || s[869] < sc || s[870] < sc || s[871] < sc ||
			s[872] < sc || s[873] < sc || s[874] < sc || s[875] < sc ||
			s[876] < sc || s[877] < sc || s[878] < sc || s[879] < sc ||
			s[880] < sc || s[881] < sc || s[882] < sc || s[883] < sc ||
			s[884] < sc || s[885] < sc || s[886] < sc || s[887] < sc ||
			s[888] < sc || s[889] < sc || s[890] < sc || s[891] < sc ||
			s[892] < sc || s[893] < sc || s[894] < sc || s[895] < sc ||
			s[896] < sc || s[897] < sc || s[898] < sc || s[899] < sc ||
			s[900] < sc || s[901] < sc || s[902] < sc || s[903] < sc ||
			s[904] < sc || s[905] < sc || s[906] < sc || s[907] < sc ||
			s[908] < sc || s[909] < sc || s[910] < sc || s[911] < sc ||
			s[912] < sc || s[913] < sc || s[914] < sc || s[915] < sc ||
			s[916] < sc || s[917] < sc || s[918] < sc || s[919] < sc ||
			s[920] < sc || s[921] < sc || s[922] < sc || s[923] < sc ||
			s[924] < sc || s[925] < sc || s[926] < sc || s[927] < sc ||
			s[928] < sc || s[929] < sc || s[930] < sc || s[931] < sc ||
			s[932] < sc || s[933] < sc || s[934] < sc || s[935] < sc ||
			s[936] < sc || s[937] < sc || s[938] < sc || s[939] < sc ||
			s[940] < sc || s[941] < sc || s[942] < sc || s[943] < sc ||
			s[944] < sc || s[945] < sc || s[946] < sc || s[947] < sc ||
			s[948] < sc || s[949] < sc || s[950] < sc || s[951] < sc ||
			s[952] < sc || s[953] < sc || s[954] < sc || s[955] < sc ||
			s[956] < sc || s[957] < sc || s[958] < sc || s[959] < sc ||
			s[960] < sc || s[961] < sc || s[962] < sc || s[963] < sc ||
			s[964] < sc || s[965] < sc || s[966] < sc || s[967] < sc ||
			s[968] < sc || s[969] < sc || s[970] < sc || s[971] < sc ||
			s[972] < sc || s[973] < sc || s[974] < sc || s[975] < sc ||
			s[976] < sc || s[977] < sc || s[978] < sc || s[979] < sc ||
			s[980] < sc || s[981] < sc || s[982] < sc || s[983] < sc ||
			s[984] < sc || s[985] < sc || s[986] < sc || s[987] < sc ||
			s[988] < sc || s[989] < sc || s[990] < sc || s[991] < sc ||
			s[992] < sc || s[993] < sc || s[994] < sc || s[995] < sc ||
			s[996] < sc || s[997] < sc || s[998] < sc || s[999] < sc ||
			s[1000] < sc || s[1001] < sc || s[1002] < sc || s[1003] < sc ||
			s[1004] < sc || s[1005] < sc || s[1006] < sc || s[1007] < sc ||
			s[1008] < sc || s[1009] < sc || s[1010] < sc || s[1011] < sc ||
			s[1012] < sc || s[1013] < sc || s[1014] < sc || s[1015] < sc ||
			s[1016] < sc || s[1017] < sc || s[1018] < sc || s[1019] < sc ||
			s[1020] < sc || s[1021] < sc || s[1022] < sc || s[1023] < sc {
			return true
		}
		s = s[1024:]
	}
	for len(s) >= 256 {
		if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
			s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc ||
			s[8] < sc || s[9] < sc || s[10] < sc || s[11] < sc ||
			s[12] < sc || s[13] < sc || s[14] < sc || s[15] < sc ||
			s[16] < sc || s[17] < sc || s[18] < sc || s[19] < sc ||
			s[20] < sc || s[21] < sc || s[22] < sc || s[23] < sc ||
			s[24] < sc || s[25] < sc || s[26] < sc || s[27] < sc ||
			s[28] < sc || s[29] < sc || s[30] < sc || s[31] < sc ||
			s[32] < sc || s[33] < sc || s[34] < sc || s[35] < sc ||
			s[36] < sc || s[37] < sc || s[38] < sc || s[39] < sc ||
			s[40] < sc || s[41] < sc || s[42] < sc || s[43] < sc ||
			s[44] < sc || s[45] < sc || s[46] < sc || s[47] < sc ||
			s[48] < sc || s[49] < sc || s[50] < sc || s[51] < sc ||
			s[52] < sc || s[53] < sc || s[54] < sc || s[55] < sc ||
			s[56] < sc || s[57] < sc || s[58] < sc || s[59] < sc ||
			s[60] < sc || s[61] < sc || s[62] < sc || s[63] < sc ||
			s[64] < sc || s[65] < sc || s[66] < sc || s[67] < sc ||
			s[68] < sc || s[69] < sc || s[70] < sc || s[71] < sc ||
			s[72] < sc || s[73] < sc || s[74] < sc || s[75] < sc ||
			s[76] < sc || s[77] < sc || s[78] < sc || s[79] < sc ||
			s[80] < sc || s[81] < sc || s[82] < sc || s[83] < sc ||
			s[84] < sc || s[85] < sc || s[86] < sc || s[87] < sc ||
			s[88] < sc || s[89] < sc || s[90] < sc || s[91] < sc ||
			s[92] < sc || s[93] < sc || s[94] < sc || s[95] < sc ||
			s[96] < sc || s[97] < sc || s[98] < sc || s[99] < sc ||
			s[100] < sc || s[101] < sc || s[102] < sc || s[103] < sc ||
			s[104] < sc || s[105] < sc || s[106] < sc || s[107] < sc ||
			s[108] < sc || s[109] < sc || s[110] < sc || s[111] < sc ||
			s[112] < sc || s[113] < sc || s[114] < sc || s[115] < sc ||
			s[116] < sc || s[117] < sc || s[118] < sc || s[119] < sc ||
			s[120] < sc || s[121] < sc || s[122] < sc || s[123] < sc ||
			s[124] < sc || s[125] < sc || s[126] < sc || s[127] < sc ||
			s[128] < sc || s[129] < sc || s[130] < sc || s[131] < sc ||
			s[132] < sc || s[133] < sc || s[134] < sc || s[135] < sc ||
			s[136] < sc || s[137] < sc || s[138] < sc || s[139] < sc ||
			s[140] < sc || s[141] < sc || s[142] < sc || s[143] < sc ||
			s[144] < sc || s[145] < sc || s[146] < sc || s[147] < sc ||
			s[148] < sc || s[149] < sc || s[150] < sc || s[151] < sc ||
			s[152] < sc || s[153] < sc || s[154] < sc || s[155] < sc ||
			s[156] < sc || s[157] < sc || s[158] < sc || s[159] < sc ||
			s[160] < sc || s[161] < sc || s[162] < sc || s[163] < sc ||
			s[164] < sc || s[165] < sc || s[166] < sc || s[167] < sc ||
			s[168] < sc || s[169] < sc || s[170] < sc || s[171] < sc ||
			s[172] < sc || s[173] < sc || s[174] < sc || s[175] < sc ||
			s[176] < sc || s[177] < sc || s[178] < sc || s[179] < sc ||
			s[180] < sc || s[181] < sc || s[182] < sc || s[183] < sc ||
			s[184] < sc || s[185] < sc || s[186] < sc || s[187] < sc ||
			s[188] < sc || s[189] < sc || s[190] < sc || s[191] < sc ||
			s[192] < sc || s[193] < sc || s[194] < sc || s[195] < sc ||
			s[196] < sc || s[197] < sc || s[198] < sc || s[199] < sc ||
			s[200] < sc || s[201] < sc || s[202] < sc || s[203] < sc ||
			s[204] < sc || s[205] < sc || s[206] < sc || s[207] < sc ||
			s[208] < sc || s[209] < sc || s[210] < sc || s[211] < sc ||
			s[212] < sc || s[213] < sc || s[214] < sc || s[215] < sc ||
			s[216] < sc || s[217] < sc || s[218] < sc || s[219] < sc ||
			s[220] < sc || s[221] < sc || s[222] < sc || s[223] < sc ||
			s[224] < sc || s[225] < sc || s[226] < sc || s[227] < sc ||
			s[228] < sc || s[229] < sc || s[230] < sc || s[231] < sc ||
			s[232] < sc || s[233] < sc || s[234] < sc || s[235] < sc ||
			s[236] < sc || s[237] < sc || s[238] < sc || s[239] < sc ||
			s[240] < sc || s[241] < sc || s[242] < sc || s[243] < sc ||
			s[244] < sc || s[245] < sc || s[246] < sc || s[247] < sc ||
			s[248] < sc || s[249] < sc || s[250] < sc || s[251] < sc ||
			s[252] < sc || s[253] < sc || s[254] < sc || s[255] < sc {
			return true
		}
		s = s[256:]
	}
	for len(s) >= 64 {
		if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
			s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc ||
			s[8] < sc || s[9] < sc || s[10] < sc || s[11] < sc ||
			s[12] < sc || s[13] < sc || s[14] < sc || s[15] < sc ||
			s[16] < sc || s[17] < sc || s[18] < sc || s[19] < sc ||
			s[20] < sc || s[21] < sc || s[22] < sc || s[23] < sc ||
			s[24] < sc || s[25] < sc || s[26] < sc || s[27] < sc ||
			s[28] < sc || s[29] < sc || s[30] < sc || s[31] < sc ||
			s[32] < sc || s[33] < sc || s[34] < sc || s[35] < sc ||
			s[36] < sc || s[37] < sc || s[38] < sc || s[39] < sc ||
			s[40] < sc || s[41] < sc || s[42] < sc || s[43] < sc ||
			s[44] < sc || s[45] < sc || s[46] < sc || s[47] < sc ||
			s[48] < sc || s[49] < sc || s[50] < sc || s[51] < sc ||
			s[52] < sc || s[53] < sc || s[54] < sc || s[55] < sc ||
			s[56] < sc || s[57] < sc || s[58] < sc || s[59] < sc ||
			s[60] < sc || s[61] < sc || s[62] < sc || s[63] < sc {
			return true
		}
		s = s[64:]
	}
	for len(s) >= 32 {
		if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
			s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc ||
			s[8] < sc || s[9] < sc || s[10] < sc || s[11] < sc ||
			s[12] < sc || s[13] < sc || s[14] < sc || s[15] < sc ||
			s[16] < sc || s[17] < sc || s[18] < sc || s[19] < sc ||
			s[20] < sc || s[21] < sc || s[22] < sc || s[23] < sc ||
			s[24] < sc || s[25] < sc || s[26] < sc || s[27] < sc ||
			s[28] < sc || s[29] < sc || s[30] < sc || s[31] < sc {
			return true
		}
		s = s[32:]
	}
	for len(s) >= 16 {
		if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
			s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc ||
			s[8] < sc || s[9] < sc || s[10] < sc || s[11] < sc ||
			s[12] < sc || s[13] < sc || s[14] < sc || s[15] < sc {
			return true
		}
		s = s[16:]
	}
	for len(s) >= 8 {
		if s[0] < sc || s[1] < sc || s[2] < sc || s[3] < sc ||
			s[4] < sc || s[5] < sc || s[6] < sc || s[7] < sc {
			return true
		}
		s = s[8:]
	}
	for i := 0; i < len(s); i++ {
		if s[i] < sc {
			return true
		}
	}
	return false
}
