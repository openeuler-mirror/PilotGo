// Copyright 2015 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package attestation

import (
	"errors"
	"fmt"

	"github.com/google/go-tspi/tspi"
	"github.com/google/go-tspi/tspiconst"
)

// This is used for object use authorization
// We let it be well known because the
// additional security gained is insignificant
// compared to the inconveniences it can cause
var wellKnownSecret [20]byte

func pad(plaintext []byte, bsize int) ([]byte, error) {
	if bsize >= 256 {
		return nil, errors.New("bsize must be < 256")
	}
	pad := bsize - (len(plaintext) % bsize)
	if pad == 0 {
		pad = bsize
	}
	for i := 0; i < pad; i++ {
		plaintext = append(plaintext, byte(pad))
	}
	return plaintext, nil
}

// AIKChallengeResponse takes the output from GenerateChallenge along with the
// encrypted AIK key blob. The TPM then decrypts the asymmetric challenge with
// its EK in order to obtain the AES key, and uses the AES key to decrypt the
// symmetrically encrypted data. It verifies that this data blob corresponds
// to the AIK it was given, and if so hands back the secret contained within
// the symmetrically encrypted data.
func AIKChallengeResponse(context *tspi.Context, aikblob []byte, asymchallenge []byte, symchallenge []byte) (secret []byte, err error) {
	srk, err := context.LoadKeyByUUID(tspiconst.TSS_PS_TYPE_SYSTEM, tspi.TSS_UUID_SRK)
	if err != nil {
		return nil, err
	}
	srkpolicy, err := srk.GetPolicy(tspiconst.TSS_POLICY_USAGE)
	if err != nil {
		return nil, err
	}
	srkpolicy.SetSecret(tspiconst.TSS_SECRET_MODE_SHA1, wellKnownSecret[:])

	tpm := context.GetTPM()
	tpmpolicy, err := context.CreatePolicy(tspiconst.TSS_POLICY_USAGE)
	if err != nil {
		return nil, err
	}
	tpm.AssignPolicy(tpmpolicy)
	tpmpolicy.SetSecret(tspiconst.TSS_SECRET_MODE_SHA1, wellKnownSecret[:])

	aik, err := context.LoadKeyByBlob(srk, aikblob)
	if err != nil {
		return nil, err
	}
	secret, err = tpm.ActivateIdentity(aik, asymchallenge, symchallenge)

	return secret, err
}

// CreateAIK asks the TPM to generate an Attestation Identity Key. It returns
// the unencrypted public half of the AIK along with an encrypted blob
// containing both halves of the key, and any error.
func CreateAIK(context *tspi.Context) ([]byte, []byte, error) {
	n := make([]byte, 2048/8)
	for i := 0; i < 2048/8; i++ {
		n[i] = 0xff
	}

	srk, err := context.LoadKeyByUUID(tspiconst.TSS_PS_TYPE_SYSTEM, tspi.TSS_UUID_SRK)
	if err != nil {
		return nil, nil, err
	}
	keypolicy, err := srk.GetPolicy(tspiconst.TSS_POLICY_USAGE)
	if err != nil {
		return nil, nil, err
	}
	err = keypolicy.SetSecret(tspiconst.TSS_SECRET_MODE_SHA1, wellKnownSecret[:])
	if err != nil {
		return nil, nil, err
	}

	tpm := context.GetTPM()
	tpmpolicy, err := tpm.GetPolicy(tspiconst.TSS_POLICY_USAGE)
	if err != nil {
		return nil, nil, err
	}
	err = tpm.AssignPolicy(tpmpolicy)
	if err != nil {
		return nil, nil, err
	}
	err = tpmpolicy.SetSecret(tspiconst.TSS_SECRET_MODE_SHA1, wellKnownSecret[:])
	if err != nil {
		return nil, nil, err
	}

	pcakey, err := context.CreateKey(tspiconst.TSS_KEY_TYPE_LEGACY | tspiconst.TSS_KEY_SIZE_2048)
	if err != nil {
		return nil, nil, err
	}

	err = pcakey.SetModulus(n)
	if err != nil {
		return nil, nil, err
	}

	aik, err := context.CreateKey(tspiconst.TSS_KEY_TYPE_IDENTITY | tspiconst.TSS_KEY_SIZE_2048)
	if err != nil {
		return nil, nil, err
	}
	_, err = tpm.CollateIdentityRequest(srk, pcakey, aik)
	if err != nil {
		return nil, nil, err
	}

	pubkey, err := aik.GetPubKeyBlob()
	if err != nil {
		return nil, nil, err
	}
	blob, err := aik.GetKeyBlob()
	if err != nil {
		return nil, nil, err
	}

	_, err = aik.GetModulus()
	return pubkey, blob, nil
}

