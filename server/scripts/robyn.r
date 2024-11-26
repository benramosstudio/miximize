# Input variables will be injected here
paid_media_spends <- c({{range $i, $v := .Variables}}{{if $i}},{{end}}{{$v}}{{end}})

# Test calculation (just doubles each value)
results <- paid_media_spends * 2

# Write results to JSON
output <- list(
  results = results,
  status = "success"
)

writeLines(jsonlite::toJSON(output), "output.json")

