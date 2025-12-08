mod crypto;

use kovi::build_bot;

fn main() {
    crypto::init();
    build_bot!(kovi_plugin_cmd, xiaoshi_kovi_plugin).run();
}