// GetQuote consumes a nonce and the aik blob and returns the quote,
// a signature, and any error
func GetQuote(context *tspi.Context, aikblob, nonce []byte) ([]byte, []byte, error) {
	srk, err := context.LoadKeyByUUID(tspiconst.TSS_PS_TYPE_SYSTEM, tspi.TSS_UUID_SRK)
	if err != nil {
		return nil, nil, fmt.Errorf("LoadKeyByUUID failed: %v", err)
	}
	srkpolicy, err := srk.GetPolicy(tspiconst.TSS_POLICY_USAGE)
	if err != nil {
		return nil, nil, fmt.Errorf("GetPolicy failed: %v", err)
	}
	srkpolicy.SetSecret(tspiconst.TSS_SECRET_MODE_SHA1, wellKnownSecret[:])

	tpm := context.GetTPM()

	aik, err := context.LoadKeyByBlob(srk, aikblob)
	if err != nil {
		return nil, nil, fmt.Errorf("LoadKeyByBlob failed: %v", err)
	}

	pcrs, err := context.CreatePCRs(tspiconst.TSS_PCRS_STRUCT_DEFAULT)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get a reference to PCRs: %v", err)
	}

	// GetQuote will only quote the PCR values that are selected by SetPCRs()
	// We want all the PCR values so we pass a slice of 0 through 23
	selectedPCRs := make([]int, 24)
	for i := 0; i < 24; i++ {
		selectedPCRs[i] = i
	}
	if err = pcrs.SetPCRs(selectedPCRs); err != nil {
		return nil, nil, fmt.Errorf("failed to set the PCR bitmap %v", err)
	}

	return tpm.GetQuote(aik, pcrs, nonce)
}

// GetEKCert reads the Endorsement Key certificate from the TPM's NVRAM and
// returns it, along with any error generated.
func GetEKCert(context *tspi.Context) (ekcert []byte, err error) {
	tpm := context.GetTPM()
	nv, err := context.CreateNV()
	if err != nil {
		return nil, err
	}
	policy, err := tpm.GetPolicy(tspiconst.TSS_POLICY_USAGE)
	if err != nil {
		return nil, err
	}
	policy.SetSecret(tspiconst.TSS_SECRET_MODE_SHA1, wellKnownSecret[:])
	nv.SetIndex(0x1000f000)
	nv.AssignPolicy(policy)
	data, err := nv.ReadValue(0, 5)
	if err != nil {
		return nil, err
	}

	tag := (uint)((uint)(data[0])<<8 | (uint)(data[1]))
	if tag != 0x1001 {
		return nil, fmt.Errorf("invalid tag: %x", tag)
	}

	if data[2] != 0 {
		return nil, fmt.Errorf("invalid certificate")
	}

	ekbuflen := (uint)(uint(data[3])<<8 | (uint)(data[4]))
	offset := (uint)(5)

	data, err = nv.ReadValue(offset, 2)

	tag = (uint)((uint)(data[0])<<8 | (uint)(data[1]))
	if tag == 0x1002 {
		offset += 2
		ekbuflen -= 2
	} else if data[0] != 0x30 {
		return nil, fmt.Errorf("invalid header: %x", tag)
	}

	ekoffset := (uint)(0)
	var ekbuf []byte
	for ekoffset < ekbuflen {
		length := (uint)(ekbuflen - ekoffset)
		if length > 128 {
			length = 128
		}
		data, err = nv.ReadValue(offset, length)
		if err != nil {
			return nil, err
		}

		ekbuf = append(ekbuf, data...)
		offset += length
		ekoffset += length
	}

	return ekbuf, nil
}
