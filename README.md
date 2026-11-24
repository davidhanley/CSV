# CSV

FLAG CSV library: [MUSE](https://github.com/davidhanley/FLAG) fetches this repo
and imports `csv.lib`.

```clojure
{:libraries ["csv"]}
```

```clojure
{:namespace "app"
 :imports   ["csv.lib"]}

(csv/read-csv-path "data.csv")
(csv/read-csv rdr {:fields-per-record -1 :lazy-quotes true})
```

Layout:

- `libraries/csv.lib` — FLAG surface (`:exports` / `:go-exports`)
- `*.go` — pure Go reader (`github.com/davidhanley/CSV`), no FLAG runtime
