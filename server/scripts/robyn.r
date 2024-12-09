# Initialize Python and nevergrad
library(reticulate)
use_python("/usr/bin/python3")
py_run_string("
import nevergrad as ng
")
assign("ng", import("nevergrad"), envir = .GlobalEnv)

library(Robyn)
library(data.table)
Sys.setenv(TZ="Etc/UTC")

Sys.setenv(R_FUTURE_FORK_ENABLE = "true")
options(future.fork.enable = TRUE)

# For debugging
options(error = function() traceback(2))

create_files <- TRUE

# Input data will be provided directly by the service
dt_simulated_weekly <- fread("{{.InputDataPath}}")
dt_simulated_weekly[, DATE := as.Date(DATE)]

# Ensure numeric columns are properly formatted
numeric_cols <- c("conversions", "facebook_s", "googlesearch_s", "facebook_i", "googlesearch_i")
dt_simulated_weekly[, (numeric_cols) := lapply(.SD, as.numeric), .SDcols = numeric_cols]

# Load holidays data
load("{{.HolidaysDataPath}}")

robyn_directory <- "{{.OutputDirectory}}"

print("=== Debug: Initial Setup Complete ===")

InputCollect <- robyn_inputs(
  dt_input = dt_simulated_weekly,
  dt_holidays = dt_prophet_holidays,
  date_var = "DATE",
  dep_var = "{{.DepVar}}",
  dep_var_type = "{{.DepVarType}}",
  prophet_vars = c({{.ProphetVars}}),
  prophet_country = "{{.ProphetCountry}}",
  paid_media_spends = c({{.PaidMediaSpends}}),
  paid_media_vars = c({{.PaidMediaVars}}),
  window_start = "{{.WindowStart}}",
  window_end = "{{.WindowEnd}}",
  adstock = "{{.AdstockType}}"
)

print("=== Debug: InputCollect Created ===")

hyperparameters <- list(
  {{.Hyperparameters}}
  train_size = c({{.TrainSizeMin}}, {{.TrainSizeMax}})
)

print("=== Debug: Hyperparameters ===")
print(hyperparameters)

# Ensure nevergrad is available in the environment
print("=== Debug: Checking nevergrad ===")
print(exists("ng"))
print(class(ng))

InputCollect <- robyn_inputs(InputCollect = InputCollect, hyperparameters = hyperparameters)

print("=== Debug: InputCollect Updated ===")

# Ensure data is properly transformed for training
if (is.null(InputCollect$dt_transform)) {
  InputCollect$dt_transform <- InputCollect$dt_mod
}

print("=== Debug: Starting Model Run ===")

OutputModels <- robyn_run(
  InputCollect = InputCollect,
  cores = {{.Cores}},
  iterations = {{.Iterations}},
  trials = {{.Trials}},
  ts_validation = as.logical("{{.TSValidation}}"),
  add_penalty_factor = as.logical("{{.AddPenaltyFactor}}")
)

print("=== Debug: Model Run Complete ===")

OutputCollect <- robyn_outputs(
  InputCollect, OutputModels,
  pareto_fronts = "{{.ParetoFronts}}",
  csv_out = "{{.CSVOut}}",
  clusters = as.logical("{{.Clusters}}"),
  export = create_files,
  plot_folder = robyn_directory,
  plot_pareto = create_files
)

# Select best model based on NRMSE and DECOMP.RSSD
all_models <- OutputCollect$allSolutions
best_model <- all_models[1] # Default to first model if no better selection criteria

# Create output directory if it doesn't exist
dir.create(file.path(robyn_directory), showWarnings = FALSE, recursive = TRUE)

# Run allocator only for the best model
tryCatch({
  AllocatorCollect <- robyn_allocator(
    InputCollect = InputCollect,
    OutputCollect = OutputCollect,
    select_model = best_model,
    channel_constr_low = 0.7,
    channel_constr_up = 1.5,
    scenario = "max_response",
    export = create_files
  )
  
  # Extract relevant metrics for the selected model
  model_metrics <- list(
    model_id = best_model,
    nrmse = OutputModels$allMetrics[[best_model]]$nrmse,
    rssd = OutputModels$allMetrics[[best_model]]$decomp.rssd,
    recommended_spend = AllocatorCollect$dt_optimOut$optmSpendUnit
  )
  
  # Write results to JSON
  json_output <- jsonlite::toJSON(
    list(
      status = "success",
      results = model_metrics
    ),
    auto_unbox = TRUE
  )
  
  writeLines(json_output, file.path(robyn_directory, "output.json"))
}, error = function(e) {
  error_json <- jsonlite::toJSON(
    list(
      status = "error",
      message = as.character(e)
    ),
    auto_unbox = TRUE
  )
  writeLines(error_json, file.path(robyn_directory, "output.json"))
}) 