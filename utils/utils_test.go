package utils_test

import (
	"errors"
	"fmt"
	"os"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/rexray/gocsi/utils"
)

var errMissingCSIEndpoint = errors.New("missing CSI_ENDPOINT")

var _ = Describe("ParseVersion", func() {
	shouldParse := func() csi.Version {
		v, ok := utils.ParseVersion(
			CurrentGinkgoTestDescription().ComponentTexts[1])
		Ω(ok).Should(BeTrue())
		return v
	}
	Context("0.0.0", func() {
		It("Should Parse", func() {
			v := shouldParse()
			Ω(v.GetMajor()).Should(Equal(uint32(0)))
			Ω(v.GetMinor()).Should(Equal(uint32(0)))
			Ω(v.GetPatch()).Should(Equal(uint32(0)))
		})
	})
	Context("0.1.0", func() {
		It("Should Parse", func() {
			v := shouldParse()
			Ω(v.GetMajor()).Should(Equal(uint32(0)))
			Ω(v.GetMinor()).Should(Equal(uint32(1)))
			Ω(v.GetPatch()).Should(Equal(uint32(0)))
		})
	})
	Context("1.1.0", func() {
		It("Should Parse", func() {
			v := shouldParse()
			Ω(v.GetMajor()).Should(Equal(uint32(1)))
			Ω(v.GetMinor()).Should(Equal(uint32(1)))
			Ω(v.GetPatch()).Should(Equal(uint32(0)))
		})
	})
})

var _ = Describe("ParseVersions", func() {
	shouldParse := func() []csi.Version {
		return utils.ParseVersions(
			CurrentGinkgoTestDescription().ComponentTexts[1])
	}
	Context("0.1.0 0.2.0 1.0.0 1.1.0", func() {
		It("Should Parse", func() {
			v := shouldParse()
			Ω(v).Should(HaveLen(4))
			Ω(v[0].Major).Should(Equal(uint32(0)))
			Ω(v[0].Minor).Should(Equal(uint32(1)))
			Ω(v[0].Patch).Should(Equal(uint32(0)))

			Ω(v[1].Major).Should(Equal(uint32(0)))
			Ω(v[1].Minor).Should(Equal(uint32(2)))
			Ω(v[1].Patch).Should(Equal(uint32(0)))

			Ω(v[2].Major).Should(Equal(uint32(1)))
			Ω(v[2].Minor).Should(Equal(uint32(0)))
			Ω(v[2].Patch).Should(Equal(uint32(0)))

			Ω(v[3].Major).Should(Equal(uint32(1)))
			Ω(v[3].Minor).Should(Equal(uint32(1)))
			Ω(v[3].Patch).Should(Equal(uint32(0)))
		})
	})
})

