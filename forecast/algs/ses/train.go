// Copyright 2018-20 PJ Engineering and Business Solutions Pty. Ltd. All rights reserved.

package ses

import (
	"context"
	"math"

	"github.com/rocketlaunchr/dataframe-go/forecast"
)

type trainingState struct {
	initialLevel   float64
	originValue    float64
	smoothingLevel float64
	rmse           float64
	zValues        map[float64]float64
}

func (se *SimpleExpSmoothing) trainSeries(ctx context.Context, start, end int) error {

	var (
		α       float64 = se.cfg.Alpha
		st      float64
		Yorigin float64
		mse     float64 // mean squared error
	)

	count := 0
	// Training smoothing Level
	for i := start; i < end+1; i++ {

		if err := ctx.Err(); err != nil {
			return err
		}

		xt := se.sf.Values[i]

		if i == start {
			st = xt
			se.tstate.initialLevel = xt
		} else if i == end { // Setting the last value in traindata as Yorigin value for bootstrapping
			Yorigin = xt
			se.tstate.originValue = Yorigin

		} else {
			st = α*xt + (1-α)*st

			mse += (xt - st) * (xt - st)
			count++
		}

	}
	mse /= float64(count)

	// calculate ZValues from confidence levels
	zVals := make(map[float64]float64, len(se.cfg.ConfidenceLevels))
	for _, l := range se.cfg.ConfidenceLevels {
		zVals[l] = forecast.ConfidenceLevelToZ(l)
	}

	se.tstate.rmse = math.Sqrt(mse)
	se.tstate.zValues = zVals
	se.tstate.smoothingLevel = st

	return nil
}
