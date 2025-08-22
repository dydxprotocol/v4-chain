use prost_build::Config;
use regex::Regex;
use std::env;
use std::fs;
use std::io;
use std::path::{Path, PathBuf};

const OUT_DIR: &str = "src";

fn features_patch(dir: impl AsRef<Path>) -> io::Result<()> {
    let regex = "impl(.+)tonic::transport(.+)";
    let replacement = "#[cfg(feature = \"grpc-transport\")]\n    \
                       impl${1}tonic::transport${2}";

    let paths = fs::read_dir(dir)?;

    for entry in paths {
        let path = entry?.path();
        let mut contents = fs::read_to_string(&path)?;

        contents = Regex::new(regex)
            .map_err(io::Error::other)?
            .replace_all(&contents, replacement)
            .to_string();

        fs::write(&path, &contents)?
    }

    Ok(())
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    if std::env::var("V4_PROTO_REBUILD").is_err() {
        return Ok(());
    }

    let mut config = Config::new();
    config.out_dir(OUT_DIR);
    config.include_file("_includes.rs");
    config.enable_type_names();
    let mut path = PathBuf::from(env::var("CARGO_MANIFEST_DIR").map_err(|e| {
        tonic_buf_build::error::TonicBufBuildError {
            message: format!("Failed to get CARGO_MANIFEST_DIR: {e}"),
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

    features_patch(OUT_DIR)?;

    Ok(())
}
