feedtr - A Go Language Feed Transformer
=======================================

feedtr fetches feed sources (Atom, RSS or anything using XML), transforms them
using XSLT and saves the outputs.

feedtr is implemented in Go as a library with a command line front end.


Command Line
------------

The source code for the command line interface can be found in
`cmd/feedtr/main.go`.

After building, run with `feedtr --config-file=config.json` or
`feedtr < config.json`.

The config file specifies a number of outputs to be produced. Each output is
created by fetching a source and transforming it.

For example, fetching the GitHub public timeline and filtering it to only
include fork events:

**config.json:**

```JSON
{
  "outputs": [
    {
      "name": "github-fork-events.atom",
      "source": "https://github.com/timeline",
      "transforms": ["github-fork-events.xslt"]
    }
  ]
}
```

If using the `--config-file` option, the files referenced in `"transforms"`
should be placed in a directory named `transforms` located in the same directory
as the config file. If loading the config from standard input, the `transforms`
directory should be in the current directory.

**transforms/github-fork-events.xslt:**

```XSLT
<xsl:stylesheet version="1.0" xmlns:xsl="http://www.w3.org/1999/XSL/Transform" xmlns:atom="http://www.w3.org/2005/Atom">
  <!-- Copy everything unless explicitly handled below -->
  <xsl:template match="@*|*|processing-instruction()|comment()">
    <xsl:copy>
      <xsl:apply-templates select="*|@*|text()|processing-instruction()|comment()"/>
    </xsl:copy>
  </xsl:template>

  <!-- Adjust the title -->
  <xsl:template match="//atom:feed/atom:title">
    <atom:title>GitHub Public Timeline Fork Events Feed</atom:title>
  </xsl:template>

  <!-- Exclude any <entry>s that don't relate to ForkEvents -->
  <xsl:template match="atom:entry[not(contains(atom:id, ':ForkEvent/'))]"/>
</xsl:stylesheet>
```

When `feedtr` runs, it will create a `cache` directory to store feeds that have
been fetched and a `outputs` directory to store the transformed outputs.

In the above example, a file named `github-fork-events.atom` will be created in
the `outputs` directory. This will contain the result of running the XSLT
transform.

Transforms will only be run if the source has changed (based on the date
returned over HTTP) since the last output file was written. If the configuration
or transforms are changed, the relevant outputs should be manually deleted.


Library
-------

Please refer to the code documentation along with the `cmd/feedtr/main.go`
source code for example usage.
