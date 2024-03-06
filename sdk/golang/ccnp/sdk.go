/*
* Copyright (c) 2024, Intel Corporation. All rights reserved.<BR>
* SPDX-License-Identifier: Apache-2.0
 */

package ccnpsdk

import (
	"errors"

	"github.com/cc-api/cc-trusted-api/common/golang/cctrusted_base"
	"github.com/cc-api/cc-trusted-api/common/golang/cctrusted_base/tdx"
	pb "github.com/hairongchen/confidential-cloud-native-primitives/sdk/golang/ccnp/proto"
)

var _ cctrusted_base.CCTrustedAPI = (*SDK)(nil)

type SDK struct {
}

// GetCCReport implements CCTrustedAPI
func (s *SDK) GetCCReport(nonce string, userData string, _ any) (pb.GetCcReportResponse, error) {
	response, err := GetCCReportFromServer(userData, nonce)
	if err != nil {
		return nil, err
	}

	switch response.cc_type {
	case cctrusted_base.TYPE_CC_TDX:
		report, err := tdx.NewTdxReportFromBytes(reportBytes)
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
	emptyRet := cctrusted_base.TcgDigest{}
	return emptyRet, nil
}

// GetMeasurementCount implements cctrusted_base.CCTrustedAPI.
func (s *SDK) GetMeasurementCount() (int, error) {
	return 4, nil
}

// ReplayCCEventLog implements cctrusted_base.CCTrustedAPI.
// func (s *SDK) ReplayCCEventLog(formatedEventLogs []cctrusted_base.FormatedTcgEvent) map[int]map[cctrusted_base.TCG_ALG][]byte {
func (s *SDK) ReplayCCEventLog(formatedEventLogs []cctrusted_base.FormatedTcgEvent) error {
	return nil
}

// GetDefaultAlgorithm implements cctrusted_base.CCTrustedAPI.
// func (s *SDK) GetDefaultAlgorithm() cctrusted_base.TCG_ALG {
func (s *SDK) GetDefaultAlgorithm() error {
	return nil
}

// SelectEventlog implements CCTrustedAPI.
// func (s *SDK) GetCCEventLog(start int32, count int32) (*cctrusted_base.EventLogger, error) {
func (s *SDK) GetCCEventLog(start int32, count int32) error {

	el, err := s.internelEventlog()

	return nil
}
