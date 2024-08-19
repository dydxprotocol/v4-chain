This directory contains files to configure version names used for upgrade tests and internal nodes. These versions are used for things such as upgrade proposals and cosmovisor directory names. When updating tests and images for a new version, most of the time you just need to change the `VERSION_CURRENT`, `VERSION_PREUPGRADE`, and `VERSION_FULL_NAME_PREUPGRADE` files.

`VERSION_CURRENT` is the upgrade name for the current version. e.g. v6.0.0 is the upgrade name for v6.0.0, v6.0.1, v6.0.0-rc0, etc.
`VERSION_PREUPGRADE` is the upgrade name for the preupgrade version. e.g. v6.0.0 is the upgrade name for v6.0.0, v6.0.1, v6.0.0-rc0, etc.
`VERSION_FULL_NAME_PREUPGRADE` is the fully qualified preupgrade version name. This is used in the download url. e.g. v6.0.0-rc0.
