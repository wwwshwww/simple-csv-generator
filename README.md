# Usage

Installation:

```bash
go install github.com/wwwshwww/simple-csv-generator@latest
```

Prepare Template YAML file such below:

```yaml
columns:
  - name: id
    type: INT
    choices:
      - 1
  - name: nickname
    type: STRING
  - name: name
    type: STRING
    choices:
      - Alice
      - Bob
      - Charlie
```

Run simple-csv-generator

```bash
simple-csv-generator gen -t {{template YAML file path}} -r {{row count num}}
```

For example:

```bash
simple-csv-generator gen -t your_csv_format.yaml -r 15
```

You can generate the CSV like below based on random.

```csv
id,nickname,name
1,biikkf,Bob
1,biikkf,Charlie
1,hovkzb,Charlie
1,biikkf,Bob
1,oyuvrt,Alice
1,ytmxpo,Bob
1,oyuvrt,Alice
1,mjdekw,Alice
1,hovkzb,Alice
1,oyuvrt,Bob
1,tclktr,Bob
1,tclktr,Alice
1,mjdekw,Alice
1,oyuvrt,Bob
1,tclktr,Bob
```

For more information about the Template YAML format, see `example.yaml`.
