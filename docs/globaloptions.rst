.. _global_options:

Global Options
==============

Global options are command-line flags that are valid for any command
and allow you to customize various aspects of the ``rack`` tool at runtime.
These options may override configuration file or environment variables you have
set previously, change output format or other aspects of the tool.

For example:

::

    rack <service> <command> <subcommand> --output json [flags]

Would result in the command returning a JSON_ formatted output.

Options
-------

``--output``
~~~~~~~~~~~~

  (string) The format in which to return the output. Options are: table, json, csv. Default is 'table'.

``json``
^^^^^^^^

  Return output in JSON.

  Given::

      rack servers instance list
      ID	        Name		Status	Public IPv4	Private IPv4	Image	Flavor
      GUID	my_server	ACTIVE	101.130.19.31	10.208.128.233	GUID	io1-30

  Adding the ``--output json`` flag returns::

    rack servers instance list --output json
    [
      {
        "Flavor": "io1-30",
        "ID": "GUID",
        "Image": "GUID",
        "Name": "my_server",
        "Private IPv4": "10.208.128.233",
        "Public IPv4": "101.130.19.31",
        "Status": "ACTIVE"
      }
    ]

  When the output pipe is **not** a tty, the JSON is no longer "pretty printed" and
  can be used when passing straight into other commands that require a JSON_
  payload.

``table``
^^^^^^^^^

  Return output in tabular format. Default output format for ``rack``.

  Given::

      rack servers instance list
      ID	        Name		Status	Public IPv4	Private IPv4	Image	Flavor
      GUID	my_server	ACTIVE	101.130.19.31	10.208.128.233	GUID	io1-30

  This presents a well formatted table with headers.

  You can add the ``--output table`` option if you have set defaults to JSON,
  CSV, and so on elsewhere. You can use the ``--no-header`` option to output
  without headers.

``csv``
^^^^^^^

  Return output in csv format.

  CSV, or comma separated output is useful for passing to other operating system
  tools, importing into Excel, Google Sheets, or another data tool.

  Given::

      rack servers instance list
      ID	        Name		Status	Public IPv4	Private IPv4	Image	Flavor
      GUID	my_server	ACTIVE	101.130.19.31	10.208.128.233	GUID	io1-30

  Adding the ``--output csv`` option::

      rack servers instance list --output csv
      ID,Name,Status,Public IPv4,Private IPv4,Image,Flavor
      GUID,my_server,ACTIVE,101.130.19.31,10.208.128.233,GUID,io1-30

  This presents a compact format with appropriate CSV headers.

``--log``
~~~~~~~~~

  (string) Log relevant information about the HTTP request. Options are: info, debug.

  Example: ``rack servers keypair list --log info``

``--username``
~~~~~~~~~~~~~~

  (string) The Rackspace username to use for authentication.

``--api-key``
~~~~~~~~~~~~~

  (string) The Rackspace API key to use for authentication.

``--auth-tenant-id``
~~~~~~~~~~~~~~~~~~~~

  (string) The tenant ID to use for authentication. May only be provided as a command-line flag.
  (Prefixed with 'auth-' so as to not collide with the ``tenant-id``` command flags.)

``--auth-token``
~~~~~~~~~~~~~~~~

  (string) The token to use for authentication. May only be provided as a command-line flag.
  Must be used with the ``auth-tenant-id`` flag.

``--region``
~~~~~~~~~~~~

  (string) The Rackspace region to use for authentication.

``--auth-url``
~~~~~~~~~~~~~~

  (string) The Rackspace URL to use for authentication. If not provided, this
  will default to the public U.S. Rackspace endpoint.

``--profile``
~~~~~~~~~~~~~

  (string) The name of the profile (in the config file) to use to look for authentication credentials.

``--no-cache``
~~~~~~~~~~~~~~

  (boolean) Don't get or set authentication credentials in the rack cache.

``--no-header``
~~~~~~~~~~~~~~~

  (boolean) Don't set the header for CSV nor tabular output. Helpful if piping output from a ``list`` command.

``--use-service-net``
~~~~~~~~~~~~~~~~~~~~~

  (boolean) Use the Rackspace internal URL to execute the request. This will only be useful when running a
  ``rack`` command from a Rackspace server.

``--help, -h``
~~~~~~~~~~~~~~

  (boolean) Show help in a given context.

Help is available on the base level; for example::

    rack --help
    NAME:
       rack - An opinionated CLI for the Rackspace cloud

    USAGE:
       rack <command> <subcommand> <action> [flags]

    VERSION:
       0.0.0

    COMMANDS:
       servers	Used for the Servers service
       help, h	Shows a list of commands or help for one command

    GLOBAL OPTIONS:
       --json			Return output in JSON format.
       --table			Return output in tabular format. This is the default output format.
       --csv			Return output in csv format.
       --help, -h			show help

And it is available per command::

    rack servers --help
    NAME:
       rack servers - Used for the Servers service

    USAGE:
       rack servers <subcommand> <action> [flags]

    VERSION:
       0.0.0

    COMMANDS:
       instance	Used for Server Instance operations
       image	Used for Server Image operations
       flavor	Used for Server Flavor operations
       keypair	Used for Server Keypair operations
       help, h	Shows a list of commands or help for one command


And again, per subcommand:

    rack servers keypair --help
    NAME:
       rack servers keypair - Used for Server Keypair operations

    USAGE:
       rack servers keypair <action> [flags]

    VERSION:
       0.0.0

    COMMANDS:
       list		rack servers keypair list [flags]
       create	rack servers keypair create <keypairName> [flags]
       get		rack [globals] servers keypair get [--name <keypairName>] [flags]
       delete	rack servers keypair delete [--name <keypairName>] [flags]
       help, h	Shows a list of commands or help for one command


.. JSON: http://json.org/
