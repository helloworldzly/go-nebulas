// Copyright (C) 2017 go-nebulas authors
//
// This file is part of the go-nebulas library.
//
// the go-nebulas library is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// the go-nebulas library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with the go-nebulas library.  If not, see <http://www.gnu.org/licenses/>.
//

package core

import (
	"reflect"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/nebulasio/go-nebulas/core/pb"
	"github.com/nebulasio/go-nebulas/util"
)

func TestTransaction(t *testing.T) {
	type fields struct {
		hash      Hash
		from      *Address
		to        *Address
		value     *util.Uint128
		nonce     uint64
		timestamp time.Time
		alg       uint8
		data      []byte
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			"full struct",
			fields(fields{
				[]byte("123455"),
				&Address{[]byte("1235")},
				&Address{[]byte("1245")},
				util.NewUint128(),
				456,
				time.Now(),
				12,
				[]byte("hwllo"),
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &Transaction{
				hash:      tt.fields.hash,
				from:      tt.fields.from,
				to:        tt.fields.to,
				value:     tt.fields.value,
				nonce:     tt.fields.nonce,
				timestamp: tt.fields.timestamp,
				alg:       tt.fields.alg,
				data:      tt.fields.data,
			}
			msg, _ := tx.ToProto()
			ir, _ := proto.Marshal(msg)
			ntx := new(Transaction)
			nMsg := new(corepb.Transaction)
			proto.Unmarshal(ir, nMsg)
			ntx.FromProto(nMsg)
			ntx.timestamp = tx.timestamp
			if !reflect.DeepEqual(tx, ntx) {
				t.Errorf("Transaction.Serialize() = %v, want %v", *tx, *ntx)
			}
		})
	}
}
