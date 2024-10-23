// This file is @generated by prost-build.
/// BlockInfo stores information about a block
#[derive(Clone, Copy, PartialEq, ::prost::Message)]
pub struct BlockInfo {
    #[prost(uint32, tag = "1")]
    pub height: u32,
    #[prost(message, optional, tag = "2")]
    pub timestamp: ::core::option::Option<::prost_types::Timestamp>,
}
impl ::prost::Name for BlockInfo {
    const NAME: &'static str = "BlockInfo";
    const PACKAGE: &'static str = "dydxprotocol.blocktime";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.blocktime.BlockInfo".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.blocktime.BlockInfo".into()
    }
}
/// AllDowntimeInfo stores information for all downtime durations.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct AllDowntimeInfo {
    /// The downtime information for each tracked duration. Sorted by duration,
    /// ascending. (i.e. the same order as they appear in DowntimeParams).
    #[prost(message, repeated, tag = "1")]
    pub infos: ::prost::alloc::vec::Vec<all_downtime_info::DowntimeInfo>,
}
/// Nested message and enum types in `AllDowntimeInfo`.
pub mod all_downtime_info {
    /// Stores information about downtime. block_info corresponds to the most
    /// recent block at which a downtime occurred.
    #[derive(Clone, Copy, PartialEq, ::prost::Message)]
    pub struct DowntimeInfo {
        #[prost(message, optional, tag = "1")]
        pub duration: ::core::option::Option<::prost_types::Duration>,
        #[prost(message, optional, tag = "2")]
        pub block_info: ::core::option::Option<super::BlockInfo>,
    }
    impl ::prost::Name for DowntimeInfo {
        const NAME: &'static str = "DowntimeInfo";
        const PACKAGE: &'static str = "dydxprotocol.blocktime";
        fn full_name() -> ::prost::alloc::string::String {
            "dydxprotocol.blocktime.AllDowntimeInfo.DowntimeInfo".into()
        }
        fn type_url() -> ::prost::alloc::string::String {
            "/dydxprotocol.blocktime.AllDowntimeInfo.DowntimeInfo".into()
        }
    }
}
impl ::prost::Name for AllDowntimeInfo {
    const NAME: &'static str = "AllDowntimeInfo";
    const PACKAGE: &'static str = "dydxprotocol.blocktime";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.blocktime.AllDowntimeInfo".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.blocktime.AllDowntimeInfo".into()
    }
}
/// DowntimeParams defines the parameters for downtime.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct DowntimeParams {
    /// Durations tracked for downtime. The durations must be sorted from
    /// shortest to longest and must all be positive.
    #[prost(message, repeated, tag = "1")]
    pub durations: ::prost::alloc::vec::Vec<::prost_types::Duration>,
}
impl ::prost::Name for DowntimeParams {
    const NAME: &'static str = "DowntimeParams";
    const PACKAGE: &'static str = "dydxprotocol.blocktime";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.blocktime.DowntimeParams".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.blocktime.DowntimeParams".into()
    }
}
/// GenesisState defines the blocktime module's genesis state.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GenesisState {
    #[prost(message, optional, tag = "1")]
    pub params: ::core::option::Option<DowntimeParams>,
}
impl ::prost::Name for GenesisState {
    const NAME: &'static str = "GenesisState";
    const PACKAGE: &'static str = "dydxprotocol.blocktime";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.blocktime.GenesisState".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.blocktime.GenesisState".into()
    }
}
/// QueryDowntimeParamsRequest is a request type for the DowntimeParams
/// RPC method.
#[derive(Clone, Copy, PartialEq, ::prost::Message)]
pub struct QueryDowntimeParamsRequest {}
impl ::prost::Name for QueryDowntimeParamsRequest {
    const NAME: &'static str = "QueryDowntimeParamsRequest";
    const PACKAGE: &'static str = "dydxprotocol.blocktime";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.blocktime.QueryDowntimeParamsRequest".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.blocktime.QueryDowntimeParamsRequest".into()
    }
}
/// QueryDowntimeParamsResponse is a response type for the DowntimeParams
/// RPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryDowntimeParamsResponse {
    #[prost(message, optional, tag = "1")]
    pub params: ::core::option::Option<DowntimeParams>,
}
impl ::prost::Name for QueryDowntimeParamsResponse {
    const NAME: &'static str = "QueryDowntimeParamsResponse";
    const PACKAGE: &'static str = "dydxprotocol.blocktime";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.blocktime.QueryDowntimeParamsResponse".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.blocktime.QueryDowntimeParamsResponse".into()
    }
}
/// QueryPreviousBlockInfoRequest is a request type for the PreviousBlockInfo
/// RPC method.
#[derive(Clone, Copy, PartialEq, ::prost::Message)]
pub struct QueryPreviousBlockInfoRequest {}
impl ::prost::Name for QueryPreviousBlockInfoRequest {
    const NAME: &'static str = "QueryPreviousBlockInfoRequest";
    const PACKAGE: &'static str = "dydxprotocol.blocktime";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.blocktime.QueryPreviousBlockInfoRequest".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.blocktime.QueryPreviousBlockInfoRequest".into()
    }
}
/// QueryPreviousBlockInfoResponse is a request type for the PreviousBlockInfo
/// RPC method.
#[derive(Clone, Copy, PartialEq, ::prost::Message)]
pub struct QueryPreviousBlockInfoResponse {
    #[prost(message, optional, tag = "1")]
    pub info: ::core::option::Option<BlockInfo>,
}
impl ::prost::Name for QueryPreviousBlockInfoResponse {
    const NAME: &'static str = "QueryPreviousBlockInfoResponse";
    const PACKAGE: &'static str = "dydxprotocol.blocktime";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.blocktime.QueryPreviousBlockInfoResponse".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.blocktime.QueryPreviousBlockInfoResponse".into()
    }
}
/// QueryAllDowntimeInfoRequest is a request type for the AllDowntimeInfo
/// RPC method.
#[derive(Clone, Copy, PartialEq, ::prost::Message)]
pub struct QueryAllDowntimeInfoRequest {}
impl ::prost::Name for QueryAllDowntimeInfoRequest {
    const NAME: &'static str = "QueryAllDowntimeInfoRequest";
    const PACKAGE: &'static str = "dydxprotocol.blocktime";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.blocktime.QueryAllDowntimeInfoRequest".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.blocktime.QueryAllDowntimeInfoRequest".into()
    }
}
/// QueryAllDowntimeInfoResponse is a request type for the AllDowntimeInfo
/// RPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryAllDowntimeInfoResponse {
    #[prost(message, optional, tag = "1")]
    pub info: ::core::option::Option<AllDowntimeInfo>,
}
impl ::prost::Name for QueryAllDowntimeInfoResponse {
    const NAME: &'static str = "QueryAllDowntimeInfoResponse";
    const PACKAGE: &'static str = "dydxprotocol.blocktime";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.blocktime.QueryAllDowntimeInfoResponse".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.blocktime.QueryAllDowntimeInfoResponse".into()
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
        /// Queries the DowntimeParams.
        pub async fn downtime_params(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryDowntimeParamsRequest>,
        ) -> std::result::Result<
            tonic::Response<super::QueryDowntimeParamsResponse>,
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
                "/dydxprotocol.blocktime.Query/DowntimeParams",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new("dydxprotocol.blocktime.Query", "DowntimeParams"),
                );
            self.inner.unary(req, path, codec).await
        }
        /// Queries the information of the previous block
        pub async fn previous_block_info(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryPreviousBlockInfoRequest>,
        ) -> std::result::Result<
            tonic::Response<super::QueryPreviousBlockInfoResponse>,
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
                "/dydxprotocol.blocktime.Query/PreviousBlockInfo",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new("dydxprotocol.blocktime.Query", "PreviousBlockInfo"),
                );
            self.inner.unary(req, path, codec).await
        }
        /// Queries all recorded downtime info.
        pub async fn all_downtime_info(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryAllDowntimeInfoRequest>,
        ) -> std::result::Result<
            tonic::Response<super::QueryAllDowntimeInfoResponse>,
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
                "/dydxprotocol.blocktime.Query/AllDowntimeInfo",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new("dydxprotocol.blocktime.Query", "AllDowntimeInfo"),
                );
            self.inner.unary(req, path, codec).await
        }
    }
}
/// MsgUpdateDowntimeParams is the Msg/UpdateDowntimeParams request type.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgUpdateDowntimeParams {
    #[prost(string, tag = "1")]
    pub authority: ::prost::alloc::string::String,
    /// Defines the parameters to update. All parameters must be supplied.
    #[prost(message, optional, tag = "2")]
    pub params: ::core::option::Option<DowntimeParams>,
}
impl ::prost::Name for MsgUpdateDowntimeParams {
    const NAME: &'static str = "MsgUpdateDowntimeParams";
    const PACKAGE: &'static str = "dydxprotocol.blocktime";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.blocktime.MsgUpdateDowntimeParams".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.blocktime.MsgUpdateDowntimeParams".into()
    }
}
/// MsgUpdateDowntimeParamsResponse is the Msg/UpdateDowntimeParams response
/// type.
#[derive(Clone, Copy, PartialEq, ::prost::Message)]
pub struct MsgUpdateDowntimeParamsResponse {}
impl ::prost::Name for MsgUpdateDowntimeParamsResponse {
    const NAME: &'static str = "MsgUpdateDowntimeParamsResponse";
    const PACKAGE: &'static str = "dydxprotocol.blocktime";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.blocktime.MsgUpdateDowntimeParamsResponse".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.blocktime.MsgUpdateDowntimeParamsResponse".into()
    }
}
/// Generated client implementations.
pub mod msg_client {
    #![allow(
        unused_variables,
        dead_code,
        missing_docs,
        clippy::wildcard_imports,
        clippy::let_unit_value,
    )]
    use tonic::codegen::*;
    use tonic::codegen::http::Uri;
    /// Msg defines the Msg service.
    #[derive(Debug, Clone)]
    pub struct MsgClient<T> {
        inner: tonic::client::Grpc<T>,
    }
    impl MsgClient<tonic::transport::Channel> {
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
    impl<T> MsgClient<T>
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
        ) -> MsgClient<InterceptedService<T, F>>
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
            MsgClient::new(InterceptedService::new(inner, interceptor))
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
        /// UpdateDowntimeParams updates the DowntimeParams in state.
        pub async fn update_downtime_params(
            &mut self,
            request: impl tonic::IntoRequest<super::MsgUpdateDowntimeParams>,
        ) -> std::result::Result<
            tonic::Response<super::MsgUpdateDowntimeParamsResponse>,
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
                "/dydxprotocol.blocktime.Msg/UpdateDowntimeParams",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new("dydxprotocol.blocktime.Msg", "UpdateDowntimeParams"),
                );
            self.inner.unary(req, path, codec).await
        }
    }
}
