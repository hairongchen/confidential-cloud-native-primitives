/*
* Copyright (c) 2024, Intel Corporation. All rights reserved.<BR>
* SPDX-License-Identifier: Apache-2.0
 */

package ccnp

import (
	"errors"

	"github.com/cc-api/cc-trusted-api/common/golang/cctrusted_base"
	"github.com/cc-api/cc-trusted-api/common/golang/cctrusted_base/tdx"
)

var _ cctrusted_base.CCTrustedAPI = (*SDK)(nil)

type SDK struct {
}

// GetCCReport implements CCTrustedAPI
func (s *SDK) GetCCReport(nonce string, userData string, _ any) (cctrusted_base.Report, error) {
	result, err := GetCCReportFromServer(userData, nonce)
	if err != nil {
		return nil, err
	}

	switch cctrusted_base.TYPE_CC_TDX { //FIXME: use type get from result
	case cctrusted_base.TYPE_CC_TDX:
		report, err := tdx.NewTdxReportFromBytes(result.CcReport)
		if err != nil {
			return nil, err
		}
		return report, nil
	default:
	}
	return nil, errors.New("[GetCCReport] get CC report failed")
}

// DumpCCReport implements cctrusted_base.CCTrustedAPI.
func (s *SDK) DumpCCReport(reportBytes []byte) error {
	return nil
}

// GetCCMeasurement implements cctrusted_base.CCTrustedAPI.
func (s *SDK) GetCCMeasurement(index int, alg cctrusted_base.TCG_ALG) (cctrusted_base.TcgDigest, error) {
	panic("not implemented!")
}

// GetMeasurementCount implements cctrusted_base.CCTrustedAPI.
func (s *SDK) GetMeasurementCount() (int, error) {
	panic("not implemented!")
}

// ReplayCCEventLog implements cctrusted_base.CCTrustedAPI.
func (s *SDK) ReplayCCEventLog(formatedEventLogs []cctrusted_base.FormatedTcgEvent) map[int]map[cctrusted_base.TCG_ALG][]byte {
	panic("not implemented!")
}

// GetDefaultAlgorithm implements cctrusted_base.CCTrustedAPI.
func (s *SDK) GetDefaultAlgorithm() cctrusted_base.TCG_ALG {
	panic("not implemented!")
}

// SelectEventlog implements CCTrustedAPI.
func (s *SDK) GetCCEventLog(start int32, count int32) (*cctrusted_base.EventLogger, error) {
	panic("not implemented!")
}
