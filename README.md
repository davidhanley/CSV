# flag-csv

FLAG CSV library. Fetched by MUSE from git and imported as `csv.lib`.

    {:libraries ["csv"]}

    {:namespace "app"
     :imports   ["csv.lib"]}

    (csv/read-csv-path "data.csv")

The FLAG surface is libraries/csv.lib. Host adapters still live in flag-lang
(runtime/csv_bind.go and libraries/csv).
