/*
* Copyright (c) 2024, Intel Corporation. All rights reserved.<BR>
* SPDX-License-Identifier: Apache-2.0
 */

package ccnpsdk

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
