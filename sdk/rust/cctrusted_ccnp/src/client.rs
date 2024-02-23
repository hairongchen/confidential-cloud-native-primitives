use anyhow::anyhow;
use crate::client::ccnp_server_pb::{
    ccnp_client::CcnpClient, GetCcEventlogRequest, GetCcEventlogResponse, GetCcMeasurementRequest,
    GetCcMeasurementResponse, GetCcReportRequest, GetCcReportResponse, GetDefaultAlgorithmRequest,
    GetDefaultAlgorithmResponse, GetMeasurementCountRequest, GetMeasurementCountResponse
};
use cctrusted_base::api_data::ExtraArgs;
use cctrusted_base::cc_type::TeeType;
use core::result::Result::Ok;
use hashbrown::HashMap;
use std::fs::File;
use std::io::Read;
use std::io::BufReader;
use std::io::read_to_string;
use tokio::net::UnixStream;
use tonic::transport::{Endpoint, Uri};
use tonic::Request;
use tower::service_fn;

//FixMe: use map from cc_type
lazy_static! {
    pub static ref TEE_VALUE_TYPE_MAP: HashMap<i32, TeeType> = {
        let mut map: HashMap<i32, TeeType> = HashMap::new();
        map.insert(-1, TeeType::PLAIN);
        map.insert(0, TeeType::TPM);
        map.insert(1, TeeType::TDX);
        map.insert(2, TeeType::SEV);
        map.insert(3, TeeType::CCA);
        map
    };
}

pub mod ccnp_server_pb {
    tonic::include_proto!("ccnp_server_pb");

    pub(crate) const FILE_DESCRIPTOR_SET: &[u8] =
        tonic::include_file_descriptor_set!("ccnp_server_descriptor");
}

pub struct CcnpServiceClient {
    pub ccnp_uds_path: String,
}

impl CcnpServiceClient {
    async fn get_cc_report_from_server_async(
        &mut self,
        nonce: Option<String>,
        data: Option<String>,
        _extra_args: ExtraArgs,
    ) -> Result<GetCcReportResponse, anyhow::Error> {
        let uds_path = self.ccnp_uds_path.parse::<Uri>().unwrap();
        let channel = Endpoint::try_from("http://[::]:0")
            .unwrap()
            .connect_with_connector(service_fn(move |_: Uri| {
                UnixStream::connect(uds_path.to_string())
            }))
            .await
            .unwrap();

        let container_id = match self.get_container_id() {
            Ok(id) => id,
            Err(e) => {
                return Err(anyhow!("[get_cc_report_from_server_async] error getting the container ID: {:?}", e));
            },
        };

        let request = Request::new(GetCcReportRequest {
            container_id,
            nonce,
            user_data: data,
        });

        let mut ccnp_client = CcnpClient::new(channel);

        let response = ccnp_client
            .get_cc_report(request)
            .await
            .unwrap()
            .into_inner();
        Ok(response)
    }

    // turn async call to sync call
    pub fn get_cc_report_from_server(
        &mut self,
        nonce: Option<String>,
        data: Option<String>,
        extra_args: ExtraArgs,
    ) -> Result<GetCcReportResponse, anyhow::Error> {
        let response = tokio::runtime::Builder::new_multi_thread()
            .enable_all()
            .build()
            .unwrap()
            .block_on(self.get_cc_report_from_server_async(nonce, data, extra_args));
        response
    }

    pub fn get_tee_type_by_value(&self, tee_id: &i32) -> TeeType {
        match TEE_VALUE_TYPE_MAP.get(tee_id) {
            Some(tee_type) => tee_type.clone(),
            None => TeeType::PLAIN,
        }
    }

    async fn get_cc_measurement_from_server_async(
        &mut self,
        index: u8,
        algo_id: u16,
    ) -> Result<GetCcMeasurementResponse, anyhow::Error> {
        let uds_path = self.ccnp_uds_path.parse::<Uri>().unwrap();
        let channel = Endpoint::try_from("http://[::]:0")
            .unwrap()
            .connect_with_connector(service_fn(move |_: Uri| {
                UnixStream::connect(uds_path.to_string())
            }))
            .await
            .unwrap();

        let container_id = match self.get_container_id() {
            Ok(id) => id,
            Err(e) => {
                return Err(anyhow!("[get_cc_measurement_from_server_async] error getting the container ID: {:?}", e));
            },
        };

        let request = Request::new(GetCcMeasurementRequest {
            container_id,
            index: index.into(),
            algo_id: algo_id.into(),
        });

        let mut ccnp_client = CcnpClient::new(channel);

        let response = ccnp_client
            .get_cc_measurement(request)
            .await
            .unwrap()
            .into_inner();
        Ok(response)
    }

    // turn async call to sync call
    pub fn get_cc_measurement_from_server(
        &mut self,
        index: u8,
        algo_id: u16,
    ) -> Result<GetCcMeasurementResponse, anyhow::Error> {
        let response = tokio::runtime::Builder::new_multi_thread()
            .enable_all()
            .build()
            .unwrap()
            .block_on(self.get_cc_measurement_from_server_async(index, algo_id));
        response
    }

