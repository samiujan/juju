// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package agent_test
import (
	gc "launchpad.net/gocheck"

	"launchpad.net/juju-core/agent"
	"launchpad.net/juju-core/environs/config"
	"launchpad.net/juju-core/state"
	"launchpad.net/juju-core/instance"
	"launchpad.net/juju-core/constraints"
	"launchpad.net/juju-core/testing"
	"launchpad.net/juju-core/version"
)

type bootstrapSuite struct {
	testing.MgoSuite
}
var _ = gc.Suite(&bootstrapSuite{})

func (s *bootstrapSuite) TestInitializeState(c *gc.C) {
	dataDir := c.MkDir()
	
	cfg, err := agent.NewAgentConfig(agent.AgentConfigParams{
		DataDir: dataDir,
		Tag: "machine-0",
		Nonce: "onceonly",
		StateAddresses: []string{testing.MgoAddr},
		CACert: []byte(testing.CACert),
		Password: testing.DefaultMongoPassword,
	})
	c.Assert(err, gc.IsNil)
	expectConstraints := constraints.MustParse("mem=1024M")
	expectHW := instance.MustParseHardware("mem=2048M")
	mcfg := agent.BootstrapMachineConfig{
		Constraints: expectConstraints,
		Jobs: []state.MachineJob{state.JobHostUnits},
		InstanceId: "i-bootstrap",
		Characteristics: expectHW,
	}
	attrs := testing.FakeConfig().Delete("admin-secret").Merge(testing.Attrs{
		"agent-version": version.Current.Number.String(),
	})
	envcfg, err := config.New(config.NoDefaults, attrs)
	c.Assert(err, gc.IsNil)

	st, m, err := cfg.InitializeState(envcfg, mcfg, state.DialOpts{})
	c.Assert(err, gc.IsNil)
	defer st.Close()
	c.Assert(m.Id(), gc.Equals, "0")
	c.Assert(m.Jobs(), gc.DeepEquals, []state.MachineJob{state.JobHostUnits})
	gotConstraints, err := m.Constraints()
	c.Assert(err, gc.IsNil)
	c.Assert(gotConstraints, gc.DeepEquals, expectConstraints)
	c.Assert(err, gc.IsNil)
	gotHW, err := m.HardwareCharacteristics()
	c.Assert(err, gc.IsNil)
	c.Assert(*gotHW, gc.DeepEquals, expectHW)
}

//test we can log in as admin
//
//read config, test that we can use it to log in to state
//test that 
