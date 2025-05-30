// This file is @generated by prost-build.
/// PerpetualFeeParams defines the parameters for perpetual fees.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct PerpetualFeeParams {
    /// Sorted fee tiers (lowest requirements first).
    #[prost(message, repeated, tag = "1")]
    pub tiers: ::prost::alloc::vec::Vec<PerpetualFeeTier>,
}
impl ::prost::Name for PerpetualFeeParams {
    const NAME: &'static str = "PerpetualFeeParams";
    const PACKAGE: &'static str = "dydxprotocol.feetiers";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.feetiers.PerpetualFeeParams".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.feetiers.PerpetualFeeParams".into()
    }
}
/// A fee tier for perpetuals
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct PerpetualFeeTier {
    /// Human-readable name of the tier, e.g. "Gold".
    #[prost(string, tag = "1")]
    pub name: ::prost::alloc::string::String,
    /// The trader's absolute volume requirement in quote quantums.
    #[prost(uint64, tag = "2")]
    pub absolute_volume_requirement: u64,
    /// The total volume share requirement.
    #[prost(uint32, tag = "3")]
    pub total_volume_share_requirement_ppm: u32,
    /// The maker volume share requirement.
    #[prost(uint32, tag = "4")]
    pub maker_volume_share_requirement_ppm: u32,
    /// The maker fee once this tier is reached.
    #[prost(sint32, tag = "5")]
    pub maker_fee_ppm: i32,
    /// The taker fee once this tier is reached.
    #[prost(sint32, tag = "6")]
    pub taker_fee_ppm: i32,
}
impl ::prost::Name for PerpetualFeeTier {
    const NAME: &'static str = "PerpetualFeeTier";
    const PACKAGE: &'static str = "dydxprotocol.feetiers";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.feetiers.PerpetualFeeTier".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.feetiers.PerpetualFeeTier".into()
    }
}
/// GenesisState defines the feetiers module's genesis state.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GenesisState {
    /// The parameters for perpetual fees.
    #[prost(message, optional, tag = "1")]
    pub params: ::core::option::Option<PerpetualFeeParams>,
}
impl ::prost::Name for GenesisState {
    const NAME: &'static str = "GenesisState";
    const PACKAGE: &'static str = "dydxprotocol.feetiers";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.feetiers.GenesisState".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.feetiers.GenesisState".into()
    }
}
/// QueryPerpetualFeeParamsRequest is a request type for the PerpetualFeeParams
/// RPC method.
#[derive(Clone, Copy, PartialEq, ::prost::Message)]
pub struct QueryPerpetualFeeParamsRequest {}
impl ::prost::Name for QueryPerpetualFeeParamsRequest {
    const NAME: &'static str = "QueryPerpetualFeeParamsRequest";
    const PACKAGE: &'static str = "dydxprotocol.feetiers";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.feetiers.QueryPerpetualFeeParamsRequest".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.feetiers.QueryPerpetualFeeParamsRequest".into()
    }
}
/// QueryPerpetualFeeParamsResponse is a response type for the PerpetualFeeParams
/// RPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryPerpetualFeeParamsResponse {
    #[prost(message, optional, tag = "1")]
    pub params: ::core::option::Option<PerpetualFeeParams>,
}
impl ::prost::Name for QueryPerpetualFeeParamsResponse {
    const NAME: &'static str = "QueryPerpetualFeeParamsResponse";
    const PACKAGE: &'static str = "dydxprotocol.feetiers";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.feetiers.QueryPerpetualFeeParamsResponse".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.feetiers.QueryPerpetualFeeParamsResponse".into()
    }
}
/// QueryUserFeeTierRequest is a request type for the UserFeeTier RPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryUserFeeTierRequest {
    #[prost(string, tag = "1")]
    pub user: ::prost::alloc::string::String,
}
impl ::prost::Name for QueryUserFeeTierRequest {
    const NAME: &'static str = "QueryUserFeeTierRequest";
    const PACKAGE: &'static str = "dydxprotocol.feetiers";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.feetiers.QueryUserFeeTierRequest".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.feetiers.QueryUserFeeTierRequest".into()
    }
}
/// QueryUserFeeTierResponse is a request type for the UserFeeTier RPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryUserFeeTierResponse {
    /// Index of the fee tier in the list queried from PerpetualFeeParams.
    #[prost(uint32, tag = "1")]
    pub index: u32,
    #[prost(message, optional, tag = "2")]
    pub tier: ::core::option::Option<PerpetualFeeTier>,
}
impl ::prost::Name for QueryUserFeeTierResponse {
    const NAME: &'static str = "QueryUserFeeTierResponse";
    const PACKAGE: &'static str = "dydxprotocol.feetiers";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.feetiers.QueryUserFeeTierResponse".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.feetiers.QueryUserFeeTierResponse".into()
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
        /// Queries the PerpetualFeeParams.
        pub async fn perpetual_fee_params(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryPerpetualFeeParamsRequest>,
        ) -> std::result::Result<
            tonic::Response<super::QueryPerpetualFeeParamsResponse>,
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
                "/dydxprotocol.feetiers.Query/PerpetualFeeParams",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new("dydxprotocol.feetiers.Query", "PerpetualFeeParams"),
                );
            self.inner.unary(req, path, codec).await
        }
        /// Queries a user's fee tier
        pub async fn user_fee_tier(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryUserFeeTierRequest>,
        ) -> std::result::Result<
            tonic::Response<super::QueryUserFeeTierResponse>,
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
                "/dydxprotocol.feetiers.Query/UserFeeTier",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("dydxprotocol.feetiers.Query", "UserFeeTier"));
            self.inner.unary(req, path, codec).await
        }
    }
}
/// MsgUpdatePerpetualFeeParams is the Msg/UpdatePerpetualFeeParams request type.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgUpdatePerpetualFeeParams {
    #[prost(string, tag = "1")]
    pub authority: ::prost::alloc::string::String,
    /// Defines the parameters to update. All parameters must be supplied.
    #[prost(message, optional, tag = "2")]
    pub params: ::core::option::Option<PerpetualFeeParams>,
}
impl ::prost::Name for MsgUpdatePerpetualFeeParams {
    const NAME: &'static str = "MsgUpdatePerpetualFeeParams";
    const PACKAGE: &'static str = "dydxprotocol.feetiers";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.feetiers.MsgUpdatePerpetualFeeParams".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.feetiers.MsgUpdatePerpetualFeeParams".into()
    }
}
/// MsgUpdatePerpetualFeeParamsResponse is the Msg/UpdatePerpetualFeeParams
/// response type.
#[derive(Clone, Copy, PartialEq, ::prost::Message)]
pub struct MsgUpdatePerpetualFeeParamsResponse {}
impl ::prost::Name for MsgUpdatePerpetualFeeParamsResponse {
    const NAME: &'static str = "MsgUpdatePerpetualFeeParamsResponse";
    const PACKAGE: &'static str = "dydxprotocol.feetiers";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.feetiers.MsgUpdatePerpetualFeeParamsResponse".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.feetiers.MsgUpdatePerpetualFeeParamsResponse".into()
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
    #[cfg(feature = "grpc-transport")]
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
        /// UpdatePerpetualFeeParams updates the PerpetualFeeParams in state.
        pub async fn update_perpetual_fee_params(
            &mut self,
            request: impl tonic::IntoRequest<super::MsgUpdatePerpetualFeeParams>,
        ) -> std::result::Result<
            tonic::Response<super::MsgUpdatePerpetualFeeParamsResponse>,
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
                "/dydxprotocol.feetiers.Msg/UpdatePerpetualFeeParams",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new(
                        "dydxprotocol.feetiers.Msg",
                        "UpdatePerpetualFeeParams",
                    ),
                );
            self.inner.unary(req, path, codec).await
        }
    }
}