var _ = Describe("GetCSIEndpoint", func() {
	var (
		err         error
		proto       string
		addr        string
		expEndpoint string
		expProto    string
		expAddr     string
	)
	BeforeEach(func() {
		expEndpoint = CurrentGinkgoTestDescription().ComponentTexts[2]
		os.Setenv(utils.CSIEndpoint, expEndpoint)
	})
	AfterEach(func() {
		proto = ""
		addr = ""
		expEndpoint = ""
		expProto = ""
		expAddr = ""
		os.Unsetenv(utils.CSIEndpoint)
	})
	JustBeforeEach(func() {
		proto, addr, err = utils.GetCSIEndpoint()
	})

	Context("Valid Endpoint", func() {
		shouldBeValid := func() {
			Ω(os.Getenv(utils.CSIEndpoint)).Should(Equal(expEndpoint))
			Ω(proto).Should(Equal(expProto))
			Ω(addr).Should(Equal(expAddr))
		}
		Context("tcp://127.0.0.1", func() {
			BeforeEach(func() {
				expProto = "tcp"
				expAddr = "127.0.0.1"
			})
			It("Should Be Valid", shouldBeValid)
		})
		Context("tcp://127.0.0.1:8080", func() {
			BeforeEach(func() {
				expProto = "tcp"
				expAddr = "127.0.0.1:8080"
			})
			It("Should Be Valid", shouldBeValid)
		})
		Context("tcp://*:8080", func() {
			BeforeEach(func() {
				expProto = "tcp"
				expAddr = "*:8080"
			})
			It("Should Be Valid", shouldBeValid)
		})
		Context("unix://path/to/sock.sock", func() {
			BeforeEach(func() {
				expProto = "unix"
				expAddr = "path/to/sock.sock"
			})
			It("Should Be Valid", shouldBeValid)
		})
		Context("unix:///path/to/sock.sock", func() {
			BeforeEach(func() {
				expProto = "unix"
				expAddr = "/path/to/sock.sock"
			})
			It("Should Be Valid", shouldBeValid)
		})
		Context("sock.sock", func() {
			BeforeEach(func() {
				expProto = "unix"
				expAddr = "sock.sock"
			})
			It("Should Be Valid", shouldBeValid)
		})
		Context("/tmp/sock.sock", func() {
			BeforeEach(func() {
				expProto = "unix"
				expAddr = "/tmp/sock.sock"
			})
			It("Should Be Valid", shouldBeValid)
		})
	})

	Context("Missing Endpoint", func() {
		Context("", func() {
			It("Should Be Missing", func() {
				Ω(err).Should(HaveOccurred())
				Ω(err).Should(Equal(errMissingCSIEndpoint))
			})
		})
		Context("    ", func() {
			It("Should Be Missing", func() {
				Ω(err).Should(HaveOccurred())
				Ω(err).Should(Equal(errMissingCSIEndpoint))
			})
		})
	})

	Context("Invalid Network Address", func() {
		shouldBeInvalid := func() {
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(Equal(fmt.Sprintf(
				"invalid network address: %s", expEndpoint)))
		}
		Context("tcp5://localhost:5000", func() {
			It("Should Be An Invalid Endpoint", shouldBeInvalid)
		})
		Context("unixpcket://path/to/sock.sock", func() {
			It("Should Be An Invalid Endpoint", shouldBeInvalid)
		})
	})

	Context("Invalid Implied Sock File", func() {
		shouldBeInvalid := func() {
			Ω(err).Should(HaveOccurred())
			Ω(err.Error()).Should(Equal(fmt.Sprintf(
				"invalid implied sock file: %[1]s: "+
					"open %[1]s: no such file or directory",
				expEndpoint)))
		}
		Context("Xtcp5://localhost:5000", func() {
			It("Should Be An Invalid Implied Sock File", shouldBeInvalid)
		})
		Context("Xunixpcket://path/to/sock.sock", func() {
			It("Should Be An Invalid Implied Sock File", shouldBeInvalid)
		})
	})
})

var _ = Describe("ParseProtoAddr", func() {
	Context("Empty Address", func() {
		It("Should Be An Empty Address", func() {
			_, _, err := utils.ParseProtoAddr("")
			Ω(err).Should(HaveOccurred())
			Ω(err).Should(Equal(utils.ErrParseProtoAddrRequired))
		})
		It("Should Be An Empty Address", func() {
			_, _, err := utils.ParseProtoAddr("   ")
			Ω(err).Should(HaveOccurred())
			Ω(err).Should(Equal(utils.ErrParseProtoAddrRequired))
		})
	})
})

var _ = Describe("ParseMap", func() {
	Context("One Pair", func() {
		It("Should Be Valid", func() {
			data := utils.ParseMap("k1=v1")
			Ω(data).Should(HaveLen(1))
			Ω(data["k1"]).Should(Equal("v1"))
		})
	})
	Context("Empty Line", func() {
		It("Should Be Valid", func() {
			data := utils.ParseMap("")
			Ω(data).Should(HaveLen(0))
		})
	})
	Context("Key Sans Value", func() {
		It("Should Be Valid", func() {
			data := utils.ParseMap("k1")
			Ω(data).Should(HaveLen(1))
		})
	})
	Context("Two Pair", func() {
		It("Should Be Valid", func() {
			data := utils.ParseMap("k1=v1, k2=v2")
			Ω(data).Should(HaveLen(2))
			Ω(data["k1"]).Should(Equal("v1"))
			Ω(data["k2"]).Should(Equal("v2"))
		})
	})
	Context("Two Pair with Quoting & Escaping", func() {
		It("Should Be Valid", func() {
			data := utils.ParseMap(`k1=v1, "k2=v2""s"`)
			Ω(data).Should(HaveLen(2))
			Ω(data["k1"]).Should(Equal("v1"))
			Ω(data["k2"]).Should(Equal(`v2"s`))
		})
		It("Should Be Valid", func() {
			data := utils.ParseMap(`k1=v1, "k2=v2\'s"`)
			Ω(data).Should(HaveLen(2))
			Ω(data["k1"]).Should(Equal("v1"))
			Ω(data["k2"]).Should(Equal(`v2\'s`))
		})
		It("Should Be Valid", func() {
			data := utils.ParseMap(`k1=v1, k2=v2's`)
			Ω(data).Should(HaveLen(2))
			Ω(data["k1"]).Should(Equal("v1"))
			Ω(data["k2"]).Should(Equal(`v2's`))
		})
	})
	Context("Two Pair with Three Spaces Between Them", func() {
		It("Should Be Valid", func() {
			data := utils.ParseMap("k1=v1,   k2=v2")
			Ω(data).Should(HaveLen(2))
			Ω(data["k1"]).Should(Equal("v1"))
			Ω(data["k2"]).Should(Equal("v2"))
		})
	})
	Context("Two Pair with One Sans Value", func() {
		It("Should Be Valid", func() {
			data := utils.ParseMap("k1=, k2=v2")
			Ω(data).Should(HaveLen(2))
			Ω(data["k1"]).Should(Equal(""))
			Ω(data["k2"]).Should(Equal("v2"))
		})
	})
	Context("Two Pair with One Sans Value & Three Spaces Between Them", func() {
		It("Should Be Valid", func() {
			data := utils.ParseMap("k1=,    k2=v2")
			Ω(data).Should(HaveLen(2))
			Ω(data["k1"]).Should(Equal(""))
			Ω(data["k2"]).Should(Equal("v2"))
		})
	})
	Context("One Pair with Quoted Value", func() {
		It("Should Be Valid", func() {
			data := utils.ParseMap("k1=v 1")
			Ω(data).Should(HaveLen(1))
			Ω(data["k1"]).Should(Equal("v 1"))
		})
	})
	Context("Three Pair with Mixed Values", func() {
		It("Should Be Valid", func() {
			data := utils.ParseMap(`"k1=v 1", "k2=v 2 ", "k3 =v3"  `)
			Ω(data).Should(HaveLen(3))
			Ω(data["k1"]).Should(Equal("v 1"))
			Ω(data["k2"]).Should(Equal("v 2 "))
			Ω(data["k3 "]).Should(Equal("v3"))
		})
	})
})

