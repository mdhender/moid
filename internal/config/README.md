# README

The Config struct holds the configuration for the entire application.

This package loads values from configuration files, the environment, and the command line.

We load values in the following order:

1. configuration files
2. environment variables
3. command line arguments

When we load values, we always overwrite any existing value.
The end result is that command line arguments are the highest priority
and configuration file values are the lowest priority.

There are two variables that must be set before creating the default configuration.

1. The environment (development, test, production). This is required; there is no default value.
2. The path to the configuration files. This is optional; it will default to the current working directory.

We will convert the path to an absolute path if it isn't already.
I've heard that can cause issues on Windows;
I'm using Mac and Linux, so I can't verify.

The environment and path are used to find and load the configuration files.

The environment files are loaded in the following order:

    Filename________________  .gitignore?  safe to store secrets?
    .env                      no           never use for secrets
    .env.{environment}        no           never use for secrets
    .env.local                yes          maybe
    .env.local.{environment}  yes          maybe

Each file will overwrite any values from prior files.
I hope that this means that if you include a value in .env but not in .env.local,
then the value in .env is kept.

After the environment files are loaded, we load values from the environment.
Any values from the environment will overwrite values from the configuration files.

Environment variables are generally named `MOID_VARIABLE_PATH`,
where `MOID_` is the prefix we're searching for,
and `VARIABLE_PATH` is the path from the config file converted
to uppercase and using `_` instead of `.` or `-` to separate elements.
For example, `meta.show-env-files` becomes `MOID_META_SHOW_ENV_FILES`.

Note that we **do not** use the environment name when looking for environment variables.

The last thing we do is use the command line arguments to overwrite values.
The flags are generally named `--variable-path`,
where `variable-path` is the path from the config file converted
to lowercase and using `-` instead of `.` to separate elements.
For example, `meta.show-env-files` becomes `--meta-show-env-files`.

Booleans, of course, get special treatment.
The value of the flag can be blank, `true` or `yes` for true,
and `false` or `no` for false.
Any other value is an error.

The working directory can be set with the environment variable `MOID_WORKING_DIRECTORY` or the command line flag `--working-directory` or the configuration value `working-directory`.
If the working directory is specified, it must be a valid path.
The application will change to the working directory after successfully loading configuration files.

## Command Line Variables
TODO: this section should document the command line variables.
For now, users must look at the configuration code for the complete list.