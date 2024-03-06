/*
* Copyright (c) 2024, Intel Corporation. All rights reserved.<BR>
* SPDX-License-Identifier: Apache-2.0
 */

package ccnpsdk

import (
	pb "github.com/hairongchen/confidential-cloud-native-primitives/sdk/golang/ccnp/proto"
)

const (
	UDS_PATH = "unix:/run/ccnp/uds/ccnp-server.sock"
)

func GetCCReportFromServer(userData string, nonce string) (GetCcReportResponse, error) {
	channel, err := grpc.Dial(UDS_PATH, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("[GetCCReportFromServer] can not connect to CCNP server UDS at %v with error: %v", UDS_PATH, err)
	}

	defer channel.Close()

	client := pb.NewCcnpClient(channel)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var ContainerId = GetContainerId();

	response, err := client.GetCcReport(ctx, &pb.GetCcReportRequest{container_id: containerId, nonce: nonce, user_data: userData})
	if err != nil {
		return nil, fmt.Errorf("[GetCCReportFromServer] fail to get cc report with error: %v", err)
	}

	return response, nil
}

func GetContainerId() (id string, error){
	var mountinfoFile string = "/proc/self/mountinfo"
	var dockerPattern string = "/docker/containers/"
	var k8sPattern string = "/kubelet/pods/"

	file, err := os.Open(mountinfoFile)
	if err != nil {
		return nil, fmt.Errorf("[GetContainerId] fail to open mountinfo file: %v", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	for _, line := range lines {
        /*
         * line format:
         *      ... /var/lib/docker/containers/{container-id}/{file} ...
         * sample:
         */
		if strings.Contains(line, dockerPattern) {
			// /var/lib/docker/containers/{container-id}/{file}
			var res = strings.Split(line, dockerPattern)
			var res1 = res[len(res)-1]
			var containerId = strings.Split(res1, '/')[0]

			return containerId, nil
		}

        /*
         * line format:
         *      ... /var/lib/kubelet/pods/{container-id}/{file} ...
         * sample:
         *      2958 2938 253:1 /var/lib/kubelet/pods/a45f46f0-20be-45ab-ace6-b77e8e2f062c/containers/busybox/8f8d892c /dev/termination-log rw,relatime - ext4 /dev/vda1 rw,discard,errors=remount-ro
         */
		if strings.Contains(line, k8sPattern) {
			// /var/lib/kubelet/pods/{container-id}/{file}
			var res = strings.Split(line, k8sPattern)
			var res1 = res[len(res)-1]
			var res2 = strings.Split(res1, '/')[0]
			var ContainerId = strings.Replace(res2, "-", "_", -1)

			return containerId, nil
		}
	}

	return nil, errors.New("[GetContainerId] no docker or kubernetes container patter found in /proc/self/mountinfo")
}