var _ = Describe("CompareVolumeInfo", func() {
	It("a == b", func() {
		a := csi.VolumeInfo{Id: "0"}
		b := csi.VolumeInfo{Id: "0"}
		Ω(utils.CompareVolumeInfo(a, b)).Should(Equal(0))
		a.CapacityBytes = 1
		b.CapacityBytes = 1
		Ω(utils.CompareVolumeInfo(a, b)).Should(Equal(0))
		a.Attributes = map[string]string{"key": "val"}
		b.Attributes = map[string]string{"key": "val"}
		Ω(utils.CompareVolumeInfo(a, b)).Should(Equal(0))
	})
	It("a > b", func() {
		a := csi.VolumeInfo{Id: "0"}
		b := csi.VolumeInfo{}
		Ω(utils.CompareVolumeInfo(a, b)).Should(Equal(1))
		b.Id = "0"
		Ω(utils.CompareVolumeInfo(a, b)).Should(Equal(0))
		a.CapacityBytes = 1
		Ω(utils.CompareVolumeInfo(a, b)).Should(Equal(1))
		b.CapacityBytes = 1
		Ω(utils.CompareVolumeInfo(a, b)).Should(Equal(0))
		a.Attributes = map[string]string{"key": "val"}
		Ω(utils.CompareVolumeInfo(a, b)).Should(Equal(1))
		b.Attributes = map[string]string{"key": "val"}
		Ω(utils.CompareVolumeInfo(a, b)).Should(Equal(0))
	})
	It("a < b", func() {
		b := csi.VolumeInfo{Id: "0"}
		a := csi.VolumeInfo{}
		Ω(utils.CompareVolumeInfo(a, b)).Should(Equal(-1))
		a.Id = "0"
		Ω(utils.CompareVolumeInfo(a, b)).Should(Equal(0))
		b.CapacityBytes = 1
		Ω(utils.CompareVolumeInfo(a, b)).Should(Equal(-1))
		a.CapacityBytes = 1
		Ω(utils.CompareVolumeInfo(a, b)).Should(Equal(0))
		b.Attributes = map[string]string{"key": "val"}
		Ω(utils.CompareVolumeInfo(a, b)).Should(Equal(-1))
		a.Attributes = map[string]string{"key": "val"}
		Ω(utils.CompareVolumeInfo(a, b)).Should(Equal(0))
	})
})

