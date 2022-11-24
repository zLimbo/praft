package main

// func (s *Server) Reply(seq int64) {
// 	zlog.Debug("Reply %d", seq)
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
// 		zlog.Warn("can't connect client %d", req.Req.ClientId)
// 		return
// 	}
// 	err := cliCli.Call("Client.ReplyRpc", replyArgs, &reply)
// 	if err != nil {
// 		zlog.Warn("Client.ReplyRpc error: %v", err)
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
// 			zlog.Warn("connect client %d error: %v", node.addr, err)
// 			return nil
// 		}
// 		s.id2cliCli[clientId] = cliCli
// 	}
// 	return cliCli
// }

// func (s *Server) closeCliCli(clientId int64) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
// 	zlog.Info("close connect with client %d", clientId)
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

// func (s *Server) RequestRpc(args *RequestArgs, reply *RequestReply) error {
// 	// 放入请求队列直接返回，后续异步通知客户端

// 	zlog.Debug("RequestRpc, from: %d", args.Req.ClientId)
// 	//构造请求
// 	req := &RequestMsg{
// 		// Operator:  make([]byte, KConfig.BatchTxNum*KConfig.TxSize),
// 		TxSet:     ycsb.GenTxSet(ycsb.Wrate, KConfig.BatchTxNum),
// 		Timestamp: time.Now().UnixNano(),
// 		ClientId:  args.Req.ClientId,
// 	}
// 	node := GetNode(args.Req.ClientId)
// 	digest := Sha256Digest(req)
// 	sign := RsaSignWithSha256(digest, node.priKey)

// 	args = &RequestArgs{
// 		Req:  req,
// 		Sign: sign,
// 	}
// 	// 验证RequestMsg
// 	// node := GetNode(args.Req.ClientId)
// 	// digest := Sha256Digest(args.Req)
// 	// ok := RsaVerifyWithSha256(digest, args.Sign, node.pubKey)
// 	// if !ok {
// 	// 	zlog.Warn("RequestMsg verify error, from: %d", args.Req.ClientId)
// 	// 	reply.Ok = false
// 	// 	return nil
// 	// }

// 	// leader 分配seq
// 	seq := s.assignSeq()
// 	//主节点新建logCert，设置参数，并存储在seq2cert中
// 	s.getCertOrNew(seq).set(args, digest, s.view, s.node.id)
// 	//向共识线程发送开始共识信号
// 	s.seqCh <- seq

// 	// 返回信息
// 	reply.Seq = seq
// 	reply.Ok = true

// 	return nil
// }
