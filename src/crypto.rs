/// Initialize rustls crypto provider
pub fn init() {
    rustls::crypto::ring::default_provider()
        .install_default()
        .expect("Failed to install rustls crypto provider");
}