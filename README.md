# murult

Bot for Materia Ultimate Raid (MUR) discord server.

## Usage

Please look at the provided `config.example.yaml` for the list of required configuration.
They are hopefully self-explanatory.
Aside from that, the server accepts some arguments:

1. `--sleep <X>` how many seconds to sleep before refreshing the channel
2. `--config <path>` alternative path to the `config.yaml` file.
   Defaults to a `config.yaml` in current working directory.
3. `--once` causes the server to run the loop once and exits.
   Mainly useful for debugging purposes.
