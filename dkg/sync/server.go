// Copyright © 2022 Obol Labs Inc.
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU General Public License as published by the Free
// Software Foundation, either version 3 of the License, or (at your option)
// any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of  MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for
// more details.
//
// You should have received a copy of the GNU General Public License along with
// this program.  If not, see <http://www.gnu.org/licenses/>.

package sync

import (
	"bufio"
	"context"

	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/network"
	"google.golang.org/protobuf/proto"

	"github.com/obolnetwork/charon/app/log"
	"github.com/obolnetwork/charon/app/z"
	pb "github.com/obolnetwork/charon/dkg/dkgpb/v1"
	"github.com/obolnetwork/charon/p2p"
)

const (
	syncProtoID = "dkg_v1.0"
	MsgSize     = 128
)

type Server struct {
	ctx        context.Context
	onFailure  func()
	tcpNode    host.Host
	peers      []p2p.Peer
	clientMsgs chan *pb.MsgSync
}

// AwaitAllConnected blocks until all peers have established a connection with this server or returns an error.
func (*Server) AwaitAllConnected() error {
	return nil
}

// AwaitAllShutdown blocks until all peers have successfully shutdown or returns an error.
// It may only be called after AwaitAllConnected.
func (*Server) AwaitAllShutdown() error {
	return nil
}

// NewServer registers a Stream Handler and returns a new Server instance.
func NewServer(ctx context.Context, tcpNode host.Host, peers []p2p.Peer, defHash []byte, onFailure func()) *Server {
	server := &Server{
		ctx:       ctx,
		tcpNode:   tcpNode,
		peers:     peers,
		onFailure: onFailure,
	}

	server.tcpNode.SetStreamHandler(syncProtoID, func(s network.Stream) {
		defer s.Close()

		buf := bufio.NewReader(s)
		b := make([]byte, MsgSize)

		// n is the number of bytes read from buffer, if n < MsgSize the other bytes will be 0
		n, err := buf.Read(b)
		if err != nil {
			log.Error(ctx, "Read client msg from stream", err)
			err = s.Reset()
			log.Error(ctx, "Stream reset", err)

			return
		}

		// Number of bytes that are read are the most important
		b = b[:n]

		msg := new(pb.MsgSync)
		if err := proto.Unmarshal(b, msg); err != nil {
			log.Error(ctx, "Unmarshal client msg", err)
			err = s.Reset()
			log.Error(ctx, "Stream reset", err)

			return
		}

		pID := s.Conn().RemotePeer()
		log.Debug(ctx, "Message received from client", z.Any("peer", p2p.PeerName(pID)))

		pubkey, err := pID.ExtractPublicKey()
		if err != nil {
			log.Error(ctx, "Get client public key", err)
			err = s.Reset()
			log.Error(ctx, "Stream reset", err)
		}

		ok, err := pubkey.Verify(defHash, msg.HashSignature)
		if err != nil {
			log.Error(ctx, "Verify defHash signature", err)
			err = s.Reset()
			log.Error(ctx, "Stream reset", err)
		}

		resp := &pb.MsgSyncResponse{
			SyncTimestamp: msg.Timestamp,
			Error:         "",
		}

		if !ok {
			resp.Error = "Invalid Signature"
		}

		resBytes, err := proto.Marshal(resp)
		if err != nil {
			log.Error(ctx, "Marshal server response", err)
			err = s.Reset()
			log.Error(ctx, "Stream reset", err)

			return
		}

		_, err = s.Write(resBytes)
		if err != nil {
			log.Error(ctx, "Send response to client", err)
			err = s.Reset()
			log.Error(ctx, "Stream reset", err)

			return
		}
	})

	return server
}