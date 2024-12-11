// This file is @generated by prost-build.
/// MarketParam represents the x/prices configuration for markets, including
/// representing price values, resolving markets on individual exchanges, and
/// generating price updates. This configuration is specific to the quote
/// currency.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MarketParam {
    /// Unique, sequentially-generated value.
    #[prost(uint32, tag = "1")]
    pub id: u32,
    /// The human-readable name of the market pair (e.g. `BTC-USD`).
    #[prost(string, tag = "2")]
    pub pair: ::prost::alloc::string::String,
    /// Static value. The exponent of the price.
    /// For example if `Exponent == -5` then a `Value` of `1,000,000,000`
    /// represents ``$10,000`. Therefore `10 ^ Exponent` represents the smallest
    /// price step (in dollars) that can be recorded.
    ///
    /// Deprecated since v8.x. This value is now determined from the marketmap.
    #[deprecated]
    #[prost(sint32, tag = "3")]
    pub exponent: i32,
    /// The minimum number of exchanges that should be reporting a live price for
    /// a price update to be considered valid.
    ///
    /// Deprecated since v8.x. This value is now determined from the marketmap.
    #[prost(uint32, tag = "4")]
    pub min_exchanges: u32,
    /// The minimum allowable change in `price` value that would cause a price
    /// update on the network. Measured as `1e-6` (parts per million).
    #[prost(uint32, tag = "5")]
    pub min_price_change_ppm: u32,
    /// A string of json that encodes the configuration for resolving the price
    /// of this market on various exchanges.
    ///
    /// Deprecated since v8.x. This is now determined from the marketmap.
    #[prost(string, tag = "6")]
    pub exchange_config_json: ::prost::alloc::string::String,
}
impl ::prost::Name for MarketParam {
    const NAME: &'static str = "MarketParam";
    const PACKAGE: &'static str = "dydxprotocol.prices";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.prices.MarketParam".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.prices.MarketParam".into()
    }
}
/// MarketPrice is used by the application to store/retrieve oracle price.
#[derive(Clone, Copy, PartialEq, ::prost::Message)]
pub struct MarketPrice {
    /// Unique, sequentially-generated value that matches `MarketParam`.
    #[prost(uint32, tag = "1")]
    pub id: u32,
    /// Static value. The exponent of the price. See the comment on the duplicate
    /// MarketParam field for more information.
    ///
    /// As of v7.1.x, this value is determined from the marketmap instead of
    /// needing to match the MarketParam field.
    #[prost(sint32, tag = "2")]
    pub exponent: i32,
    /// The variable value that is updated by oracle price updates. `0` if it has
    /// never been updated, `>0` otherwise.
    #[prost(uint64, tag = "3")]
    pub price: u64,
}
impl ::prost::Name for MarketPrice {
    const NAME: &'static str = "MarketPrice";
    const PACKAGE: &'static str = "dydxprotocol.prices";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.prices.MarketPrice".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.prices.MarketPrice".into()
    }
}
/// GenesisState defines the prices module's genesis state.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct GenesisState {
    #[prost(message, repeated, tag = "1")]
    pub market_params: ::prost::alloc::vec::Vec<MarketParam>,
    #[prost(message, repeated, tag = "2")]
    pub market_prices: ::prost::alloc::vec::Vec<MarketPrice>,
}
impl ::prost::Name for GenesisState {
    const NAME: &'static str = "GenesisState";
    const PACKAGE: &'static str = "dydxprotocol.prices";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.prices.GenesisState".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.prices.GenesisState".into()
    }
}
/// QueryMarketPriceRequest is request type for the Query/Params `MarketPrice`
/// RPC method.
#[derive(Clone, Copy, PartialEq, ::prost::Message)]
pub struct QueryMarketPriceRequest {
    #[prost(uint32, tag = "1")]
    pub id: u32,
}
impl ::prost::Name for QueryMarketPriceRequest {
    const NAME: &'static str = "QueryMarketPriceRequest";
    const PACKAGE: &'static str = "dydxprotocol.prices";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.prices.QueryMarketPriceRequest".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.prices.QueryMarketPriceRequest".into()
    }
}
/// QueryMarketPriceResponse is response type for the Query/Params `MarketPrice`
/// RPC method.
#[derive(Clone, Copy, PartialEq, ::prost::Message)]
pub struct QueryMarketPriceResponse {
    #[prost(message, optional, tag = "1")]
    pub market_price: ::core::option::Option<MarketPrice>,
}
impl ::prost::Name for QueryMarketPriceResponse {
    const NAME: &'static str = "QueryMarketPriceResponse";
    const PACKAGE: &'static str = "dydxprotocol.prices";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.prices.QueryMarketPriceResponse".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.prices.QueryMarketPriceResponse".into()
    }
}
/// QueryAllMarketPricesRequest is request type for the Query/Params
/// `AllMarketPrices` RPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryAllMarketPricesRequest {
    #[prost(message, optional, tag = "1")]
    pub pagination: ::core::option::Option<
        super::super::cosmos::base::query::v1beta1::PageRequest,
    >,
}
impl ::prost::Name for QueryAllMarketPricesRequest {
    const NAME: &'static str = "QueryAllMarketPricesRequest";
    const PACKAGE: &'static str = "dydxprotocol.prices";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.prices.QueryAllMarketPricesRequest".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.prices.QueryAllMarketPricesRequest".into()
    }
}
/// QueryAllMarketPricesResponse is response type for the Query/Params
/// `AllMarketPrices` RPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryAllMarketPricesResponse {
    #[prost(message, repeated, tag = "1")]
    pub market_prices: ::prost::alloc::vec::Vec<MarketPrice>,
    #[prost(message, optional, tag = "2")]
    pub pagination: ::core::option::Option<
        super::super::cosmos::base::query::v1beta1::PageResponse,
    >,
}
impl ::prost::Name for QueryAllMarketPricesResponse {
    const NAME: &'static str = "QueryAllMarketPricesResponse";
    const PACKAGE: &'static str = "dydxprotocol.prices";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.prices.QueryAllMarketPricesResponse".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.prices.QueryAllMarketPricesResponse".into()
    }
}
/// QueryMarketParamsRequest is request type for the Query/Params `MarketParams`
/// RPC method.
#[derive(Clone, Copy, PartialEq, ::prost::Message)]
pub struct QueryMarketParamRequest {
    #[prost(uint32, tag = "1")]
    pub id: u32,
}
impl ::prost::Name for QueryMarketParamRequest {
    const NAME: &'static str = "QueryMarketParamRequest";
    const PACKAGE: &'static str = "dydxprotocol.prices";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.prices.QueryMarketParamRequest".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.prices.QueryMarketParamRequest".into()
    }
}
/// QueryMarketParamResponse is response type for the Query/Params `MarketParams`
/// RPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryMarketParamResponse {
    #[prost(message, optional, tag = "1")]
    pub market_param: ::core::option::Option<MarketParam>,
}
impl ::prost::Name for QueryMarketParamResponse {
    const NAME: &'static str = "QueryMarketParamResponse";
    const PACKAGE: &'static str = "dydxprotocol.prices";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.prices.QueryMarketParamResponse".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.prices.QueryMarketParamResponse".into()
    }
}
/// QueryAllMarketParamsRequest is request type for the Query/Params
/// `AllMarketParams` RPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryAllMarketParamsRequest {
    #[prost(message, optional, tag = "1")]
    pub pagination: ::core::option::Option<
        super::super::cosmos::base::query::v1beta1::PageRequest,
    >,
}
impl ::prost::Name for QueryAllMarketParamsRequest {
    const NAME: &'static str = "QueryAllMarketParamsRequest";
    const PACKAGE: &'static str = "dydxprotocol.prices";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.prices.QueryAllMarketParamsRequest".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.prices.QueryAllMarketParamsRequest".into()
    }
}
/// QueryAllMarketParamsResponse is response type for the Query/Params
/// `AllMarketParams` RPC method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct QueryAllMarketParamsResponse {
    #[prost(message, repeated, tag = "1")]
    pub market_params: ::prost::alloc::vec::Vec<MarketParam>,
    #[prost(message, optional, tag = "2")]
    pub pagination: ::core::option::Option<
        super::super::cosmos::base::query::v1beta1::PageResponse,
    >,
}
impl ::prost::Name for QueryAllMarketParamsResponse {
    const NAME: &'static str = "QueryAllMarketParamsResponse";
    const PACKAGE: &'static str = "dydxprotocol.prices";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.prices.QueryAllMarketParamsResponse".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.prices.QueryAllMarketParamsResponse".into()
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
        /// Queries a MarketPrice by id.
        pub async fn market_price(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryMarketPriceRequest>,
        ) -> std::result::Result<
            tonic::Response<super::QueryMarketPriceResponse>,
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
                "/dydxprotocol.prices.Query/MarketPrice",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("dydxprotocol.prices.Query", "MarketPrice"));
            self.inner.unary(req, path, codec).await
        }
        /// Queries a list of MarketPrice items.
        pub async fn all_market_prices(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryAllMarketPricesRequest>,
        ) -> std::result::Result<
            tonic::Response<super::QueryAllMarketPricesResponse>,
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
                "/dydxprotocol.prices.Query/AllMarketPrices",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("dydxprotocol.prices.Query", "AllMarketPrices"));
            self.inner.unary(req, path, codec).await
        }
        /// Queries a MarketParam by id.
        pub async fn market_param(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryMarketParamRequest>,
        ) -> std::result::Result<
            tonic::Response<super::QueryMarketParamResponse>,
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
                "/dydxprotocol.prices.Query/MarketParam",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("dydxprotocol.prices.Query", "MarketParam"));
            self.inner.unary(req, path, codec).await
        }
        /// Queries a list of MarketParam items.
        pub async fn all_market_params(
            &mut self,
            request: impl tonic::IntoRequest<super::QueryAllMarketParamsRequest>,
        ) -> std::result::Result<
            tonic::Response<super::QueryAllMarketParamsResponse>,
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
                "/dydxprotocol.prices.Query/AllMarketParams",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("dydxprotocol.prices.Query", "AllMarketParams"));
            self.inner.unary(req, path, codec).await
        }
    }
}
/// MsgCreateOracleMarket is a message used by x/gov for creating a new oracle
/// market.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgCreateOracleMarket {
    /// The address that controls the module.
    #[prost(string, tag = "1")]
    pub authority: ::prost::alloc::string::String,
    /// `params` defines parameters for the new oracle market.
    #[prost(message, optional, tag = "2")]
    pub params: ::core::option::Option<MarketParam>,
}
impl ::prost::Name for MsgCreateOracleMarket {
    const NAME: &'static str = "MsgCreateOracleMarket";
    const PACKAGE: &'static str = "dydxprotocol.prices";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.prices.MsgCreateOracleMarket".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.prices.MsgCreateOracleMarket".into()
    }
}
/// MsgCreateOracleMarketResponse defines the CreateOracleMarket response type.
#[derive(Clone, Copy, PartialEq, ::prost::Message)]
pub struct MsgCreateOracleMarketResponse {}
impl ::prost::Name for MsgCreateOracleMarketResponse {
    const NAME: &'static str = "MsgCreateOracleMarketResponse";
    const PACKAGE: &'static str = "dydxprotocol.prices";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.prices.MsgCreateOracleMarketResponse".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.prices.MsgCreateOracleMarketResponse".into()
    }
}
/// MsgUpdateMarketPrices is a request type for the UpdateMarketPrices method.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgUpdateMarketPrices {
    #[prost(message, repeated, tag = "1")]
    pub market_price_updates: ::prost::alloc::vec::Vec<
        msg_update_market_prices::MarketPrice,
    >,
}
/// Nested message and enum types in `MsgUpdateMarketPrices`.
pub mod msg_update_market_prices {
    /// MarketPrice represents a price update for a single market
    #[derive(Clone, Copy, PartialEq, ::prost::Message)]
    pub struct MarketPrice {
        /// The id of market to update
        #[prost(uint32, tag = "1")]
        pub market_id: u32,
        /// The updated price
        #[prost(uint64, tag = "2")]
        pub price: u64,
    }
    impl ::prost::Name for MarketPrice {
        const NAME: &'static str = "MarketPrice";
        const PACKAGE: &'static str = "dydxprotocol.prices";
        fn full_name() -> ::prost::alloc::string::String {
            "dydxprotocol.prices.MsgUpdateMarketPrices.MarketPrice".into()
        }
        fn type_url() -> ::prost::alloc::string::String {
            "/dydxprotocol.prices.MsgUpdateMarketPrices.MarketPrice".into()
        }
    }
}
impl ::prost::Name for MsgUpdateMarketPrices {
    const NAME: &'static str = "MsgUpdateMarketPrices";
    const PACKAGE: &'static str = "dydxprotocol.prices";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.prices.MsgUpdateMarketPrices".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.prices.MsgUpdateMarketPrices".into()
    }
}
/// MsgUpdateMarketPricesResponse defines the MsgUpdateMarketPrices response
/// type.
#[derive(Clone, Copy, PartialEq, ::prost::Message)]
pub struct MsgUpdateMarketPricesResponse {}
impl ::prost::Name for MsgUpdateMarketPricesResponse {
    const NAME: &'static str = "MsgUpdateMarketPricesResponse";
    const PACKAGE: &'static str = "dydxprotocol.prices";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.prices.MsgUpdateMarketPricesResponse".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.prices.MsgUpdateMarketPricesResponse".into()
    }
}
/// MsgUpdateMarketParam is a message used by x/gov for updating the parameters
/// of an oracle market.
#[derive(Clone, PartialEq, ::prost::Message)]
pub struct MsgUpdateMarketParam {
    #[prost(string, tag = "1")]
    pub authority: ::prost::alloc::string::String,
    /// The market param to update. Each field must be set.
    #[prost(message, optional, tag = "2")]
    pub market_param: ::core::option::Option<MarketParam>,
}
impl ::prost::Name for MsgUpdateMarketParam {
    const NAME: &'static str = "MsgUpdateMarketParam";
    const PACKAGE: &'static str = "dydxprotocol.prices";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.prices.MsgUpdateMarketParam".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.prices.MsgUpdateMarketParam".into()
    }
}
/// MsgUpdateMarketParamResponse defines the UpdateMarketParam response type.
#[derive(Clone, Copy, PartialEq, ::prost::Message)]
pub struct MsgUpdateMarketParamResponse {}
impl ::prost::Name for MsgUpdateMarketParamResponse {
    const NAME: &'static str = "MsgUpdateMarketParamResponse";
    const PACKAGE: &'static str = "dydxprotocol.prices";
    fn full_name() -> ::prost::alloc::string::String {
        "dydxprotocol.prices.MsgUpdateMarketParamResponse".into()
    }
    fn type_url() -> ::prost::alloc::string::String {
        "/dydxprotocol.prices.MsgUpdateMarketParamResponse".into()
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
        /// UpdateMarketPrices updates the oracle price of a market relative to
        /// quoteCurrency.
        pub async fn update_market_prices(
            &mut self,
            request: impl tonic::IntoRequest<super::MsgUpdateMarketPrices>,
        ) -> std::result::Result<
            tonic::Response<super::MsgUpdateMarketPricesResponse>,
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
                "/dydxprotocol.prices.Msg/UpdateMarketPrices",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new("dydxprotocol.prices.Msg", "UpdateMarketPrices"),
                );
            self.inner.unary(req, path, codec).await
        }
        /// CreateOracleMarket creates a new oracle market.
        pub async fn create_oracle_market(
            &mut self,
            request: impl tonic::IntoRequest<super::MsgCreateOracleMarket>,
        ) -> std::result::Result<
            tonic::Response<super::MsgCreateOracleMarketResponse>,
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
                "/dydxprotocol.prices.Msg/CreateOracleMarket",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(
                    GrpcMethod::new("dydxprotocol.prices.Msg", "CreateOracleMarket"),
                );
            self.inner.unary(req, path, codec).await
        }
        /// UpdateMarketParams allows governance to update the parameters of an
        /// oracle market.
        pub async fn update_market_param(
            &mut self,
            request: impl tonic::IntoRequest<super::MsgUpdateMarketParam>,
        ) -> std::result::Result<
            tonic::Response<super::MsgUpdateMarketParamResponse>,
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
                "/dydxprotocol.prices.Msg/UpdateMarketParam",
            );
            let mut req = request.into_request();
            req.extensions_mut()
                .insert(GrpcMethod::new("dydxprotocol.prices.Msg", "UpdateMarketParam"));
            self.inner.unary(req, path, codec).await
        }
    }
}