package network

import (
	"context"
	tls "crypto/tls"
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/choria-io/go-choria/choria"
	"github.com/choria-io/go-choria/srvcache"
	gomock "github.com/golang/mock/gomock"
	logrus "github.com/sirupsen/logrus"

	"github.com/choria-io/go-choria/config"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	nats "github.com/nats-io/nats.go"
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Broker/Network")
}

var _ = Describe("Network Broker", func() {
	var (
		mockctl   *gomock.Controller
		cfg       *config.Config
		clientcfg *config.Config
		fw        *MockChoriaFramework
		bi        *MockBuildInfoProvider
		srv       *Server
		err       error
		logger    *logrus.Entry
	)

	BeforeEach(func() {
		mockctl = gomock.NewController(GinkgoT())
		bi = NewMockBuildInfoProvider(mockctl)
		fw = NewMockChoriaFramework(mockctl)

		cfg, err = config.NewDefaultConfig()
		Expect(err).ToNot(HaveOccurred())
		clientcfg, err = config.NewDefaultConfig()
		Expect(err).ToNot(HaveOccurred())

		cfg.Choria.SSLDir = "testdata/ssl"

		logger = logrus.NewEntry(logrus.New())
		logger.Logger.SetLevel(logrus.DebugLevel)
		logger.Logger.Out = GinkgoWriter

		fw.EXPECT().Configuration().Return(cfg).AnyTimes()
		fw.EXPECT().Logger(gomock.Any()).Return(logger).AnyTimes()
		bi.EXPECT().MaxBrokerClients().Return(50000).AnyTimes()
	})

	AfterEach(func() {
		mockctl.Finish()
	})

	Describe("NewServer", func() {
		It("Should initialize the server correctly", func() {
			cfg.Choria.NetworkListenAddress = "0.0.0.0"
			cfg.Choria.NetworkClientPort = 8080
			cfg.Choria.NetworkWriteDeadline = time.Duration(10 * time.Second)
			cfg.LogLevel = "error"
			cfg.Choria.StatsPort = 8081
			cfg.Choria.StatsListenAddress = "192.168.1.1"
			cfg.Choria.NetworkPeerPort = 8082
			cfg.Choria.NetworkPeerUser = "bob"
			cfg.Choria.NetworkPeerPassword = "secret"
			cfg.Choria.NetworkPeers = []string{"nats://localhost:9000", "nats://localhost:9001", "nats://localhost:8082"}

			fw.EXPECT().NetworkBrokerPeers().Return(srvcache.NewServers(
				srvcache.NewServer("localhost", 9000, "nats"),
				srvcache.NewServer("localhost", 9001, "nats"),
				srvcache.NewServer("localhost", 8082, "nats"),
			), nil)

			fw.EXPECT().TLSConfig().Return(&tls.Config{}, nil)

			srv, err = NewServer(fw, bi, false)
			Expect(err).ToNot(HaveOccurred())

			Expect(srv.opts.Host).To(Equal("0.0.0.0"))
			Expect(srv.opts.Port).To(Equal(8080))
			Expect(srv.opts.Logtime).To(BeFalse())
			Expect(srv.opts.MaxConn).To(Equal(50000))
			Expect(srv.opts.WriteDeadline).To(Equal(time.Duration(10 * time.Second)))
			Expect(srv.opts.NoSigs).To(BeTrue())
			Expect(srv.opts.Debug).To(BeFalse())
			Expect(srv.opts.HTTPHost).To(Equal("192.168.1.1"))
			Expect(srv.opts.HTTPPort).To(Equal(8081))
			Expect(srv.opts.Cluster.Host).To(Equal("0.0.0.0"))
			Expect(srv.opts.Cluster.NoAdvertise).To(BeTrue())
			Expect(srv.opts.Cluster.Port).To(Equal(8082))
			Expect(srv.opts.Cluster.Username).To(Equal("bob"))
			Expect(srv.opts.Cluster.Password).To(Equal("secret"))
			Expect(srv.opts.Routes).To(HaveLen(2))
			Expect(srv.opts.Routes[0].Host).To(Equal("localhost:9000"))
			Expect(srv.opts.Routes[1].Host).To(Equal("localhost:9001"))
			Expect(srv.opts.TLS).To(BeTrue())
			Expect(srv.opts.TLSVerify).To(BeTrue())
			Expect(srv.opts.TLSTimeout).To(Equal(float64(2)))
			Expect(srv.opts.LeafNode.Host).To(Equal(""))
			Expect(srv.opts.LeafNode.Port).To(Equal(0))
		})

		It("Should initialize the server, with TLS correctly", func() {
			cfg.Choria.NetworkListenAddress = "localhost"
			cfg.Choria.NetworkClientPort = -1
			cfg.Choria.NetworkWriteDeadline = time.Duration(10 * time.Second)
			cfg.LogLevel = "debug"
			cfg.Choria.NetworkPeerPort = -1
			cfg.Choria.FileSecurityCA = "../../../../../providers/security/testdata/intermediate/ca.pem"
			cfg.Choria.FileSecurityCertificate = "../../../../../providers/security/testdata/intermediate/localhost.pem"
			cfg.Choria.FileSecurityKey = "../../../../../providers/security/testdata/intermediate/localhost-key.pem"

			clientcfg.Choria.FileSecurityCA = "../../providers/security/testdata/good/certs/ca.pem"
			clientcfg.Choria.FileSecurityCertificate = "../../providers/security/testdata/good/certs/rip.mcollective.pem"
			clientcfg.Choria.FileSecurityKey = "../../providers/security/testdata/good/private_keys/rip.mcollective.pem"

			lfw, err := choria.NewWithConfig(cfg)
			cfw, err := choria.NewWithConfig(clientcfg)

			srv, err = NewServer(lfw, bi, false)
			Expect(err).ToNot(HaveOccurred())

			ctx, cancel := context.WithCancel(context.Background())

			wg := sync.WaitGroup{}
			wg.Add(1)
			go srv.Start(ctx, &wg)
			if !srv.gnatsd.ReadyForConnections(5 * time.Second) {
				Fail("Broker didnt start")
			}

			tlsc, err := cfw.TLSConfig()
			Expect(err).ToNot(HaveOccurred())
			_, err = nats.Connect(srv.gnatsd.ClientURL(), nats.Secure(tlsc), nats.MaxReconnects(1))
			Expect(err.Error()).ToNot(HaveOccurred())

			cancel()

		})

		// It("Should support disabling TLS Verify", func() {
		// 	cfg.DisableTLSVerify = true

		// 	fw, err = choria.NewWithConfig(cfg)
		// 	Expect(err).ToNot(HaveOccurred())

		// 	srv, err = NewServer(fw, false)
		// 	Expect(err).ToNot(HaveOccurred())
		// 	Expect(srv.opts.TLSVerify).To(BeFalse())
		// })

		It("Should support disabling TLS", func() {
			cfg.DisableTLS = true
			fw.EXPECT().NetworkBrokerPeers().Return(srvcache.NewServers(), nil)

			srv, err = NewServer(fw, bi, false)
			Expect(err).ToNot(HaveOccurred())
			Expect(srv.opts.TLS).To(BeFalse())
		})

		It("Should support forcing client TLS on while framework TLS is off", func() {
			cfg.DisableTLS = true
			cfg.Choria.NetworkClientTLSForce = true

			fw.EXPECT().NetworkBrokerPeers().Return(srvcache.NewServers(), nil)
			fw.EXPECT().ValidateSecurity().Return([]string{}, true)
			fw.EXPECT().TLSConfig().Return(&tls.Config{}, nil)

			srv, err = NewServer(fw, bi, false)
			Expect(err).ToNot(HaveOccurred())
			Expect(srv.opts.TLS).To(BeTrue())
		})

		Describe("Gateways", func() {
			BeforeEach(func() {
				fw = NewMockChoriaFramework(mockctl)

				fw.EXPECT().TLSConfig().Return(&tls.Config{}, nil)
				fw.EXPECT().NetworkBrokerPeers().Return(srvcache.NewServers(), nil)
				fw.EXPECT().Logger(gomock.Any()).Return(logger).AnyTimes()
			})

			It("Should require a name and remotes", func() {
				config, err := config.NewConfig("testdata/gateways/noremotes.cfg")
				Expect(err).ToNot(HaveOccurred())

				fw.EXPECT().Configuration().Return(config).AnyTimes()

				srv, err = NewServer(fw, bi, false)
				Expect(err).ToNot(HaveOccurred())

				err = srv.setupGateways()
				Expect(err).To(MatchError("network Gateways require at least one remote"))
			})

			It("Should support remote gateways", func() {
				config, err := config.NewConfig("testdata/gateways/remotes.cfg")
				config.DisableTLSVerify = false
				Expect(err).ToNot(HaveOccurred())

				fw.EXPECT().Configuration().Return(config).AnyTimes()

				srv, err = NewServer(fw, bi, false)
				Expect(err).ToNot(HaveOccurred())

				Expect(srv.config.Choria.NetworkGatewayRemotes).To(Equal([]string{"C1", "C2"}))

				Expect(srv.opts.Gateway.Name).To(Equal("CHORIA"))
				Expect(srv.opts.Gateway.RejectUnknown).To(BeTrue())

				remotes := srv.opts.Gateway.Gateways
				Expect(remotes).To(HaveLen(2))
				Expect(remotes[0].Name).To(Equal("C1"))
				Expect(remotes[0].URLs).To(HaveLen(2))
				Expect(remotes[0].URLs[0].String()).To(Equal("nats://c1-1.example.net:7222"))
				Expect(remotes[0].URLs[1].String()).To(Equal("nats://c1-2.example.net:7222"))
				Expect(remotes[1].Name).To(Equal("C2"))
				Expect(remotes[1].TLSConfig).To(BeNil())
				Expect(remotes[1].URLs).To(HaveLen(2))
				Expect(remotes[1].URLs[0].String()).To(Equal("nats://c2-1.example.net:7222"))
				Expect(remotes[1].URLs[1].String()).To(Equal("nats://c2-2.example.net:7222"))
			})

			It("Should handle missing custom TLS", func() {
				config, err := config.NewConfig("testdata/gateways/missingtls.cfg")
				config.DisableTLSVerify = false
				Expect(err).ToNot(HaveOccurred())

				fw.EXPECT().Configuration().Return(config).AnyTimes()

				srv, err = NewServer(fw, bi, false)
				Expect(err).ToNot(HaveOccurred())

				Expect(srv.opts.Gateway.Gateways).To(HaveLen(1))
				Expect(srv.opts.Gateway.Gateways[0].Name).To(Equal("C2"))
			})

			It("Should support custom TLS", func() {
				config, err := config.NewConfig("testdata/gateways/customtls.cfg")
				config.DisableTLSVerify = false
				Expect(err).ToNot(HaveOccurred())

				fw.EXPECT().Configuration().Return(config).AnyTimes()

				srv, err = NewServer(fw, bi, false)
				Expect(err).ToNot(HaveOccurred())

				Expect(srv.opts.Gateway.Gateways).To(HaveLen(2))
				Expect(srv.opts.Gateway.Gateways[1].TLSConfig).ToNot(BeNil())
			})
		})

		Describe("Leafnodes", func() {
			BeforeEach(func() {
				fw = NewMockChoriaFramework(mockctl)

				fw.EXPECT().TLSConfig().Return(&tls.Config{}, nil)
				fw.EXPECT().NetworkBrokerPeers().Return(srvcache.NewServers(), nil)
				fw.EXPECT().Logger(gomock.Any()).Return(logger).AnyTimes()
			})

			It("Should support basic listening only leafnodes mode", func() {
				config, err := config.NewConfig("testdata/leafnodes/listening.cfg")
				Expect(err).ToNot(HaveOccurred())

				fw.EXPECT().Configuration().Return(config).AnyTimes()

				srv, err = NewServer(fw, bi, false)
				Expect(err).ToNot(HaveOccurred())
				Expect(srv.opts.LeafNode.Port).To(Equal(6222))
				Expect(srv.opts.LeafNode.Remotes).To(HaveLen(0))
			})

			It("Should support connecting to leafnodes", func() {
				config, err := config.NewConfig("testdata/leafnodes/remotes.cfg")
				Expect(err).ToNot(HaveOccurred())
				Expect(config.Choria.NetworkLeafRemotes).To(Equal([]string{"ln1", "ln2"}))
				fw.EXPECT().Configuration().Return(config).AnyTimes()

				srv, err = NewServer(fw, bi, false)
				Expect(err).ToNot(HaveOccurred())
				Expect(srv.IsTLS()).To(BeTrue())
				Expect(srv.opts.LeafNode.Port).To(Equal(6222))
				Expect(srv.opts.LeafNode.Remotes).To(HaveLen(2))
				Expect(srv.opts.LeafNode.Remotes[0].URLs[0].String()).To(Equal("leafnode://ln1-1.example.net:6222"))
				Expect(srv.opts.LeafNode.Remotes[0].TLS).ToNot(BeFalse())
				Expect(srv.opts.LeafNode.Remotes[1].URLs[0].String()).To(Equal("leafnode://ln2.example.net:6222"))
			})

			It("Should handle missing custom TLS", func() {
				config, err := config.NewConfig("testdata/leafnodes/missingtls.cfg")
				Expect(err).ToNot(HaveOccurred())
				Expect(config.Choria.NetworkLeafRemotes).To(Equal([]string{"ln1", "ln2"}))
				fw.EXPECT().Configuration().Return(config).AnyTimes()

				srv, err = NewServer(fw, bi, false)
				Expect(err).ToNot(HaveOccurred())
				Expect(srv.opts.LeafNode.Port).To(Equal(6222))
				Expect(srv.opts.LeafNode.Remotes).To(HaveLen(1))
				Expect(srv.opts.LeafNode.Remotes[0].URLs[0].String()).To(Equal("leafnode://ln2.example.net:6222"))
			})

			It("Should handle custom TLS", func() {
				config, err := config.NewConfig("testdata/leafnodes/customtls.cfg")
				Expect(err).ToNot(HaveOccurred())
				fw.EXPECT().Configuration().Return(config).AnyTimes()

				srv, err = NewServer(fw, bi, false)
				Expect(err).ToNot(HaveOccurred())
				Expect(srv.opts.LeafNode.Port).To(Equal(6222))
				Expect(srv.opts.LeafNode.Remotes).To(HaveLen(1))
				Expect(srv.opts.LeafNode.Remotes[0].URLs[0].String()).To(Equal("leafnode://ln1.example.net:6222"))
				Expect(srv.opts.LeafNode.Remotes[0].TLS).To(BeTrue())
				Expect(srv.opts.LeafNode.Remotes[0].TLSConfig).ToNot(BeNil())
			})
		})

		Describe("Accounts", func() {
			It("Should support JWT accounts", func() {
				cfg.Choria.NetworkAccountOperator = "choria_operator"
				cfg.ConfigFile = "testdata/broker.cfg"
				cfg.DisableTLS = true
				fw.EXPECT().NetworkBrokerPeers().Return(srvcache.NewServers(), nil)

				srv, err = NewServer(fw, bi, false)
				Expect(err).ToNot(HaveOccurred())
				Expect(srv.as).ToNot(BeNil())
			})

			It("Should fail when starting JWT accounts fails", func() {
				cfg.Choria.NetworkAccountOperator = "choria_operator"
				cfg.ConfigFile = "testdata/nonexisting/broker.cfg"
				cfg.DisableTLS = true

				srv, err = NewServer(fw, bi, false)
				Expect(err).To(HaveOccurred())
			})

			It("Should support setting system accounts", func() {
				cfg.Choria.NetworkAccountOperator = "choria_operator"
				cfg.Choria.NetworkSystemAccount = "ADMB22B4NQU27GI3KP6XUEFM5RSMOJY4O75NCP2P5JPQC2NGQNG6NJX2"
				cfg.ConfigFile = "testdata/broker.cfg"
				cfg.DisableTLS = true

				fw.EXPECT().NetworkBrokerPeers().Return(srvcache.NewServers(), nil)

				srv, err = NewServer(fw, bi, false)
				Expect(err).ToNot(HaveOccurred())
				Expect(srv.as).ToNot(BeNil())

				Expect(srv.opts.SystemAccount).To(Equal("ADMB22B4NQU27GI3KP6XUEFM5RSMOJY4O75NCP2P5JPQC2NGQNG6NJX2"))
			})

			It("Should support jetstream", func() {
				td, err := ioutil.TempDir("", "")
				Expect(err).ToNot(HaveOccurred())
				defer os.RemoveAll(td)

				fw.EXPECT().NetworkBrokerPeers().Return(srvcache.NewServers(), nil)
				cfg.DisableTLS = true
				cfg.Choria.NetworkStreamStore = td
				srv, err = NewServer(fw, bi, false)

				Expect(err).ToNot(HaveOccurred())
				Expect(srv.opts.JetStream).To(BeTrue())
				Expect(srv.opts.StoreDir).To(Equal(td))
			})

		})
	})
})