    async fn get_cc_eventlog_from_server_async(
        &mut self,
        start: Option<u32>,
        count: Option<u32>,
    ) -> Result<GetCcEventlogResponse, anyhow::Error> {
        let uds_path = self.ccnp_uds_path.parse::<Uri>().unwrap();
        let channel = Endpoint::try_from("http://[::]:0")
            .unwrap()
            .connect_with_connector(service_fn(move |_: Uri| {
                UnixStream::connect(uds_path.to_string())
            }))
            .await
            .unwrap();

        let container_id = match self.get_container_id() {
            Ok(id) => id,
            Err(e) => {
                return Err(anyhow!("[get_cc_eventlog_from_server_async] error getting the container ID: {:?}", e));
            },
        };

        let request = Request::new(GetCcEventlogRequest {
            container_id,
            start,
            count,
        });

        let mut ccnp_client = CcnpClient::new(channel);

        let response = ccnp_client
            .get_cc_eventlog(request)
            .await
            .unwrap()
            .into_inner();
        Ok(response)
    }

    // turn async call to sync call
    pub fn get_cc_eventlog_from_server(
        &mut self,
        start: Option<u32>,
        count: Option<u32>,
    ) -> Result<GetCcEventlogResponse, anyhow::Error> {
        let response = tokio::runtime::Builder::new_multi_thread()
            .enable_all()
            .build()
            .unwrap()
            .block_on(self.get_cc_eventlog_from_server_async(start, count));
        response
    }

    async fn get_cc_measurement_count_from_server_async(
        &mut self,
    ) -> Result<GetMeasurementCountResponse, anyhow::Error> {
        let uds_path = self.ccnp_uds_path.parse::<Uri>().unwrap();
        let channel = Endpoint::try_from("http://[::]:0")
            .unwrap()
            .connect_with_connector(service_fn(move |_: Uri| {
                UnixStream::connect(uds_path.to_string())
            }))
            .await
            .unwrap();

        let request = Request::new(GetMeasurementCountRequest {});

        let mut ccnp_client = CcnpClient::new(channel);

        let response = ccnp_client
            .get_measurement_count(request)
            .await
            .unwrap()
            .into_inner();
        Ok(response)
    }

    // turn async call to sync call
    pub fn get_cc_measurement_count_from_server(
        &mut self,
    ) -> Result<GetMeasurementCountResponse, anyhow::Error> {
        let response = tokio::runtime::Builder::new_multi_thread()
            .enable_all()
            .build()
            .unwrap()
            .block_on(self.get_cc_measurement_count_from_server_async());
        response
    }

    async fn get_cc_default_algorithm_from_server_async(
        &mut self,
    ) -> Result<GetDefaultAlgorithmResponse, anyhow::Error> {
        let uds_path = self.ccnp_uds_path.parse::<Uri>().unwrap();
        let channel = Endpoint::try_from("http://[::]:0")
            .unwrap()
            .connect_with_connector(service_fn(move |_: Uri| {
                UnixStream::connect(uds_path.to_string())
            }))
            .await
            .unwrap();

        let request = Request::new(GetDefaultAlgorithmRequest {});

        let mut ccnp_client = CcnpClient::new(channel);

        let response = ccnp_client
            .get_default_algorithm(request)
            .await
            .unwrap()
            .into_inner();
        Ok(response)
    }

    // turn async call to sync call
    pub fn get_cc_default_algorithm_from_server(
        &mut self,
    ) -> Result<GetDefaultAlgorithmResponse, anyhow::Error> {
        let response = tokio::runtime::Builder::new_multi_thread()
            .enable_all()
            .build()
            .unwrap()
            .block_on(self.get_cc_default_algorithm_from_server_async());
        response
    }

    pub fn get_container_id(&self) -> Result<String, anyhow::Error> {
        let mountinfo = "/proc/self/mountinfo";
        let docker_pattern = "/docker/containers/";
        let k8s_pattern = "/kubelet/pods/";

        let mut data_lines = Vec::new();

        // let file = File::open(mountinfo)?;
        // let mut file_reader = BufReader::new(file);
        // let mut file_string = String::new();
        // let _ = file_reader.read_to_string(&mut file_string);
        data_lines = read_to_string(mountinfo)
            .unwrap()
            .lines()
            .map(String::from);

        for line in data_lines {
            /*
             * /var/lib/docker/containers/{container-id}/{file}
             * example: 
             */
            if line.contains(docker_pattern) {
                if let e = line.split(docker_pattern).last() {
                    Ok(e.split("/").first())
                } else {
                    Err(anyhow!("[get_container_id] incorrect docker container info in /proc/self/mountinfo!"))
                }
            }

            /*
             * line: /var/lib/kubelet/pods/{container-id}/{file}
             * example: 2958 2938 253:1 /var/lib/kubelet/pods/a45f46f0-20be-45ab-ace6-b77e8e2f062c/containers/busybox/8f8d892c /dev/termination-log rw,relatime - ext4 /dev/vda1 rw,discard,errors=remount-ro
             */
            if line.contains(k8s_pattern){
                if let e = line.split(k8s_pattern).last() {
                    Ok(e.split("/").first().replace("-", "-"))
                } else {
                    Err(anyhow!("[get_container_id] incorrect k8s pod container info in /proc/self/mountinfo!"))
                }
            }
        }

        Err(anyhow!("[get_container_id] no container info in /proc/self/mountinfo!"))

    }
}
