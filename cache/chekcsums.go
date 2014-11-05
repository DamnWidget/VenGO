/*
   Copyright (C) 2014  Oscar Campos <oscar.campos@member.fsf.org>

   This program is free software; you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation; either version 2 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License along
   with this program; if not, write to the Free Software Foundation, Inc.,
   51 Franklin Street, Fifth Floor, Boston, MA 02110-1301 USA.

   See LICENSE file for more details.
*/

package cache

import "fmt"

var checksums map[string]string = map[string]string{
	// sources

	// old
	"1.1":   "a464704ebbbdd552a39b5f9429b059c117d165b3",
	"1.1.1": "f365aed8183e487a48a66ace7bf36e5974dffbb3",
	"1.1.2": "f5ab02bbfb0281b6c19520f44f7bc26f9da563fb",
	"1.2":   "7dd2408d40471aeb30a9e0b502c6717b5bf383a5",
	"1.2.1": "6a4b9991eddd8039438438d6aa25126ab7e07f2f",

	// stable
	"1.2.2": "3ce0ac4db434fc1546fec074841ff40dc48c1167",
	"1.3":   "9f9dfcbcb4fa126b2b66c0830dc733215f2f056e",
	"1.3.1": "bc296c9c305bacfbd7bff9e1b54f6f66ae421e6e",
	"1.3.2": "67d3a692588c259f9fe9dca5b80109e5b99271df",
	"1.3.3": "b54b7deb7b7afe9f5d9a3f5dd830c7dede35393a",

	// unstable
	"1.4beta1": "f2fece0c9f9cdc6e8a85ab56b7f1ffcb57c3e7cd",
	"1.3rc2":   "53a5b75c8bb2399c36ed8fe14f64bd2df34ca4d9",
	"1.3rc1":   "6a9dac2e65c07627fe51899e0031e298560b0097",

	// binaries

	// old
	"1.1.linux-amd64":   "58087060236614f13d7089d798055b4255ed02ec",
	"1.1.linux-386":     "8e48c8d6738316c07ad8a3d3714c30fb67afba1d",
	"1.1.freebsd-amd64": "d0c54a17e15747d67b9cda4f00fb4152233b8736",
	"1.1.freebsd-386":   "099d7aa5cee7800376031cd078d058b4503d88b0",
	"1.1.windows-amd64": "c09c4e97588572edbaa1b69b4d5e223279cb3f1f",
	"1.1.windows-386":   "01832d85c529c3104bf0a3702313730f76dde434",

	"1.1.1.linux-amd64":   "a35374bb290cb6c2a8e3d464ee192fa05a500018",
	"1.1.1.linux-386":     "05c427490583b675c0d19065d7860b52a8662ad1",
	"1.1.1.freebsd-amd64": "a4389c7508f705a54370366e6c7e7a457af9a147",
	"1.1.1.freebsd-386":   "6043ed2f216e9e7e9960b9372c6466caab36a937",
	"1.1.1.windows-amd64": "4f9dbc941af05dae5b9a01483644907ec2c4f875",
	"1.1.1.windows-386":   "dee46d45e27193c14e07e695aa4c13d7fdc9de9b",

	"1.1.2.linux-amd64":   "42634e25f98a5db1e8a2a8270c3604fcf8fed38d",
	"1.1.2.linux-386":     "42334112e5ba7fcd1a58de0a85ff2668e446cd0b",
	"1.1.2.freebsd-amd64": "8573ec318d1f7624233c51d69e94bc4d6bafb74c",
	"1.1.2.freebsd-386":   "ad7642eadb744ba6088c388b17c68eca860a43f7",
	"1.1.2.windows-amd64": "ab0a37ac5a71720e9a5a7b0a5f67afc432e872c6",
	"1.1.2.windows-386":   "f98d0fb3d31beb6587c265d220da7b17e302de4c",

	"1.2.linux-amd64":          "664e5025eae91412a96a10f4ed1a8af6f0f32b7d",
	"1.2.linux-386":            "63d67f66b766817b9760d0972451a7ee8dc05a2d",
	"1.2.freebsd-amd64":        "e77db02f5706ba31e45e385e51c13e4773c4d292",
	"1.2.freebsd-386":          "1e87f725570b67c21e176f4f77c4e39e7545fb2a",
	"1.2.darwin-amd64-osx10.6": "311a79a5ba8e258a4b17d941c0649bdae1741b13",
	"1.2.darwin-amd64-osx10.8": "71604d32cfaaab31b72707ca82287313e4bbb4e0",
	"1.2.darwin-386-osx10.6":   "22bca55e0477c11a532e7f1e5c4ff1828b1cb0a0",
	"1.2.darwin-386-osx10.8":   "d24cae107bbd4ae150fa41bef217ddc2e0cf0a39",
	"1.2.windows-amd64":        "3d48ddf918583c4ef0b2b817191b606959489a5b",
	"1.2.windows-386":          "23543124a1db2409bc12018c376f6e0eb8681f5f",

	"1.2.1.linux-amd64":          "7605f577ff6ac2d608a3a4e829b255ae2ebc8dcf",
	"1.2.1.linux-386":            "4f2611bdb51c1b1555da1e99e847e53ad9d52f42",
	"1.2.1.freebsd-amd64":        "7fd9444831bf1313733ef23404d85af882010206",
	"1.2.1.freebsd-386":          "380872dcf2ca30a8c942143d2d3ab2f6fc0144ec",
	"1.2.1.darwin-amd64-osx10.6": "3447ac38656c85ebc98a83d97663e06fa312bf65",
	"1.2.1.darwin-amd64-osx10.8": "da48bfa2f31c32d779dc99564324ec25a20f8e3c",
	"1.2.1.darwin-386-osx10.6":   "b8641da682fe7b4053f68e8d85cb5a78ab9049b9",
	"1.2.1.darwin-386-osx10.8":   "0f9ab33a35994f2ac603947fc99884295de2b5b3",
	"1.2.1.windows-amd64":        "60fea3a6b36dab5e79909c338fd3c1e92cccfbe1",
	"1.2.1.windows-386":          "29c10e196f06644b2e918cf2ee5a1118202f068c",

	// stable
	"1.2.2.linux-amd64":          "6bd151ca49c435462c8bf019477a6244b958ebb5",
	"1.2.2.linux-386":            "d16f892173b0589945d141cefb22adce57e3be9c",
	"1.2.2.freebsd-amd64":        "858744ab8ff9661d42940486af63d451853914a0",
	"1.2.2.freebsd-386":          "d226b8e1c3f75d31fa426df63aa776d7e08cddac",
	"1.2.2.darwin-amd64-osx10.6": "24c182718fd61b2621692dcdfc34937a6b5ee369",
	"1.2.2.darwin-amd64-osx10.8": "19be1eca8fc01b32bb6588a70773b84cdce6bed1",
	"1.2.2.darwin-386-osx10.6":   "360ec6cbfdec9257de029f918a881b9944718d7c",
	"1.2.2.darwin-386-osx10.8":   "4219b464e82e7c23d9dc02c193e7a0a28a09af1a",
	"1.2.2.windows-amd64":        "9ee22fe6c4d98124d582046aab465ab69eaab048",
	"1.2.2.windows-386":          "560bb33ec70ab733f31ff15f1a48fe35963983b9",

	"1.3.linux-amd64":          "b6b154933039987056ac307e20c25fa508a06ba6",
	"1.3.linux-386":            "22db33b0c4e242ed18a77b03a60582f8014fd8a6",
	"1.3.freebsd-amd64":        "71214bafabe2b5f52ee68afce96110031b446f0c",
	"1.3.freebsd-386":          "8afa9574140cdd5fc97883a06a11af766e7f0203",
	"1.3.darwin-amd64-osx10.6": "82ffcfb7962ca7114a1ee0a96cac51c53061ea05",
	"1.3.darwin-amd64-osx10.8": "8d768f10cd00e0b152490291d9cd6179a8ccf0a7",
	"1.3.darwin-386-osx10.6":   "159d2797bee603a80b829c4404c1fb2ee089cc00",
	"1.3.darwin-386-osx10.8":   "bade975462b5610781f6a9fe8ac13031b3fb7aa6",
	"1.3.windows-amd64":        "1e4888e1494aed7f6934acb5c4a1ffb0e9a022b1",
	"1.3.windows-386":          "e4e5279ce7d8cafdf210a522a70677d5b9c7589d",

	"1.3.1.linux-amd64":          "3af011cc19b21c7180f2604fd85fbc4ddde97143",
	"1.3.1.linux-386":            "36f87ce21cdb4cb8920bb706003d8655b4e1fc81",
	"1.3.1.freebsd-amd64":        "99e23fdd33860d837912e8647ed2a4b3d2b09d3c",
	"1.3.1.freebsd-386":          "586debe95542b3b56841f6bd2e5257e301a1ffdc",
	"1.3.1.darwin-amd64-osx10.6": "40716361d352c4b40252e79048e8bc084c3f3d1b",
	"1.3.1.darwin-amd64-osx10.8": "a7271cbdc25173d0f8da66549258ff65cca4bf06",
	"1.3.1.darwin-386-osx10.6":   "84f70a4c83be24cea696654a5b55331ea32f8a3f",
	"1.3.1.darwin-386-osx10.8":   "244dfba1f4239b8e2eb9c3abae5ad63fc32c807a",
	"1.3.1.windows-amd64":        "4548785cfa3bc228d18d2d06e39f58f0e4e014f1",
	"1.3.1.windows-386":          "64f99e40e79e93a622e73d7d55a5b8340f07747f",

	"1.3.2.linux-amd64":          "0e4b6120eee6d45e2e4374dac4fe7607df4cbe42",
	"1.3.2.linux-386":            "3cbfd62d401a6ca70779856fa8ad8c4d6c35c8cc",
	"1.3.2.freebsd-amd64":        "95b633f45156fbbe79076638f854e76b9cd01301",
	"1.3.2.freebsd-386":          "fea3ef264120b5c3b4c50a8929d56f47a8366503",
	"1.3.2.darwin-amd64-osx10.6": "36ca7e8ac9af12e70b1e01182c7ffc732ff3b876",
	"1.3.2.darwin-amd64-osx10.8": "323bf8088614d58fee2b4d2cb07d837063d7d77e",
	"1.3.2.darwin-386-osx10.6":   "d1652f6e0ed3063b7b43d2bc12981d927bc85deb",
	"1.3.2.darwin-386-osx10.8":   "d040c85698c749fdbe25e8568c4d71648a5e3a75",
	"1.3.2.windows-amd64":        "7f7147484b1bc9e52cf034de816146977d0137f6",
	"1.3.2.windows-386":          "86160c478436253f51241ac1905577d337577ce0",

	"1.3.3.linux-amd64":          "14068fbe349db34b838853a7878621bbd2b24646",
	"1.3.3.linux-386":            "9eb426d5505de55729e2656c03d85722795dd85e",
	"1.3.3.freebsd-amd64":        "8531ae5e745c887f8dad1a3f00ca873cfcace56e",
	"1.3.3.freebsd-386":          "875a5515dd7d3e5826c7c003bb2450f3129ccbad",
	"1.3.3.darwin-amd64-osx10.6": "dfe68de684f6e8d9c371d01e6d6a522efe3b8942",
	"1.3.3.darwin-amd64-osx10.8": "be686ec7ba68d588735cc2094ccab8bdd651de9e",
	"1.3.3.darwin-386-osx10.6":   "04b3e38549183e984f509c07ad40d8bcd577a702",
	"1.3.3.darwin-386-osx10.8":   "88f35d3327a84107aac4f2f24cb0883e5fdbe0e5",
	"1.3.3.windows-amd64":        "5f0b3b104d3db09edd32ef1d086ba20bafe01ada",
	"1.3.3.windows-386":          "ba99083b22e0b22b560bb2d28b9b99b405d01b6b",

	// unstable
	"1.4beta1.linux-amd64":          "d2712acdaa4469ce2dc57c112a70900667269ca0",
	"1.4beta1.linux-386":            "122ea6cae37d9b62c69efa3e21cc228e41006b75",
	"1.4beta1.freebsd-amd64":        "42fbd5336437dde85b34d774bfed111fe579db88",
	"1.4beta1.freebsd-386":          "65045b7a5d2a991a45b1e86ad11252bc84043651",
	"1.4beta1.darwin-amd64-osx10.6": "ad8798fe744bb119f0e8eeacf97be89763c5f12a",
	"1.4beta1.darwin-amd64-osx10.8": "e08df216d9761c970e438295129721ec8374654a",
	"1.4beta1.darwin-386-osx10.6":   "a360e7c8f1d528901e721d0cc716461f8a636823",
	"1.4beta1.darwin-386-osx10.8":   "d863907870e8e79850a7a725b398502afd1163d8",
	"1.4beta1.windows-amd64":        "386deea0a7c384178aedfe48e4ee2558a8cd43d8",
	"1.4beta1.windows-386":          "a6d75ca59b70226087104b514389e48d49854ed4",
}

// check if a given version is supported by VenGO to auto donwload/compile
// if the version is valid, it returns it's SHA1 fingerprint, error is
// returned otherwise
func Checksum(version string) (string, error) {
	if sha1, ok := checksums[version]; ok {
		return sha1, nil
	}
	return "", fmt.Errorf("%s is not a VenGO supported version you must donwload and compile it yourself", version)
}
