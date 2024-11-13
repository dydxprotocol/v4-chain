// This file is @generated by prost-build.
/// EpochInfo stores metadata of an epoch timer.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct EpochInfo {
    /// name is the unique identifier.
    #[prost(string, tag = "1")]
    pub name: ::prost::alloc::string::String,
    /// next_tick indicates when the next epoch starts (in Unix Epoch seconds),
    /// if `EpochInfo` has been initialized.
    /// If `EpochInfo` is not initialized yet, `next_tick` indicates the earliest
    /// initialization time (see `is_initialized` below).
    #[prost(uint32, tag = "2")]
    pub next_tick: u32,
    /// duration of the epoch in seconds.
    #[prost(uint32, tag = "3")]
    pub duration: u32,
    /// current epoch is the number of the current epoch.
    /// 0 if `next_tick` has never been reached, positive otherwise.
    #[prost(uint32, tag = "4")]
    pub current_epoch: u32,
    /// current_epoch_start_block indicates the block height when the current
    /// epoch started. 0 if `current_epoch` is 0.
    #[prost(uint32, tag = "5")]
    pub current_epoch_start_block: u32,
    /// is_initialized indicates whether the `EpochInfo` has been initialized
    /// and started ticking.
    /// An `EpochInfo` is initialized when all below conditions are true:
    /// - Not yet initialized
    /// - `BlockHeight` >= 2
    /// - `BlockTime` >= `next_tick`
    #[prost(bool, tag = "6")]
    pub is_initialized: bool,
    /// fast_forward_next_tick specifies whether during initialization, `next_tick`
    /// should be fast-forwarded to be greater than the current block time.
    /// If `false`, the original `next_tick` value is
    /// unchanged during initialization.
    /// If `true`, `next_tick` will be set to the smallest value `x` greater than
    /// the current block time such that `(x - next_tick) % duration = 0`.
    #[prost(bool, tag = "7")]
    pub fast_forward_next_tick: bool,
}
impl ::prost::Name for EpochInfo {
    const NAME: &'static str = "EpochInfo";
    const PACKAGE: &'static str = "dydxprotocol.epochs";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.epochs.EpochInfo".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.epochs.EpochInfo".into()
    }
}
/// GenesisState defines the epochs module's genesis state.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GenesisState {
    /// this line is used by starport scaffolding # genesis/proto/state
    #[prost(message, repeated, tag = "1")]
    pub epoch_info_list: ::prost::alloc::vec::Vec<EpochInfo>,
}
impl ::prost::Name for GenesisState {
    const NAME: &'static str = "GenesisState";
    const PACKAGE: &'static str = "dydxprotocol.epochs";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.epochs.GenesisState".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.epochs.GenesisState".into()
    }
}
/// QueryGetEpochInfoRequest is request type for the GetEpochInfo RPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryGetEpochInfoRequest {
    #[prost(string, tag = "1")]
    pub name: ::prost::alloc::string::String,
}
impl ::prost::Name for QueryGetEpochInfoRequest {
    const NAME: &'static str = "QueryGetEpochInfoRequest";
    const PACKAGE: &'static str = "dydxprotocol.epochs";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.epochs.QueryGetEpochInfoRequest".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.epochs.QueryGetEpochInfoRequest".into()
    }
}
/// QueryEpochInfoResponse is response type for the GetEpochInfo RPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryEpochInfoResponse {
    #[prost(message, optional, tag = "1")]
    pub epoch_info: ::core::option::Option<EpochInfo>,
}
impl ::prost::Name for QueryEpochInfoResponse {
    const NAME: &'static str = "QueryEpochInfoResponse";
    const PACKAGE: &'static str = "dydxprotocol.epochs";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.epochs.QueryEpochInfoResponse".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.epochs.QueryEpochInfoResponse".into()
    }
}
/// QueryAllEpochInfoRequest is request type for the AllEpochInfo RPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryAllEpochInfoRequest {
    #[prost(message, optional, tag = "1")]
    pub pagination: ::core::option::Option<
        super::super::cosmos::base::query::v1beta1::PageRequest,
    >,
}
impl ::prost::Name for QueryAllEpochInfoRequest {
    const NAME: &'static str = "QueryAllEpochInfoRequest";
    const PACKAGE: &'static str = "dydxprotocol.epochs";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.epochs.QueryAllEpochInfoRequest".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.epochs.QueryAllEpochInfoRequest".into()
    }
}
/// QueryEpochInfoAllResponse is response type for the AllEpochInfo RPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryEpochInfoAllResponse {
    #[prost(message, repeated, tag = "1")]
    pub epoch_info: ::prost::alloc::vec::Vec<EpochInfo>,
    #[prost(message, optional, tag = "2")]
    pub pagination: ::core::option::Option<
        super::super::cosmos::base::query::v1beta1::PageResponse,
    >,
}
impl ::prost::Name for QueryEpochInfoAllResponse {
    const NAME: &'static str = "QueryEpochInfoAllResponse";
    const PACKAGE: &'static str = "dydxprotocol.epochs";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.epochs.QueryEpochInfoAllResponse".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.epochs.QueryEpochInfoAllResponse".into()
    }
}
/// Generated client implementations.
pub mod query_client {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value,
    )]
    use tonic::codegen::*;
    use tonic::codegen::http::Uri;
    /// Query defines the gRPC querier service.
    #[derive(Debug, Clone)]
    pub struct QueryClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    #[cfg(feature = "grpc-transport")]
    impl QueryClient<tonic::transport::Channel> {
        /// Attempt to create a new client by connecting to a given endpoint.
        pub async fn connect<D>(dst: D) -> Result<Self, tonic::transport::Error>
        where
            D: TryInto<tonic::transport::Endpoint>,
            D::Error: Into<StdError>,
        {
            let conn = tonic::transport::Endpoint::new(dst)?.connect().await?;
            Ok(Self::new(conn))
        }
    }
    impl<T> QueryClient<T>
    where
        T: tonic::client::GrpcService<tonic::body::BoxBody>,
        T::Error: Into<StdError>,
        T::ResponseBody: Body<Data = Bytes> + std::marker::Send + 'static,
        <T::ResponseBody as Body>::Error: Into<StdError> + std::marker::Send,
    {
        pub fn new(inner: T) -> Self {
            let inner = tonic::client::Grpc::new(inner);
            Self { inner }
        }
        pub fn with_origin(inner: T, origin: Uri) -> Self {
            let inner = tonic::client::Grpc::with_origin(inner, origin);
            Self { inner }
        }
        pub fn with_interceptor<F>(
            inner: T,
            interceptor: F,
        ) -> QueryClient<InterceptedService<T, F>>
        where
            F: tonic::service::Interceptor,
            T::ResponseBody: Default,
            T: tonic::codegen::Service<
                http::Request<tonic::body::BoxBody>,
                Response = http::Response<
                    <T as tonic::client::GrpcService<tonic::body::BoxBody>>::ResponseBody,
                >,
            >,
            <T as tonic::codegen::Service<
                http::Request<tonic::body::BoxBody>,
            >>::Error: Into<StdError> + std::marker::Send + std::marker::Sync,
        {
            QueryClient::new(InterceptedService::new(inner, interceptor))
        }
        /// Compress requests with the given encoding.
        ///
        /// This requires the server to support it otherwise it might respond with an
        /// error.
        #[must_use]
        pub fn send_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.send_compressed(encoding);
            self
        }
        /// Enable decompressing responses.
        #[must_use]
        pub fn accept_compressed(mut self, encoding: CompressionEncoding) -> Self {
            self.inner = self.inner.accept_compressed(encoding);
            self
        }
        /// Limits the maximum size of a decoded message.
        ///
        /// Default: `4MB`
        #[must_use]
        pub fn max_decoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_decoding_message_size(limit);
            self
        }
        /// Limits the maximum size of an encoded message.
        ///
        /// Default: `usize::MAX`
        #[must_use]
        pub fn max_encoding_message_size(mut self, limit: usize) -> Self {
            self.inner = self.inner.max_encoding_message_size(limit);
            self
        }
        /// Queries a EpochInfo by name.
        pub async fn epoch_info(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryGetEpochInfoRequest>,
        ) -> std::result::Result<
            tonic::Response<super::QueryEpochInfoResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/dydxprotocol.epochs.Query/EpochInfo",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("dydxprotocol.epochs.Query", "EpochInfo"));
            self.inner.unary(req, path, codec).await
        }
        /// Queries a list of EpochInfo items.
        pub async fn epoch_info_all(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryAllEpochInfoRequest>,
        ) -> std::result::Result<
            tonic::Response<super::QueryEpochInfoAllResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::unknown(
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/dydxprotocol.epochs.Query/EpochInfoAll",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("dydxprotocol.epochs.Query", "EpochInfoAll"));
            self.inner.unary(req, path, codec).await
        }
    }
}
