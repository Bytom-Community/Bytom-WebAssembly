// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package common

import "math/big"

// Bytes2Big
//
func BytesToBig(data []byte) *big.Int {
	n := new(big.Int)
	n.SetBytes(data)

	return n
}
func Bytes2Big(data []byte) *big.Int { return BytesToBig(data) }

// Big to bytes
//
// Returns the bytes of a big integer with the size specified by **base**
// Attempts to pad the byte array with zeros.
func BigToBytes(num *big.Int, base int) []byte {
	ret := make([]byte, base/8)

	if len(num.Bytes()) > base/8 {
		return num.Bytes()
	}

	return append(ret[:len(ret)-len(num.Bytes())], num.Bytes()...)
}
