package main

// func (s *Server) Reply(seq int64) {
// 	Debug("Reply %d", seq)
// 	req, _, _, _ := s.getCertOrNew(seq).get()
// 	msg := &ReplyMsg{
// 		Seq:       seq,
// 		Timestamp: time.Now().UnixNano(),
// 		ClientId:  req.Req.ClientId,
// 		NodeId:    s.node.id,
// 		// Result:    req.Req.Operator,
// 	}
// 	digest := Sha256Digest(msg)
// 	sign := RsaSignWithSha256(digest, s.node.priKey)
// 	replyArgs := &ReplyArgs{
// 		Msg:  msg,
// 		Sign: sign,
// 	}
// 	var reply bool
// 	cliCli := s.getCliCli(req.Req.ClientId)
// 	if cliCli == nil {
// 		Warn("can't connect client %d", req.Req.ClientId)
// 		return
// 	}
// 	err := cliCli.Call("Client.ReplyRpc", replyArgs, &reply)
// 	if err != nil {
// 		Warn("Client.ReplyRpc error: %v", err)
// 		s.closeCliCli(req.Req.ClientId)
// 	}
// }

// func (s *Server) getCliCli(clientId int64) *rpc.Client {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
// 	cliCli, ok := s.id2cliCli[clientId]
// 	if !ok || cliCli == nil {
// 		node := GetNode(clientId)
// 		var err error
// 		cliCli, err = rpc.DialHTTP("tcp", node.addr)
// 		if err != nil {
// 			Warn("connect client %d error: %v", node.addr, err)
// 			return nil
// 		}
// 		s.id2cliCli[clientId] = cliCli
// 	}
// 	return cliCli
// }

// func (s *Server) closeCliCli(clientId int64) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
// 	Info("close connect with client %d", clientId)
// 	cliCli, ok := s.id2cliCli[clientId]
// 	if ok && cliCli != nil {
// 		cliCli = nil
// 		delete(s.id2cliCli, clientId)
// 	}
// }

// func (s *Server) CloseCliCliRPC(args *CloseCliCliArgs, reply *bool) error {
// 	s.closeCliCli(args.ClientId)
// 	*reply = true
// 	return nil
// }