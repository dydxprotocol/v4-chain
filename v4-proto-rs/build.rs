use prost_build::Config;
use std::env;
use std::path::PathBuf;

fn main() -> Result<(), tonic_buf_build::error::TonicBufBuildError> {
    if std::env::var("V4_PROTO_REBUILD").is_ok() {
        let mut config = Config::new();
        config.out_dir("src");
        config.include_file("_includes.rs");
        config.enable_type_names();
        let mut path = PathBuf::from(env::var("CARGO_MANIFEST_DIR").map_err(|e| {
            tonic_buf_build::error::TonicBufBuildError {
                message: format!("Failed to get CARGO_MANIFEST_DIR: {}", e),
                cause: None,
            }
        })?);
        path.pop();
        tonic_buf_build::compile_from_buf_workspace_with_config(
            tonic_build::configure().build_server(false),
            Some(config),
            tonic_buf_build::TonicBufConfig {
                buf_dir: Some(path),
            },
        )?;
    }

    Ok(())
}
