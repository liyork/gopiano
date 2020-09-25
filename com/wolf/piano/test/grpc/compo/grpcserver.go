package compo

//
//import (
//	"fmt"
//	"google.golang.org/appengine/log"
//	"io"
//	"strings"
//)
//
//// 展示如何处理grpc的异常情况
//
//const (
//	ResponseFailCode = 0
//	ResponseSuccCode = 1
//	xmethod          = "x"
//	ymethod          = "z"
//)
//
//type SidecarRelyService struct {
//}
//
//func (s *SidecarRelyService) Subscribe(stream x.xx) error {
//	return s.optSidecarCrd(stream, xmethod)
//}
//
//func (s *SidecarRelyService) UnSubscribe(stream x.xx) error {
//	return s.optSidecarCrd(stream, ymethod)
//}
//
//func (s *SidecarRelyService) optSidecarCrd(stream x.xx, optMethod string) error {
//	firstConn := true
//	for {
//		in, err := stream.Recv()
//		if err == io.EOF {
//			return nil
//		}
//		if err != nil {
//			return err
//		}
//		if len(in.Ip) == 0 {
//			s.sendAndLog(stream, optMethod, ResponseFailCode,
//				fmt.Sprintf("SidecarRelyService %v err, request.Ip is nil", optMethod))
//			return nil
//		}
//		if len(in.Services) == 0 {
//			s.sendAndLog(stream, optMethod, ResponseFailCode,
//				fmt.Sprintf("SidecarRelyService %v err, request.Services is empty", optMethod))
//			return nil
//		}
//
//		nodeNs := s.transNs(in)
//		if firstConn {
//			defer s.clearSidecarCrd(nodeNs)
//		}
//
//		var revision string
//		if optMethod == xmethod {
//			revision, err = s.x.createOrAppendSidecarCrd(nodeNs, in.Services)
//		} else if optMethod == ymethod {
//			revision, err = s.x.removeOrDeleteSidecarCrd(nodeNs, in.Services)
//		}
//
//		if err != nil {
//			s.sendAndLog(stream, optMethod, ResponseFailCode,
//				fmt.Sprintf("x.optSidecarCrd err, nodeNs:%v, service:%v, revision:%v, err:%v",
//					nodeNs, in.Services, revision, err))
//			return nil
//		}
//
//		s.sendAndLog(stream, optMethod, ResponseSuccCode,
//			fmt.Sprintf("SidecarRelyService.%v success, nodeNs:%s, service:%v", optMethod, nodeNs, in.Services))
//		firstConn = false
//	}
//}
//
//func (s *SidecarRelyService) transNs(in *controlplane.SidecarRelyRequest) string {
//	return strings.Replace(in.Ip, ".", "-", -1)
//}
//
//func (s *SidecarRelyService) clearSidecarCrd(nodeNs string) {
//	log.Warnf("will clearSidecarCrd in nodeNs, nodeNs:%s", nodeNs)
//	err := s.x.deleteSidecarCrd(nodeNs)
//	if err != nil {
//		log.Errorf("clearSidecarCrd err, nodeNs:%s, err:%v", nodeNs, err)
//	}
//}
//
//func (s *SidecarRelyService) sendAndLog(stream controlplane.SidecarRelyService_SubscribeServer,
//	optMethod string, code int32, msg string) {
//	if err := stream.Send(&controlplane.SidecarRelyResponse{Code: code, Msg: msg}); err != nil {
//		log.Errorf("stream.Send err, msg:%v, err:%v", msg, err)
//		return
//	}
//	if code == ResponseSuccCode {
//		log.Infof("record %v succ msg:%v", optMethod, msg)
//	} else {
//		log.Errorf("record %v err msg:%v", optMethod, msg)
//	}
//}
