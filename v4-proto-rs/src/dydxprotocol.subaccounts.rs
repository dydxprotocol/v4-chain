// This file is @generated by prost-build.
/// AssetPositions define an account’s positions of an `Asset`.
/// Therefore they hold any information needed to trade on Spot and Margin.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct AssetPosition {
    /// The `Id` of the `Asset`.
    #[prost(uint32, tag = "1")]
    pub asset_id: u32,
    /// The absolute size of the position in base quantums.
    #[prost(bytes = "vec", tag = "2")]
    pub quantums: ::prost::alloc::vec::Vec<u8>,
    /// The `Index` (either `LongIndex` or `ShortIndex`) of the `Asset` the last
    /// time this position was settled
    /// TODO(DEC-582): pending margin trading being added.
    #[prost(uint64, tag = "3")]
    pub index: u64,
}
impl ::prost::Name for AssetPosition {
    const NAME: &'static str = "AssetPosition";
    const PACKAGE: &'static str = "dydxprotocol.subaccounts";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.subaccounts.AssetPosition".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.subaccounts.AssetPosition".into()
    }
}
/// PerpetualPositions are an account’s positions of a `Perpetual`.
/// Therefore they hold any information needed to trade perpetuals.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct PerpetualPosition {
    /// The `Id` of the `Perpetual`.
    #[prost(uint32, tag = "1")]
    pub perpetual_id: u32,
    /// The size of the position in base quantums.
    #[prost(bytes = "vec", tag = "2")]
    pub quantums: ::prost::alloc::vec::Vec<u8>,
    /// The funding_index of the `Perpetual` the last time this position was
    /// settled.
    #[prost(bytes = "vec", tag = "3")]
    pub funding_index: ::prost::alloc::vec::Vec<u8>,
}
impl ::prost::Name for PerpetualPosition {
    const NAME: &'static str = "PerpetualPosition";
    const PACKAGE: &'static str = "dydxprotocol.subaccounts";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.subaccounts.PerpetualPosition".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.subaccounts.PerpetualPosition".into()
    }
}
/// SubaccountId defines a unique identifier for a Subaccount.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct SubaccountId {
    /// The address of the wallet that owns this subaccount.
    #[prost(string, tag = "1")]
    pub owner: ::prost::alloc::string::String,
    /// The unique number of this subaccount for the owner.
    /// Currently limited to 128*1000 subaccounts per owner.
    #[prost(uint32, tag = "2")]
    pub number: u32,
}
impl ::prost::Name for SubaccountId {
    const NAME: &'static str = "SubaccountId";
    const PACKAGE: &'static str = "dydxprotocol.subaccounts";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.subaccounts.SubaccountId".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.subaccounts.SubaccountId".into()
    }
}
/// Subaccount defines a single sub-account for a given address.
/// Subaccounts are uniquely indexed by a subaccountNumber/owner pair.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct Subaccount {
    /// The Id of the Subaccount
    #[prost(message, optional, tag = "1")]
    pub id: ::core::option::Option<SubaccountId>,
    /// All `AssetPosition`s associated with this subaccount.
    /// Always sorted ascending by `asset_id`.
    #[prost(message, repeated, tag = "2")]
    pub asset_positions: ::prost::alloc::vec::Vec<AssetPosition>,
    /// All `PerpetualPosition`s associated with this subaccount.
    /// Always sorted ascending by `perpetual_id.
    #[prost(message, repeated, tag = "3")]
    pub perpetual_positions: ::prost::alloc::vec::Vec<PerpetualPosition>,
    /// Set by the owner. If true, then margin trades can be made in this
    /// subaccount.
    #[prost(bool, tag = "4")]
    pub margin_enabled: bool,
}
impl ::prost::Name for Subaccount {
    const NAME: &'static str = "Subaccount";
    const PACKAGE: &'static str = "dydxprotocol.subaccounts";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.subaccounts.Subaccount".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.subaccounts.Subaccount".into()
    }
}
/// GenesisState defines the subaccounts module's genesis state.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GenesisState {
    #[prost(message, repeated, tag = "1")]
    pub subaccounts: ::prost::alloc::vec::Vec<Subaccount>,
}
impl ::prost::Name for GenesisState {
    const NAME: &'static str = "GenesisState";
    const PACKAGE: &'static str = "dydxprotocol.subaccounts";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.subaccounts.GenesisState".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.subaccounts.GenesisState".into()
    }
}
/// QueryGetSubaccountRequest is request type for the Query RPC method.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryGetSubaccountRequest {
    #[prost(string, tag = "1")]
    pub owner: ::prost::alloc::string::String,
    #[prost(uint32, tag = "2")]
    pub number: u32,
}
impl ::prost::Name for QueryGetSubaccountRequest {
    const NAME: &'static str = "QueryGetSubaccountRequest";
    const PACKAGE: &'static str = "dydxprotocol.subaccounts";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.subaccounts.QueryGetSubaccountRequest".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.subaccounts.QueryGetSubaccountRequest".into()
    }
}
/// QuerySubaccountResponse is response type for the Query RPC method.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QuerySubaccountResponse {
    #[prost(message, optional, tag = "1")]
    pub subaccount: ::core::option::Option<Subaccount>,
}
impl ::prost::Name for QuerySubaccountResponse {
    const NAME: &'static str = "QuerySubaccountResponse";
    const PACKAGE: &'static str = "dydxprotocol.subaccounts";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.subaccounts.QuerySubaccountResponse".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.subaccounts.QuerySubaccountResponse".into()
    }
}
/// QueryAllSubaccountRequest is request type for the Query RPC method.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryAllSubaccountRequest {
    #[prost(message, optional, tag = "1")]
    pub pagination: ::core::option::Option<
        super::super::cosmos::base::query::v1beta1::PageRequest,
    >,
}
impl ::prost::Name for QueryAllSubaccountRequest {
    const NAME: &'static str = "QueryAllSubaccountRequest";
    const PACKAGE: &'static str = "dydxprotocol.subaccounts";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.subaccounts.QueryAllSubaccountRequest".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.subaccounts.QueryAllSubaccountRequest".into()
    }
}
/// QuerySubaccountAllResponse is response type for the Query RPC method.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QuerySubaccountAllResponse {
    #[prost(message, repeated, tag = "1")]
    pub subaccount: ::prost::alloc::vec::Vec<Subaccount>,
    #[prost(message, optional, tag = "2")]
    pub pagination: ::core::option::Option<
        super::super::cosmos::base::query::v1beta1::PageResponse,
    >,
}
impl ::prost::Name for QuerySubaccountAllResponse {
    const NAME: &'static str = "QuerySubaccountAllResponse";
    const PACKAGE: &'static str = "dydxprotocol.subaccounts";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.subaccounts.QuerySubaccountAllResponse".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.subaccounts.QuerySubaccountAllResponse".into()
    }
}
/// QueryGetWithdrawalAndTransfersBlockedInfoRequest is a request type for
/// fetching information about whether withdrawals and transfers are blocked for
/// a collateral pool associated with the passed in perpetual id.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryGetWithdrawalAndTransfersBlockedInfoRequest {
    #[prost(uint32, tag = "1")]
    pub perpetual_id: u32,
}
impl ::prost::Name for QueryGetWithdrawalAndTransfersBlockedInfoRequest {
    const NAME: &'static str = "QueryGetWithdrawalAndTransfersBlockedInfoRequest";
    const PACKAGE: &'static str = "dydxprotocol.subaccounts";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.subaccounts.QueryGetWithdrawalAndTransfersBlockedInfoRequest"
            .into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.subaccounts.QueryGetWithdrawalAndTransfersBlockedInfoRequest"
            .into()
    }
}
/// QueryGetWithdrawalAndTransfersBlockedInfoRequest is a response type for
/// fetching information about whether withdrawals and transfers are blocked.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryGetWithdrawalAndTransfersBlockedInfoResponse {
    #[prost(uint32, tag = "1")]
    pub negative_tnc_subaccount_seen_at_block: u32,
    #[prost(uint32, tag = "2")]
    pub chain_outage_seen_at_block: u32,
    #[prost(uint32, tag = "3")]
    pub withdrawals_and_transfers_unblocked_at_block: u32,
}
impl ::prost::Name for QueryGetWithdrawalAndTransfersBlockedInfoResponse {
    const NAME: &'static str = "QueryGetWithdrawalAndTransfersBlockedInfoResponse";
    const PACKAGE: &'static str = "dydxprotocol.subaccounts";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.subaccounts.QueryGetWithdrawalAndTransfersBlockedInfoResponse"
            .into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.subaccounts.QueryGetWithdrawalAndTransfersBlockedInfoResponse"
            .into()
    }
}
/// QueryCollateralPoolAddressRequest is the request type for fetching the
/// account address of the collateral pool associated with the passed in
/// perpetual id.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryCollateralPoolAddressRequest {
    #[prost(uint32, tag = "1")]
    pub perpetual_id: u32,
}
impl ::prost::Name for QueryCollateralPoolAddressRequest {
    const NAME: &'static str = "QueryCollateralPoolAddressRequest";
    const PACKAGE: &'static str = "dydxprotocol.subaccounts";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.subaccounts.QueryCollateralPoolAddressRequest".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.subaccounts.QueryCollateralPoolAddressRequest".into()
    }
}
/// QueryCollateralPoolAddressResponse is a response type for fetching the
/// account address of the collateral pool associated with the passed in
/// perpetual id.
#[allow(clippy::derive_partial_eq_without_eq)]
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryCollateralPoolAddressResponse {
    #[prost(string, tag = "1")]
    pub collateral_pool_address: ::prost::alloc::string::String,
}
impl ::prost::Name for QueryCollateralPoolAddressResponse {
    const NAME: &'static str = "QueryCollateralPoolAddressResponse";
    const PACKAGE: &'static str = "dydxprotocol.subaccounts";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.subaccounts.QueryCollateralPoolAddressResponse".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.subaccounts.QueryCollateralPoolAddressResponse".into()
    }
}
/// Generated client implementations.
pub mod query_client {
    #![allow(unused_variables, dead_code, missing_docs, clippy::let_unit_value)]
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
        T::ResponseBody: Body<Data = Bytes> + Send + 'static,
        <T::ResponseBody as Body>::Error: Into<StdError> + Send,
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
            >>::Error: Into<StdError> + Send + Sync,
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
        /// Queries a Subaccount by id
        pub async fn subaccount(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryGetSubaccountRequest>,
        ) -> std::result::Result<
            tonic::Response<super::QuerySubaccountResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/dydxprotocol.subaccounts.Query/Subaccount",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("dydxprotocol.subaccounts.Query", "Subaccount"));
            self.inner.unary(req, path, codec).await
        }
        /// Queries a list of Subaccount items.
        pub async fn subaccount_all(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryAllSubaccountRequest>,
        ) -> std::result::Result<
            tonic::Response<super::QuerySubaccountAllResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/dydxprotocol.subaccounts.Query/SubaccountAll",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new("dydxprotocol.subaccounts.Query", "SubaccountAll"),
                );
            self.inner.unary(req, path, codec).await
        }
        /// Queries information about whether withdrawal and transfers are blocked, and
        /// if so which block they are re-enabled on.
        pub async fn get_withdrawal_and_transfers_blocked_info(
            &mut self,
            request: impl tonic::IntoRequest<
                super::QueryGetWithdrawalAndTransfersBlockedInfoRequest,
            >,
        ) -> std::result::Result<
            tonic::Response<super::QueryGetWithdrawalAndTransfersBlockedInfoResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/dydxprotocol.subaccounts.Query/GetWithdrawalAndTransfersBlockedInfo",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new(
                        "dydxprotocol.subaccounts.Query",
                        "GetWithdrawalAndTransfersBlockedInfo",
                    ),
                );
            self.inner.unary(req, path, codec).await
        }
        /// Queries the collateral pool account address for a perpetual id.
        pub async fn collateral_pool_address(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryCollateralPoolAddressRequest>,
        ) -> std::result::Result<
            tonic::Response<super::QueryCollateralPoolAddressResponse>,
            tonic::Status,
        > {
            self.inner
                .ready()
                .await
                .map_err(|e| {
                    tonic::Status::new(
                        tonic::Code::Unknown,
                        format!("Service was not ready: {}", e.into()),
                    )
                })?;
            let codec = tonic::codec::ProstCodec::default();
            let path = http::uri::PathAndQuery::from_static(
                "/dydxprotocol.subaccounts.Query/CollateralPoolAddress",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new(
                        "dydxprotocol.subaccounts.Query",
                        "CollateralPoolAddress",
                    ),
                );
            self.inner.unary(req, path, codec).await
        }
    }
}