var _ = Describe("EqualVolumeCapability", func() {
	It("a == b", func() {
		a := &csi.VolumeCapability{
			AccessMode: &csi.VolumeCapability_AccessMode{
				Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
			},
			AccessType: &csi.VolumeCapability_Block{
				Block: &csi.VolumeCapability_BlockVolume{},
			},
		}
		b := &csi.VolumeCapability{
			AccessMode: &csi.VolumeCapability_AccessMode{
				Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
			},
			AccessType: &csi.VolumeCapability_Block{
				Block: &csi.VolumeCapability_BlockVolume{},
			},
		}
		Ω(utils.EqualVolumeCapability(a, b)).Should(BeTrue())
		a.AccessMode.Mode = csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY
		Ω(utils.EqualVolumeCapability(a, b)).Should(BeFalse())
		b.AccessMode.Mode = csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY
		Ω(utils.EqualVolumeCapability(a, b)).Should(BeTrue())
		a.AccessMode = nil
		Ω(utils.EqualVolumeCapability(a, b)).Should(BeFalse())
		b.AccessMode = nil
		Ω(utils.EqualVolumeCapability(a, b)).Should(BeTrue())
		a = nil
		Ω(utils.EqualVolumeCapability(nil, b)).Should(BeFalse())
		b = nil
		Ω(utils.EqualVolumeCapability(a, b)).Should(BeFalse())

		aAT := &csi.VolumeCapability_Mount{
			Mount: &csi.VolumeCapability_MountVolume{
				FsType:     "ext4",
				MountFlags: []string{"rw"},
			},
		}
		bAT := &csi.VolumeCapability_Mount{
			Mount: &csi.VolumeCapability_MountVolume{
				FsType:     "ext4",
				MountFlags: []string{"rw"},
			},
		}

		a = &csi.VolumeCapability{
			AccessMode: &csi.VolumeCapability_AccessMode{
				Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
			},
			AccessType: aAT,
		}
		b = &csi.VolumeCapability{
			AccessMode: &csi.VolumeCapability_AccessMode{
				Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
			},
			AccessType: bAT,
		}
		Ω(utils.EqualVolumeCapability(a, b)).Should(BeTrue())
		aAT.Mount.FsType = "xfs"
		Ω(utils.EqualVolumeCapability(a, b)).Should(BeFalse())
		bAT.Mount.FsType = "xfs"
		Ω(utils.EqualVolumeCapability(a, b)).Should(BeTrue())
		aAT.Mount.MountFlags = append(aAT.Mount.MountFlags, "nosuid")
		Ω(utils.EqualVolumeCapability(a, b)).Should(BeFalse())
		bAT.Mount.MountFlags = append(bAT.Mount.MountFlags, "nosuid")
		Ω(utils.EqualVolumeCapability(a, b)).Should(BeTrue())
		aAT.Mount.MountFlags[0] = "ro"
		Ω(utils.EqualVolumeCapability(a, b)).Should(BeFalse())
		bAT.Mount.MountFlags[0] = "ro"
		Ω(utils.EqualVolumeCapability(a, b)).Should(BeTrue())
		aAT.Mount.MountFlags = nil
		Ω(utils.EqualVolumeCapability(a, b)).Should(BeFalse())
		bAT.Mount.MountFlags = nil
		Ω(utils.EqualVolumeCapability(a, b)).Should(BeTrue())
	})
})

var _ = Describe("AreVolumeCapabilitiesCompatible", func() {
	It("compatible", func() {
		aMountAT := &csi.VolumeCapability_Mount{
			Mount: &csi.VolumeCapability_MountVolume{
				FsType:     "ext4",
				MountFlags: []string{"rw"},
			},
		}
		a := []*csi.VolumeCapability{
			&csi.VolumeCapability{
				AccessMode: &csi.VolumeCapability_AccessMode{
					Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
				},
				AccessType: aMountAT,
			},
		}

		b := []*csi.VolumeCapability{
			&csi.VolumeCapability{
				AccessMode: &csi.VolumeCapability_AccessMode{
					Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
				},
				AccessType: &csi.VolumeCapability_Block{
					Block: &csi.VolumeCapability_BlockVolume{},
				},
			},
			&csi.VolumeCapability{
				AccessMode: &csi.VolumeCapability_AccessMode{
					Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_WRITER,
				},
				AccessType: &csi.VolumeCapability_Mount{
					Mount: &csi.VolumeCapability_MountVolume{
						FsType:     "ext4",
						MountFlags: []string{"rw"},
					},
				},
			},
		}

		Ω(utils.AreVolumeCapabilitiesCompatible(a, b)).Should(BeTrue())
		aMountAT.Mount.MountFlags[0] = "ro"
		Ω(utils.AreVolumeCapabilitiesCompatible(a, b)).Should(BeFalse())
		a[0].AccessType = &csi.VolumeCapability_Block{
			Block: &csi.VolumeCapability_BlockVolume{},
		}
		Ω(utils.AreVolumeCapabilitiesCompatible(a, b)).Should(BeTrue())
	})
})
