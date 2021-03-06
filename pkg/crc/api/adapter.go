package api

import (
	"context"

	"github.com/code-ready/crc/pkg/crc/api/client"
	"github.com/code-ready/crc/pkg/crc/logging"
	"github.com/code-ready/crc/pkg/crc/machine"
	"github.com/code-ready/machine/libmachine/state"
)

type AdaptedClient interface {
	GetName() string
	Delete() client.Result
	GetConsoleURL() client.ConsoleResult
	Start(ctx context.Context, startConfig machine.StartConfig) client.StartResult
	Status() client.ClusterStatusResult
	Stop() client.Result
}

type Adapter struct {
	Underlying machine.Client
}

func (a *Adapter) GetName() string {
	return a.Underlying.GetName()
}

func (a *Adapter) Delete() client.Result {
	err := a.Underlying.Delete()
	if err != nil {
		logging.Error(err)
		return client.Result{
			Name:    a.Underlying.GetName(),
			Success: false,
			Error:   err.Error(),
		}
	}
	return client.Result{
		Name:    a.Underlying.GetName(),
		Success: true,
	}
}

func (a *Adapter) GetConsoleURL() client.ConsoleResult {
	res, err := a.Underlying.GetConsoleURL()
	if err != nil {
		logging.Error(err)
		return client.ConsoleResult{
			Success: false,
			Error:   err.Error(),
		}
	}
	return client.ConsoleResult{
		ClusterConfig: res.ClusterConfig,
		Success:       true,
	}
}

func (a *Adapter) Start(ctx context.Context, startConfig machine.StartConfig) client.StartResult {
	res, err := a.Underlying.Start(ctx, startConfig)
	if err != nil {
		logging.Error(err)
		return client.StartResult{
			Name:  a.Underlying.GetName(),
			Error: err.Error(),
		}
	}
	return client.StartResult{
		Name:           a.Underlying.GetName(),
		Status:         res.Status.String(),
		ClusterConfig:  res.ClusterConfig,
		KubeletStarted: res.KubeletStarted,
	}
}

func (a *Adapter) Status() client.ClusterStatusResult {
	res, err := a.Underlying.Status()
	if err != nil {
		logging.Error(err)
		return client.ClusterStatusResult{
			Name:    a.Underlying.GetName(),
			Error:   err.Error(),
			Success: false,
		}
	}
	return client.ClusterStatusResult{
		Name:             a.Underlying.GetName(),
		CrcStatus:        res.CrcStatus.String(),
		OpenshiftStatus:  res.OpenshiftStatus,
		OpenshiftVersion: res.OpenshiftVersion,
		DiskUse:          res.DiskUse,
		DiskSize:         res.DiskSize,
		Success:          true,
	}
}

func (a *Adapter) Stop() client.Result {
	vmState, err := a.Underlying.Stop()
	if err != nil {
		logging.Error(err)
		if vmState == state.Running {
			err := a.Underlying.PowerOff()
			if err != nil {
				logging.Error(err)
				return client.Result{
					Name:    a.Underlying.GetName(),
					Success: false,
					Error:   err.Error(),
				}
			}
		}
	}
	return client.Result{
		Name:    a.Underlying.GetName(),
		Success: true,
	}
}
