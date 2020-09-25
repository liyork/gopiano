package compo

//
//import (
//	"context"
//	"google.golang.org/grpc"
//	"io"
//	"os"
//	"sync/atomic"
//	"time"
//)
//
//// 演示如何进行grpc客户端，带有retry功能
//
//const (
//	subScribeMethod   = "Subscribe"
//	unSubScribeMethod = "UnSubscribe"
//	registerMethod    = "Register"
//	request           = "Request"
//	receive           = "Receive"
//)
//
//type relyConfig struct {
//	Cancel      context.CancelFunc
//	stream      x.x
//	requestChan chan []string
//	methodName  string
//	direction   string
//}
//
//type relyCommon struct {
//	serviceClient x.SidecarRelyServiceClient
//	conn          *grpc.ClientConn
//
//	needReconChan          chan chan struct{}
//	callbackChan           chan struct{}
//	waitForCloseReconChan  chan struct{}
//	waitForCloseReconCount uint32
//
//	stopChan         chan struct{}
//	waitForStopChan  chan struct{}
//	waitForStopCount uint32
//}
//
//type SidecarRelyClient struct {
//	*relyCommon
//	subClient      *SidecarRelySubClient
//	unSubClient    *SidecarRelyUnSubClient
//	registerClient *SidecarRelyRegisterClient
//}
//
//type SidecarRelySubClient struct {
//	*relyCommon
//	relyConfig *relyConfig
//}
//
//type SidecarRelyUnSubClient struct {
//	*relyCommon
//	relyConfig *relyConfig
//}
//
//type SidecarRelyRegisterClient struct {
//	*relyCommon
//	requestChan chan *x.SidecarRegisterRequest
//	methodName  string
//	direction   string
//}
//
//// ======relyCommon
//
//func (s *relyCommon) connClose() {
//	if s.conn == nil {
//		return
//	}
//	s.conn.Close()
//	s.conn = nil
//}
//
//func (s *relyCommon) start(relyConfig *relyConfig) {
//	x.GoWithRecover(func() {
//		s.requestSidecarRely(relyConfig)
//	}, nil)
//	x.GoWithRecover(func() {
//		s.receiveSidecarRely(relyConfig)
//	}, nil)
//}
//
//func (s *SidecarRelyRegisterClient) start() {
//	x.GoWithRecover(func() {
//		s.register()
//	}, nil)
//}
//
//func (s *relyCommon) requestSidecarRely(relyConfig *relyConfig) {
//	relyConfig.direction = request
//	requestMethodName := relyConfig.methodName + " " + request
//	log.DefaultLogger.Infof("[sidecarrely client] %s routine invoke", requestMethodName)
//	localIp := os.Getenv(types.LocalIpKey)
//	s.invokeWithRetry(
//		func() {
//			services := <-relyConfig.requestChan
//
//			request := &x.SidecarRelyRequest{
//				Ip:       localIp,
//				Services: services,
//			}
//			stream := relyConfig.stream
//			if stream == nil {
//				return
//			}
//
//			if err := stream.Send(request); err != nil {
//				log.DefaultLogger.Errorf("[sidecarrely client] %s failed to send request, err:%v", requestMethodName, err)
//				s.failCleanup(relyConfig)
//				return
//			}
//			log.DefaultLogger.Infof("[sidecarrely client] %s success sent request, services:%v", requestMethodName, services)
//		},
//		func() bool {
//			return relyConfig.stream == nil
//		},
//		func() {
//			log.DefaultLogger.Warnf("[sidecarrely client] %s routine receive graceful shut down signal", requestMethodName)
//			close(relyConfig.requestChan)
//			s.waitForCloseReconChan <- struct{}{}
//		},
//	)
//}
//
//func (s *relyCommon) failCleanup(config *relyConfig) {
//	config.stream = nil
//}
//
//func (s *SidecarRelyRegisterClient) x() {
//	s.direction = request
//	requestMethodName := s.methodName + " " + request
//	log.DefaultLogger.Infof("[sidecarrely client] [%s] routine invoke", requestMethodName)
//	var isNeedReconnect = false
//	s.invokeWithRetry(
//		func() {
//			registerRequest, ok := <-s.requestChan
//			if !ok {
//				isNeedReconnect = true
//				return
//			}
//			ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
//
//			response, err := s.serviceClient.x(ctx, registerRequest)
//			if err != nil {
//				log.DefaultLogger.Errorf("[sidecarrely client] [%s] failed to register, err:%v", requestMethodName, err)
//				isNeedReconnect = true
//				return
//			}
//			isNeedReconnect = false
//			log.DefaultLogger.Infof("[sidecarrely client] [%s] success sent request, resp:%v", requestMethodName, response)
//		},
//		func() bool {
//			return isNeedReconnect
//		},
//		func() {
//			log.DefaultLogger.Warnf("[sidecarrely client] [%s] routine receive graceful shut down signal", requestMethodName)
//			s.waitForCloseReconChan <- struct{}{}
//		},
//	)
//}
//
//func (s *relyCommon) receiveSidecarRely(relyConfig *relyConfig) {
//	relyConfig.direction = "receive"
//	receiveMethodName := relyConfig.methodName + " " + relyConfig.direction
//	log.DefaultLogger.Infof("[sidecarrely client] %s routine invoke", receiveMethodName)
//	s.invokeWithRetry(
//		func() {
//			stream := relyConfig.stream
//			if stream == nil {
//				return
//			}
//			resp, err := stream.Recv()
//			if err == io.EOF {
//				log.DefaultLogger.Warnf("[sidecarrely client] [%s] receive stream eof", receiveMethodName)
//				s.failCleanup(relyConfig)
//			}
//			if err != nil {
//				log.DefaultLogger.Errorf("[sidecarrely client] [%s] receive stream err, err:%v", receiveMethodName, err)
//				s.failCleanup(relyConfig)
//				return
//			}
//			log.DefaultLogger.Infof("[sidecarrely client] [%s] success got message, msg:%v ", receiveMethodName, resp)
//		},
//		func() bool {
//			return relyConfig.stream == nil
//		},
//		func() {
//			log.DefaultLogger.Warnf("[sidecarrely client] %s routine receive graceful shut down signal", receiveMethodName)
//			s.waitForCloseReconChan <- struct{}{}
//		})
//}
//
//func (s *relyCommon) addWaitForStopCount() {
//	atomic.AddUint32(&s.waitForStopCount, 1)
//}
//
//func (s *relyCommon) addWaitForCloseReconCount() {
//	atomic.AddUint32(&s.waitForCloseReconCount, 1)
//}
//
//func (s *relyCommon) invokeWithRetry(onConnect func(), isNeedReconnect func() bool, stopCallback func()) {
//	s.addWaitForStopCount()
//	s.addWaitForCloseReconCount()
//	callbackChan := s.callbackChan
//	for {
//		if s.checkStop(
//			func() {
//				stopCallback()
//				s.waitForStopChan <- struct{}{}
//			}) {
//			return
//		}
//
//		if isNeedReconnect() {
//			s.needReconChan <- callbackChan
//			<-callbackChan
//			callbackChan = s.callbackChan
//		}
//
//		onConnect()
//	}
//}
//func (s *relyCommon) checkStop(closeCallback func()) bool {
//	select {
//	case <-s.stopChan:
//		closeCallback()
//		return true
//	default:
//	}
//	return false
//}
//
//func (s *relyCommon) waitAndCloseNeedReconChan() {
//	log.DefaultLogger.Infof("before waitAndCloseNeedReconChan....%d", s.waitForCloseReconCount)
//	for i := uint32(1); i <= s.waitForCloseReconCount; i++ {
//		<-s.waitForCloseReconChan
//	}
//	log.DefaultLogger.Infof("after waitAndCloseNeedReconChan....%d", s.waitForCloseReconCount)
//	close(s.waitForCloseReconChan)
//	close(s.needReconChan)
//}
//
//func (s *relyCommon) makeRelyConfig(
//	relyService func(context.Context) (x.SidecarRelyService_SubscribeClient, error),
//	methodName string) *relyConfig {
//	out := &relyConfig{
//		requestChan: make(chan []string),
//		methodName:  methodName,
//	}
//	s.createStream(out, relyService)
//	return out
//}
//
//func (s *relyCommon) createStream(relyConfig *relyConfig,
//	relyService func(context.Context) (x.SidecarRelyService_SubscribeClient, error)) bool {
//	ctx, cancel := context.WithCancel(context.Background())
//	relyConfig.Cancel = cancel
//	stream, err := relyService(ctx)
//	if err != nil {
//		log.DefaultLogger.Warnf("[sidecarrely client] %s invoke relyService err, err:%v", relyConfig.methodName, err)
//		s.connClose()
//		return false
//	}
//
//	relyConfig.stream = stream
//	log.DefaultLogger.Infof("[sidecarrely client] %s invoke relyService succ", relyConfig.methodName)
//	return true
//}
//
//// ======SidecarRelySubClient
//
//func (s *SidecarRelySubClient) makeSidecarRelyConfig() {
//	s.relyConfig = s.makeRelyConfig(s.relyService, subScribeMethod)
//}
//
//func (s *SidecarRelySubClient) relyService(ctx context.Context) (x.SidecarRelyService_SubscribeClient, error) {
//	return s.serviceClient.Subscribe(ctx)
//}
//
//func (s *SidecarRelyRegisterClient) makeSidecarRelyConfig() {
//	s.methodName = registerMethod
//	s.requestChan = make(chan *x.SidecarRegisterRequest, 0)
//}
//
//// ======SidecarRelyUnSubClient
//
//func (s *SidecarRelyUnSubClient) makeSidecarRelyConfig() {
//	s.relyConfig = s.makeRelyConfig(s.relyService, unSubScribeMethod)
//}
//
//func (s *SidecarRelyUnSubClient) relyService(ctx context.Context) (x.SidecarRelyService_SubscribeClient, error) {
//	return s.serviceClient.UnSubscribe(ctx)
//}
//
//// ======SidecarRelyClient
//
//func (s *SidecarRelyClient) makeSidecarRelySubClient() *SidecarRelySubClient {
//	return &SidecarRelySubClient{relyCommon: s.relyCommon}
//}
//
//func (s *SidecarRelyClient) makeSidecarRelyUnSubClient() *SidecarRelyUnSubClient {
//	return &SidecarRelyUnSubClient{relyCommon: s.relyCommon}
//}
//func (s *SidecarRelyClient) makeSidecarRelyRegisterClient() *SidecarRelyRegisterClient {
//	return &SidecarRelyRegisterClient{relyCommon: s.relyCommon}
//}
//
//func (s *SidecarRelyClient) Start() error {
//	log.DefaultLogger.Infof("[sidecarrely client] start")
//	s.createRelyCommon()
//	s.createRelyServiceClient()
//	s.createSubClient()
//	s.createUnSubClient()
//	s.createRegisterClient()
//	s.reconnectClient()
//
//	return nil
//}
//
//func (s *SidecarRelyClient) createSubClient() {
//	client := s.makeSidecarRelySubClient()
//	client.makeSidecarRelyConfig()
//	client.start(client.relyConfig)
//	s.subClient = client
//}
//func (s *SidecarRelyClient) createUnSubClient() {
//	client := s.makeSidecarRelyUnSubClient()
//	client.makeSidecarRelyConfig()
//	client.start(client.relyConfig)
//	s.unSubClient = client
//}
//
//func (s *SidecarRelyClient) createRegisterClient() {
//	client := s.makeSidecarRelyRegisterClient()
//	client.makeSidecarRelyConfig()
//	client.start()
//	s.registerClient = client
//}
//
//func (s *SidecarRelyClient) createRelyServiceClient() {
//	sidecarRelyServer := os.Getenv(types.SIDECARELYSERVER)
//	conn, err := grpc.Dial(sidecarRelyServer, grpc.WithInsecure())
//	if err != nil {
//		log.DefaultLogger.Warnf("[sidecarrely client] can not connect server, sidecarRelyServer:%s, err:%v", sidecarRelyServer, err)
//		return
//	}
//	serviceClient := controlplane.NewSidecarRelyServiceClient(conn)
//
//	s.relyCommon.serviceClient = serviceClient
//	s.relyCommon.conn = conn
//}
//
//func (s *SidecarRelyClient) reconnectClient() {
//	utils.GoWithRecover(func() {
//		s.reConnRelyClient()
//	}, nil)
//}
//
//func (s *SidecarRelyClient) createRelyCommon() {
//	s.relyCommon = &relyCommon{
//		needReconChan:         make(chan chan struct{}),
//		callbackChan:          make(chan struct{}),
//		waitForCloseReconChan: make(chan struct{}),
//		stopChan:              make(chan struct{}),
//		waitForStopChan:       make(chan struct{}),
//		waitForStopCount:      0,
//	}
//}
//
//func (s *SidecarRelyClient) closeSidecarRelyClient() {
//	s.closeSubStream()
//	s.closeUnSubStream()
//
//	if nil == s.subClient.relyConfig.stream && nil == s.unSubClient.relyConfig.stream {
//		log.DefaultLogger.Warnf("[sidecarrely client] closeSidecarRelyClient subclient and unsubclient is nil")
//	}
//
//	s.relyCommon.connClose()
//	log.DefaultLogger.Infof("[sidecarrely client] reConnRelyClient closed stream")
//}
//
//func (s *SidecarRelyClient) closeUnSubStream() {
//	unSubConfig := s.unSubClient.relyConfig
//	unSubConfig.stream = nil
//	unSubConfig.Cancel()
//}
//
//func (s *SidecarRelyClient) closeSubStream() {
//	subConfig := s.subClient.relyConfig
//	subConfig.stream = nil
//	subConfig.Cancel()
//}
//
//func (s *SidecarRelyClient) reConnRelyClient() {
//	s.addWaitForStopCount()
//	var lastCallbackChan chan struct{}
//	for {
//		select {
//		case curCallbackChan, ok := <-s.needReconChan:
//			if !ok {
//				log.DefaultLogger.Warnf("[sidecarrely client] reConnRelyClient needReconChan has closed")
//				s.waitForStopChan <- struct{}{}
//				return
//			}
//
//			if lastCallbackChan == curCallbackChan {
//				log.DefaultLogger.Warnf("[sidecarrely client] reConnRelyClient ignore same reconClient request")
//				continue
//			}
//
//			lastCallbackChan = curCallbackChan
//			s.closeSidecarRelyClient()
//
//			interval := time.Second
//			for {
//				if s.checkStop(func() {
//					close(curCallbackChan)
//					log.DefaultLogger.Warnf("[sidecarrely client] reConnRelyClient routine reconnecting receive graceful shut down signal")
//				}) {
//					break
//				}
//
//				s.createRelyServiceClient()
//				isCreateSubSucc := s.createStream(s.subClient.relyConfig, s.subClient.relyService)
//				isCreateUnSubSucc := s.createStream(s.unSubClient.relyConfig, s.unSubClient.relyService)
//
//				if isCreateSubSucc && isCreateUnSubSucc {
//					s.callbackChan = make(chan struct{})
//					close(curCallbackChan)
//					break
//				}
//
//				log.DefaultLogger.Infof("[sidecarrely client] reconnectClient failed, retry after %v", interval)
//				time.Sleep(interval + time.Duration(rand.Intn(1000))*time.Millisecond)
//				interval = v2.ComputeInterval(interval)
//			}
//		}
//	}
//}
//
//func (s *SidecarRelyClient) Subscribe(services []string) bool {
//	client := s.subClient
//	if nil == client {
//		log.DefaultLogger.Warnf("subscribe fail, s.subClient is nil")
//		return false
//	}
//	config := client.relyConfig
//	config.requestChan <- services
//	return true
//}
//
//func (s *SidecarRelyClient) UnSubscribe(services []string) bool {
//	client := s.unSubClient
//	if nil == client {
//		log.DefaultLogger.Warnf("unSubClient fail, s.unSubClient is nil")
//		return false
//	}
//	config := client.relyConfig
//	config.requestChan <- services
//	return true
//}
//
//func (s *SidecarRelyClient) Register(serviceName, alias string, port int, protocol int, attrs map[string]string) bool {
//	client := s.registerClient
//	if nil == client {
//		log.DefaultLogger.Warnf("register fail, s.registerClient is nil")
//		return false
//	}
//
//	registerRequest := &controlplane.SidecarRegisterRequest{
//		ServiceName: serviceName,
//		Alias:       alias,
//		Port:        int32(port),
//		Protocol:    int32(protocol),
//		Attrs:       attrs,
//	}
//	client.requestChan <- registerRequest
//	return true
//}
//
//func (s *SidecarRelyClient) Stop() {
//	log.DefaultLogger.Infof("prepare to stop sidecarrely client")
//	close(s.stopChan)
//	log.DefaultLogger.Infof("before s.waitForStopCount:%d", s.waitForStopCount)
//
//	closeConnChan := make(chan int)
//	go func() {
//		for {
//			select {
//			case <-closeConnChan:
//				log.DefaultLogger.Infof("reveive closeConnChan return")
//				return
//			default:
//				time.Sleep(1 * time.Second)
//				s.closeSidecarRelyClient()
//				s.closeRegister()
//				log.DefaultLogger.Infof("Stop closeSidecarRelyClient...")
//			}
//		}
//	}()
//	s.waitAndCloseNeedReconChan()
//	s.waitForActiveChan()
//
//	closeConnChan <- 1
//	log.DefaultLogger.Infof("sidecarrely client stop finish")
//}
//
//func (s *SidecarRelyClient) waitForActiveChan() {
//	for i := uint32(1); i <= s.waitForStopCount; i++ {
//		log.DefaultLogger.Infof("before <-s.waitForStopChan:%d", i)
//		<-s.waitForStopChan
//		log.DefaultLogger.Infof("end <-s.waitForStopChan:%d", i)
//	}
//	close(s.waitForStopChan)
//}
//
//func (s *SidecarRelyClient) closeRegister() {
//	if s.registerClient.requestChan != nil {
//		close(s.registerClient.requestChan)
//	}
//	s.registerClient.requestChan = nil
//}
//
////  ======
//
//func MakeSidecarRelyClient() *SidecarRelyClient {
//	return &SidecarRelyClient{}
//}
